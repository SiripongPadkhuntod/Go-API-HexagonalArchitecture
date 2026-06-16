package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"strings"
)

func (s *appService) FindAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *appService) FindByID(ctx context.Context, id string) (domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return domain.User{}, fmt.Errorf(errIDRequired, domain.ErrInvalidInput)
	}

	return s.repo.FindByID(ctx, id)
}
