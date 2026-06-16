package usecase

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

	ERROR_MESSAGE_INVALID_REQUEST_PARAMS ErrorMessage = "Invalid request parameters."
	ERROR_MESSAGE_INVALID_INPUT          ErrorMessage = "Invalid input"
	ERROR_MESSAGE_USER_NOT_FOUND         ErrorMessage = "User not found"
	ERROR_MESSAGE_USER_ALREADY_EXISTS    ErrorMessage = "Email is already in use"
	ERROR_MESSAGE_INTERNAL_SERVER_ERROR  ErrorMessage = "Internal Server Error"

	ERROR_CODE_BAD_REQUEST           ErrorCode = "9988"
	ERROR_CODE_INVALID_INPUT         ErrorCode = "9987"
	ERROR_CODE_USER_NOT_FOUND        ErrorCode = "9984"
	ERROR_CODE_USER_ALREADY_EXISTS   ErrorCode = "9983"
	ERROR_CODE_INTERNAL_SERVER_ERROR ErrorCode = "5000"
)

type AppError struct {
	Kind    ErrorKind
	Code    ErrorCode
	Message ErrorMessage
	Cause   error
}

func NewBusinessError(code ErrorCode, message ErrorMessage, cause error) *AppError {
	return &AppError{
		Kind:    ErrorKindBusiness,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func NewTechnicalError(cause error) *AppError {
	return &AppError{
		Kind:    ErrorKindTechnical,
		Code:    ERROR_CODE_INTERNAL_SERVER_ERROR,
		Message: ERROR_MESSAGE_INTERNAL_SERVER_ERROR,
		Cause:   cause,
	}
}

func (e *AppError) Error() string {
	if e.Cause == nil {
		return string(e.Message)
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

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
		return NewBusinessError(ERROR_CODE_INVALID_INPUT, ERROR_MESSAGE_INVALID_INPUT, err)
	case errors.Is(err, domain.ErrUserNotFound):
		return NewBusinessError(ERROR_CODE_USER_NOT_FOUND, ERROR_MESSAGE_USER_NOT_FOUND, err)
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return NewBusinessError(ERROR_CODE_USER_ALREADY_EXISTS, ERROR_MESSAGE_USER_ALREADY_EXISTS, err)
	default:
		return NewTechnicalError(err)
	}
}
