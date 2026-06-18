package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"hexagonalarchitecture/internal/core/port"
)

// Bind handles the binding of URI, Query, and JSON Body, then calls the provided service function.
// It maps errors and success responses automatically.
func Bind[Req any, Res any](
	c *gin.Context,
	successStatus int,
	f func(ctx context.Context, req Req) (Res, error),
) {
	var req Req

	// 1) Bind Request (URI, Query, Body)
	if err := bindRequest(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(port.ErrCodeBadRequest, port.ErrMessageInvalidRequestParams))
		return
	}

	// 2) Call handler
	resp, err := f(c.Request.Context(), req)
	if err != nil {
		RespondError(c, err)
		return
	}

	// 3) Success response
	RespondSuccess(c, successStatus, resp)
}

// bindRequest intelligently binds URI, Query, and Body depending on the HTTP method
func bindRequest[T any](c *gin.Context, req *T) error {
	// 1) Bind URI params
	if err := c.ShouldBindUri(req); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			return err
		}
	}

	// 2) Bind Query params
	if err := c.ShouldBindQuery(req); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			return err
		}
	}

	// 3) Bind JSON body
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodDelete {
		if c.Request.Body != nil && c.Request.ContentLength != 0 {
			return c.ShouldBindJSON(req)
		}
		return validateStruct(req)
	}

	// For GET/DELETE, manually validate struct since ShouldBindJSON wasn't called
	return validateStruct(req)
}

func validateStruct(obj any) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.Struct(obj)
	}
	return nil
}

// ---- Response & Error Formatting ----

// SuccessResponse represents a standard success response wrapper
type SuccessResponse struct {
	Status string `json:"status" example:"success"`
	Data   any    `json:"data"`
}

// RespondSuccess is a helper function to standardise success responses
func RespondSuccess(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, SuccessResponse{
		Status: "success",
		Data:   data,
	})
}

type ErrorResponse struct {
	Code    port.ErrorCode    `json:"code" example:"9988"`
	Message port.ErrorMessage `json:"message" example:"Invalid request parameters."`
}

func NewErrorResponse(code port.ErrorCode, message port.ErrorMessage) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func RespondError(c *gin.Context, err error) {
	var appErr *port.AppError
	if errors.As(err, &appErr) {
		c.JSON(httpStatusFromAppError(appErr), NewErrorResponse(appErr.Code, appErr.Message))
		return
	}

	c.JSON(
		http.StatusInternalServerError,
		NewErrorResponse(port.ErrCodeInternalServer, port.ErrMessageInternalServer),
	)
}

func httpStatusFromAppError(err *port.AppError) int {
	if err.Kind == port.ErrorKindTechnical {
		return http.StatusInternalServerError
	}

	switch err.Code {
	case port.ErrCodeUserNotFound:
		return http.StatusNotFound
	case port.ErrCodeUserAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
