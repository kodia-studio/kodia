package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/adapters/http/dto"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"github.com/kodia-studio/kodia/pkg/response"
	"github.com/kodia-studio/kodia/pkg/validation"
	"go.uber.org/zap"
)

// UserHandler handles user management HTTP requests.
type UserHandler struct {
	userService ports.UserService
	validate    *validation.Validator
	log         *zap.Logger
}

func NewUserHandler(userService ports.UserService, validate *validation.Validator, log *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
		log:         log,
	}
}

// GetAll godoc
// @Summary      List all users
// @Description  Returns a paginated list of all users (admin only)
// @Tags         users
// @Security     BearerAuth
// @Param        page     query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(15)
// @Success      200 {object} response.Response{data=[]dto.UserResponse}
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	params := pagination.FromContext(c)
	users, total, err := h.userService.GetAll(c.Request.Context(), params)
	if err != nil {
		h.log.Error("Failed to get users", zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	meta := response.NewMeta(params.Page, params.PerPage, total)
	response.OKWithMeta(c, "Users retrieved", dto.MapUsersToResponse(users), meta)
}

// GetByID godoc
// @Summary      Get user by ID
// @Tags         users
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Failure      404 {object} response.Response
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.OK(c, "User retrieved", dto.MapUserToResponse(user))
}

// GetMe godoc
// @Summary      Get authenticated user profile
// @Tags         users
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}
	response.OK(c, "Profile retrieved", dto.MapUserToResponse(user))
}

// Update godoc
// @Summary      Update user profile
// @Tags         users
// @Security     BearerAuth
// @Param        id   path string              true "User ID"
// @Param        body body dto.UpdateUserRequest true "Profile update data"
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Router       /users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Only allow users to update their own profile (unless admin)
	requestingUserID := middleware.GetUserID(c)
	requestingRole := middleware.GetUserRole(c)
	if id != requestingUserID && requestingRole != string(domain.RoleAdmin) {
		response.Forbidden(c, "You can only update your own profile")
		return
	}

	var req dto.UpdateUserRequest
	if !validation.BindAndValidate(c, h.validate, &req) {
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, ports.UpdateUserInput{
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}

	response.OK(c, "Profile updated", dto.MapUserToResponse(user))
}

// Delete godoc
// @Summary      Delete user
// @Tags         users
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      204
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.NoContent(c)
}

// ChangePassword godoc
// @Summary      Change password
// @Tags         users
// @Security     BearerAuth
// @Param        body body dto.ChangePasswordRequest true "Password change data"
// @Success      200 {object} response.Response
// @Router       /users/me/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if !validation.BindAndValidate(c, h.validate, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	err := h.userService.ChangePassword(c.Request.Context(), userID, ports.ChangePasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.BadRequest(c, "Current password is incorrect", nil)
			return
		}
		response.InternalServerError(c, "")
		return
	}

	response.OK(c, "Password changed successfully", nil)
}

// DeleteMe godoc
// @Summary      Delete own account
// @Description  Permanently delete the authenticated user's account
// @Tags         users
// @Security     BearerAuth
// @Success      204
// @Failure      401 {object} response.Response
// @Router       /users/me [delete]
func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.userService.Delete(c.Request.Context(), userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.NoContent(c)
}

// Create godoc
// @Summary      Create a new user
// @Tags         users
// @Security     BearerAuth
// @Param        body body dto.CreateUserRequest true "User creation data"
// @Success      201 {object} response.Response{data=dto.UserResponse}
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if !validation.BindAndValidate(c, h.validate, &req) {
		return
	}

	id := uuid.NewString()
	user, err := h.userService.Create(c.Request.Context(), ports.CreateUserInput{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		response.InternalServerError(c, "Failed to create user")
		return
	}

	c.JSON(201, response.Response{
		Success: true,
		Message: "User created successfully",
		Data:    dto.MapUserToResponse(user),
	})
}

// GetTrashed godoc
// @Summary      List trashed (soft-deleted) users
// @Tags         users
// @Security     BearerAuth
// @Param        page     query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(15)
// @Success      200 {object} response.Response{data=[]dto.UserResponse}
// @Router       /users/trashed [get]
func (h *UserHandler) GetTrashed(c *gin.Context) {
	params := pagination.FromContext(c)
	users, total, err := h.userService.GetTrashed(c.Request.Context(), params)
	if err != nil {
		h.log.Error("Failed to get trashed users", zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	meta := response.NewMeta(params.Page, params.PerPage, total)
	response.OKWithMeta(c, "Trashed users retrieved", dto.MapUsersToResponse(users), meta)
}

// Restore godoc
// @Summary      Restore a trashed user
// @Tags         users
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} response.Response
// @Router       /users/{id}/restore [post]
func (h *UserHandler) Restore(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.Restore(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.OK(c, "User restored successfully", nil)
}

// ForceDelete godoc
// @Summary      Permanently delete a user
// @Tags         users
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      204
// @Router       /users/{id}/force-delete [post]
func (h *UserHandler) ForceDelete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.ForceDelete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "")
		return
	}
	response.NoContent(c)
}
