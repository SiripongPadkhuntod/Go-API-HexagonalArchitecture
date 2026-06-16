package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

const errIDRequired = "%w: id is required"

type userService struct {
	repo      port.UserRepository
	publisher port.UserEventPublisher
	logger    port.Logger
	ids       port.IDGenerator
	clock     port.Clock
}

type UserServiceDeps struct {
	Repo      port.UserRepository
	Publisher port.UserEventPublisher
	Logger    port.Logger
	IDs       port.IDGenerator
	Clock     port.Clock
}

func NewUserService(deps UserServiceDeps) port.UserService {
	return &userService{
		repo:      deps.Repo,
		publisher: deps.Publisher,
		logger:    deps.Logger,
		ids:       deps.IDs,
		clock:     deps.Clock,
	}
}

func (s *userService) Create(ctx context.Context, input port.CreateUserInput) (domain.User, error) {
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, port.ToAppError(err)
	}

	user := domain.NewUser(s.ids.NewID(), input.Name, input.Email, s.clock.Now())
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, port.ToAppError(err)
	}

	if err := s.publisher.PublishUserCreated(ctx, createdUser); err != nil {
		s.logger.Error("failed to publish user created event", "user_id", createdUser.ID, "error", err)
	}

	return createdUser, nil
}

func (s *userService) FindAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, port.ToAppError(err)
	}

	return users, nil
}

func (s *userService) FindByID(ctx context.Context, id string) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, port.ToAppError(fmt.Errorf(errIDRequired, domain.ErrInvalidInput))
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, port.ToAppError(err)
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id string, input port.UpdateUserInput) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, port.ToAppError(fmt.Errorf(errIDRequired, domain.ErrInvalidInput))
	}
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, port.ToAppError(err)
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, port.ToAppError(err)
	}

	user.Update(input.Name, input.Email, s.clock.Now())
	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return domain.User{}, port.ToAppError(err)
	}

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return port.ToAppError(fmt.Errorf(errIDRequired, domain.ErrInvalidInput))
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return port.ToAppError(err)
	}

	return nil
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
