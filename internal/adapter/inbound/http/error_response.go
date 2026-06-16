package http

import "hexagonalarchitecture/internal/core/usecase"

type ErrorResponse struct {
	Code    usecase.ErrorCode    `json:"code" example:"9988"`
	Message usecase.ErrorMessage `json:"message" example:"Invalid request parameters."`
}

func newErrorResponse(code usecase.ErrorCode, message usecase.ErrorMessage) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
