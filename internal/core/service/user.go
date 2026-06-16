package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
	"strings"
)

func (s *appService) Create(ctx context.Context, input port.CreateUserInput) (domain.User, error) {
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, err
	}

	user := domain.NewUser(s.ids.NewID(), input.Name, input.Email, s.clock.Now())
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	if err := s.publisher.PublishUserCreated(ctx, createdUser); err != nil {
		s.logger.Error("failed to publish user created event", "user_id", createdUser.ID, "error", err)
	}

	return createdUser, nil
}

func (s *appService) FindAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *appService) FindByID(ctx context.Context, id string) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, fmt.Errorf(errIDRequired, domain.ErrInvalidInput)
	}

	return s.repo.FindByID(ctx, id)
}

func (s *appService) Update(ctx context.Context, id string, input port.UpdateUserInput) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, fmt.Errorf(errIDRequired, domain.ErrInvalidInput)
	}
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, err
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	user.Update(input.Name, input.Email, s.clock.Now())
	return s.repo.Update(ctx, user)
}

func (s *appService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf(errIDRequired, domain.ErrInvalidInput)
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
