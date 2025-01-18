package repository

import (
	"context"

	"github.com/muhrifqii/tuskar/domain"
)

type (
	UserRepository interface {
		GetByUsername(ctx context.Context, username string) (*domain.User, error)
		CreateUser(ctx context.Context, user *domain.User) error
	}

	TaskRepository interface {
		GetAllTasks(ctx context.Context, params *domain.TaskQueryParams) ([]domain.Task, int, error)
		GetByID(ctx context.Context, id string) (domain.Task, error)
		CreateTask(ctx context.Context, task *domain.Task) error
		UpdateTask(ctx context.Context, task *domain.Task) error
		DeleteTask(ctx context.Context, id string) error
	}
)
