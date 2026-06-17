package noop

import (
	"context"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

var _ port.UserEventPublisher = (*Client)(nil)

func (c *Client) PublishUserCreated(ctx context.Context, user domain.User) error {
	return nil
}
