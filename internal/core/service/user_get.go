package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/usecase"
	"strings"
)

func (s *appService) FindAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, usecase.ToAppError(err)
	}

	return users, nil
}

func (s *appService) FindByID(ctx context.Context, id string) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, usecase.ToAppError(fmt.Errorf(errIDRequired, domain.ErrInvalidInput))
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, usecase.ToAppError(err)
	}

	return user, nil
}
