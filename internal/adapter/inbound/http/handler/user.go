package handler

import (
	"context"

	"hexagonalarchitecture/internal/adapter/inbound/http/handler/dto"
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
// @Router /api/v1/users [post]
func (h *UserHandler) Create(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	user, err := h.users.Create(ctx, port.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.ToUserResponse(user), nil
}

// FindAllUsers godoc
// @Summary List users
// @Description Get all users.
// @Tags Users
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Router /api/v1/users [get]
func (h *UserHandler) FindAll(ctx context.Context, _ struct{}) ([]dto.UserResponse, error) {
	users, err := h.users.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponses(users), nil
}

// FindUserByID godoc
// @Summary Get user by ID
// @Description Get one user by ID.
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) FindByID(ctx context.Context, req dto.GetUserRequest) (dto.UserResponse, error) {
	user, err := h.users.FindByID(ctx, req.ID)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.ToUserResponse(user), nil
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
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, err := h.users.Update(ctx, req.ID, port.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.ToUserResponse(user), nil
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete an existing user by ID.
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(ctx context.Context, req dto.DeleteUserRequest) (map[string]string, error) {
	if err := h.users.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return map[string]string{"message": "User deleted successfully"}, nil
}
