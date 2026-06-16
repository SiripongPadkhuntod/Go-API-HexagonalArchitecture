package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/usecase"
	"strings"
)

func (s *appService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return usecase.ToAppError(fmt.Errorf(errIDRequired, domain.ErrInvalidInput))
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return usecase.ToAppError(err)
	}

	return nil
}
