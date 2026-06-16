package port

import (
	"context"

	"hexagonalarchitecture/internal/core/domain"
)

// UserRepository is the outbound port used by application services.
type AppRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	Delete(ctx context.Context, id string) error
}
