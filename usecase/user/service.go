package user

import (
	"context"

	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/repository"
)

type Service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) GetUserByIdentifier(ctx context.Context, identifier string) (*domain.User, error) {
	return s.userRepository.GetByUsername(ctx, identifier)
}

func (s *Service) CreateUser(ctx context.Context, user domain.User) error {
	return s.userRepository.CreateUser(ctx, &user)
}
