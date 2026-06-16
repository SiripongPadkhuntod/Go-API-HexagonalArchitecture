package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/usecase"
	"strings"
)

func (s *appService) Update(ctx context.Context, id string, input usecase.UpdateUserInput) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, usecase.ToAppError(fmt.Errorf(errIDRequired, domain.ErrInvalidInput))
	}
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, usecase.ToAppError(err)
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, usecase.ToAppError(err)
	}

	user.Update(input.Name, input.Email, s.clock.Now())
	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return domain.User{}, usecase.ToAppError(err)
	}

	return updatedUser, nil
}
