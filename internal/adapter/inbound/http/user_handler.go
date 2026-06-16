package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"hexagonalarchitecture/internal/adapter/inbound/http/dto"
	"hexagonalarchitecture/internal/core/port"
)

type UserHandler struct {
	users port.UserService
}

func NewUserHandler(users port.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// CreateUser godoc
// @Summary Create user
// @Description Create a new user.
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create user payload"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(port.ErrCodeBadRequest, port.ErrMessageInvalidRequestParams))
		return
	}

	user, err := h.users.Create(c.Request.Context(), port.CreateUserInput{
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

// FindAllUsers godoc
// @Summary List users
// @Description Get all users.
// @Tags Users
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) FindAll(c *gin.Context) {
	users, err := h.users.FindAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponses(users))
}

// FindUserByID godoc
// @Summary Get user by ID
// @Description Get one user by ID.
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) FindByID(c *gin.Context) {
	user, err := h.users.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user by ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param payload body dto.UpdateUserRequest true "Update user payload"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(port.ErrCodeBadRequest, port.ErrMessageInvalidRequestParams))
		return
	}

	user, err := h.users.Update(c.Request.Context(), c.Param("id"), port.UpdateUserInput{
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete an existing user by ID.
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.users.Delete(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func respondError(c *gin.Context, err error) {
	var appErr *port.AppError
	if errors.As(err, &appErr) {
		c.JSON(httpStatusFromAppError(appErr), newErrorResponse(appErr.Code, appErr.Message))
		return
	}

	c.JSON(
		http.StatusInternalServerError,
		newErrorResponse(port.ErrCodeInternalServer, port.ErrMessageInternalServer),
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
