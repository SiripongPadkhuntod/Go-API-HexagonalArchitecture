package port

import (
	"context"

	"hexagonalarchitecture/internal/core/domain"
)

type AppService interface {
	Create(ctx context.Context, input CreateUserInput) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, id string, input UpdateUserInput) (domain.User, error)
	Delete(ctx context.Context, id string) error
}

type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserInput struct {
	Name  string
	Email string
}
