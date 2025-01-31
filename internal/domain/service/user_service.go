package service

import (
	"context"

	"github.com/Fox1N69/iq-testtask/internal/domain/entity"
	"github.com/Fox1N69/iq-testtask/internal/repository"
)

type UserService interface {
	UserByID(ctx context.Context, userID int64) (*entity.User, error)
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) UserByID(ctx context.Context, userID int64) (*entity.User, error) {
	return s.repository.UserByID(ctx, userID)
}
