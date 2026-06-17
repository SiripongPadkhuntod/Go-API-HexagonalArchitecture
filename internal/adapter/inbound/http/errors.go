package http

import "hexagonalarchitecture/internal/core/port"

type ErrorResponse struct {
	Code    port.ErrorCode    `json:"code" example:"9988"`
	Message port.ErrorMessage `json:"message" example:"Invalid request parameters."`
}

func newErrorResponse(code port.ErrorCode, message port.ErrorMessage) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
