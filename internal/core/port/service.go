package port

import (
	"context"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/usecase"
)

type AppService interface {
	Create(ctx context.Context, input usecase.CreateUserInput) (domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, id string, input usecase.UpdateUserInput) (domain.User, error)
	Delete(ctx context.Context, id string) error
}

type UserEventPublisher interface {
	PublishUserCreated(ctx context.Context, user domain.User) error
}
