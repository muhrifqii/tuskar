package provision

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/repository"
	"github.com/muhrifqii/tuskar/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	userRepository repository.UserRepository
	log            *zap.Logger
}

func NewService(
	zap *zap.Logger,
	userRepository repository.UserRepository,
) *Service {
	return &Service{
		userRepository: userRepository,
		log:            zap,
	}
}

func (s *Service) CreateSystemUser(ctx context.Context, username, password string) error {
	exist, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return err
	} else if exist != nil {
		log.Debug("User already exist. Skipping...", zap.String("username", username))
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := &domain.User{
		Username: username,
		Password: hashedPassword,
	}
	return s.userRepository.CreateUser(ctx, user)
}
