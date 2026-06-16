package port

import (
	"errors"
	"fmt"

	"hexagonalarchitecture/internal/core/domain"
)

type ErrorCode string
type ErrorMessage string
type ErrorKind string

const (
	ErrorKindBusiness  ErrorKind = "business"
	ErrorKindTechnical ErrorKind = "technical"

	ErrMessageInvalidRequestParams ErrorMessage = "Invalid request parameters."
	ErrMessageInvalidInput         ErrorMessage = "Invalid input"
	ErrMessageUserNotFound         ErrorMessage = "User not found"
	ErrMessageUserAlreadyExists    ErrorMessage = "Email is already in use"
	ErrMessageInternalServer       ErrorMessage = "Internal Server Error"

	ErrCodeBadRequest        ErrorCode = "9988"
	ErrCodeInvalidInput      ErrorCode = "9987"
	ErrCodeUserNotFound      ErrorCode = "9984"
	ErrCodeUserAlreadyExists ErrorCode = "9983"
	ErrCodeInternalServer    ErrorCode = "5000"
)

type AppError struct {
	Kind    ErrorKind
	Code    ErrorCode
	Message ErrorMessage
	Cause   error
}

func NewBusinessError(code ErrorCode, message ErrorMessage, cause error) *AppError {
	return &AppError{Kind: ErrorKindBusiness, Code: code, Message: message, Cause: cause}
}

func NewTechnicalError(cause error) *AppError {
	return &AppError{
		Kind:    ErrorKindTechnical,
		Code:    ErrCodeInternalServer,
		Message: ErrMessageInternalServer,
		Cause:   cause,
	}
}

func (e *AppError) Error() string {
	if e.Cause == nil {
		return string(e.Message)
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *AppError) Unwrap() error { return e.Cause }

func ToAppError(err error) error {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return err
	}
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return NewBusinessError(ErrCodeInvalidInput, ErrMessageInvalidInput, err)
	case errors.Is(err, domain.ErrUserNotFound):
		return NewBusinessError(ErrCodeUserNotFound, ErrMessageUserNotFound, err)
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return NewBusinessError(ErrCodeUserAlreadyExists, ErrMessageUserAlreadyExists, err)
	default:
		return NewTechnicalError(err)
	}
}
