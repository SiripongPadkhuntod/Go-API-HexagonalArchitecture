package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// UserService is the inbound application port for user use cases.
type UserService interface {
	Create(ctx context.Context, input CreateUserInput) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, id string, input UpdateUserInput) (domain.User, error)
	Delete(ctx context.Context, id string) error
}

type CreateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UpdateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type userService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (domain.User, error) {
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, err
	}

	user := domain.NewUser(newID(), input.Name, input.Email)
	return s.repo.Create(ctx, user)
}

func (s *userService) FindAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *userService) FindByID(ctx context.Context, id string) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}

	return s.repo.FindByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id string, input UpdateUserInput) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, err
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user.Update(input.Name, input.Email)
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func validateUser(name, email string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if !emailPattern.MatchString(strings.TrimSpace(email)) {
		return fmt.Errorf("%w: email is invalid", domain.ErrInvalidInput)
	}

	return nil
}

func newID() string {
	return fmt.Sprintf("usr_%d", time.Now().UTC().UnixNano())
}
