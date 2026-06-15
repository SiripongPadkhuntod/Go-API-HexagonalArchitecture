package port

import (
	"context"

	"hexagonalarchitecture/internal/core/domain"
)

type UserEventPublisher interface {
	PublishUserCreated(ctx context.Context, user domain.User) error
}
