package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"strings"
)

func (s *appService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf(errIDRequired, domain.ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}
