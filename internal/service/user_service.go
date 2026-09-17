package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
)

// UserService encapsulates business logic for User CRUD operations.
type UserService interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, p repository.Pagination) ([]domain.User, int64, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	users repository.UserRepository
}

// NewUserService builds a UserService backed by the given repository.
func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	user := req.ToDomain()
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context, p repository.Pagination) ([]domain.User, int64, error) {
	return s.users.List(ctx, p)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ApplyTo(user)

	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.users.Delete(ctx, id)
}
