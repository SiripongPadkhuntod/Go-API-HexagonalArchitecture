package noop

import (
	"context"
	"net/http"

	"hexagonalarchitecture/internal/core/port"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Do(ctx context.Context, request port.OutboundAPIRequest) (port.OutboundAPIResponse, error) {
	return port.OutboundAPIResponse{StatusCode: http.StatusNoContent}, nil
}
