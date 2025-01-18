package task

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/repository"
	"github.com/muhrifqii/tuskar/internal/repository/redis"
	pureRedis "github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

type Service struct {
	taskRepository repository.TaskRepository
	redis          *redis.StorageRedis
	log            *zap.Logger
}

func NewService(taskRepository repository.TaskRepository, redis *redis.StorageRedis, zap *zap.Logger) *Service {
	return &Service{
		taskRepository: taskRepository,
		redis:          redis,
		log:            zap,
	}
}

func cacheKey(id string) string {
	return "task:" + id
}

func (s *Service) GetAllTasks(ctx context.Context, params *domain.TaskQueryParams) ([]domain.Task, int, error) {
	return s.taskRepository.GetAllTasks(ctx, params)
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.Task, error) {
	var taskToBeCached domain.Task
	cacheKey := cacheKey(id)

	cached, err := s.redis.Get(cacheKey)
	if err != nil {
		return domain.Task{}, err
	}

	if cached != nil {
		err = sonic.Unmarshal(cached, &taskToBeCached)
		return taskToBeCached, err
	}

	stored, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	byted, err := sonic.Marshal(stored)

	s.redis.Set(cacheKey, byted, time.Minute*10)
	return stored, err
}

func (s *Service) CreateTask(ctx context.Context, task *domain.Task) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- s.taskRepository.CreateTask(ctx, task)
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *Service) UpdateTask(ctx context.Context, task *domain.Task) error {
	err := s.taskRepository.UpdateTask(ctx, task)
	if err == nil {
		cacheKey := cacheKey(task.ID)
		if cacheErr := s.redis.Delete(cacheKey); cacheErr != nil && cacheErr != pureRedis.Nil {
			s.log.Warn("Failed to invalidate cache", zap.Error(cacheErr))
		}
	}
	return err
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	err := s.taskRepository.DeleteTask(ctx, id)

	if err == nil {
		cacheKey := cacheKey(id)
		if cacheErr := s.redis.Delete(cacheKey); cacheErr != nil && cacheErr != pureRedis.Nil {
			s.log.Warn("Failed to invalidate cache", zap.Error(cacheErr))
		}
	}
	return err
}
