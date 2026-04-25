package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/response"
	"github.com/kodia-studio/kodia/pkg/validation"
	"go.uber.org/zap"
)

// RoleHandler handles role management HTTP requests.
type RoleHandler struct {
	roleService ports.RoleService
	validate    *validation.Validator
	log         *zap.Logger
}

func NewRoleHandler(roleService ports.RoleService, validate *validation.Validator, log *zap.Logger) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		validate:    validate,
		log:         log,
	}
}

// CreateRoleRequest is the request body for creating a role.
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required,min=1,max=255"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// RoleResponse is the HTTP response for a role.
type RoleResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Create a new role with permissions (admin only)
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body CreateRoleRequest true "Create role request"
// @Success      201 {object} response.Response{data=RoleResponse}
// @Router       /api/admin/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	// Verify admin
	middleware.RequireRole("admin")(c)

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), req.Name, req.Description, req.Permissions)
	if err != nil {
		h.log.Error("Failed to create role", zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	response.Created(c, "Role created", h.mapRoleToResponse(role))
}

// GetRoles godoc
// @Summary      List all roles
// @Description  Get all defined roles (admin only)
// @Tags         roles
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=[]RoleResponse}
// @Router       /api/admin/roles [get]
func (h *RoleHandler) GetRoles(c *gin.Context) {
	// Verify admin
	middleware.RequireRole("admin")(c)

	roles, err := h.roleService.GetAllRoles(c.Request.Context())
	if err != nil {
		h.log.Error("Failed to get roles", zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = h.mapRoleToResponse(role)
	}

	response.OK(c, "Roles retrieved", responses)
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Soft-delete a role by ID (admin only)
// @Tags         roles
// @Security     BearerAuth
// @Param        id path string true "Role ID"
// @Success      200 {object} response.Response
// @Router       /api/admin/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	// Verify admin
	middleware.RequireRole("admin")(c)

	id := c.Param("id")
	if err := h.roleService.DeleteRole(c.Request.Context(), id); err != nil {
		h.log.Error("Failed to delete role", zap.String("role_id", id), zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	response.OK(c, "Role deleted", nil)
}

// AssignRoleRequest is the request body for assigning a role to a user.
type AssignRoleRequest struct {
	RoleName string `json:"role_name" binding:"required,min=1"`
}

// AssignRole godoc
// @Summary      Assign a role to a user
// @Description  Assign a role to a user (admin only)
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Param        user_id path string true "User ID"
// @Param        request body AssignRoleRequest true "Assign role request"
// @Success      200 {object} response.Response
// @Router       /api/admin/users/{user_id}/roles [post]
func (h *RoleHandler) AssignRole(c *gin.Context) {
	// Verify admin
	middleware.RequireRole("admin")(c)

	userID := c.Param("user_id")
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	if err := h.roleService.AssignRole(c.Request.Context(), userID, req.RoleName); err != nil {
		h.log.Error("Failed to assign role", zap.String("user_id", userID), zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	response.OK(c, "Role assigned", nil)
}

// RevokeRole godoc
// @Summary      Revoke a role from a user
// @Description  Revoke a role from a user (admin only)
// @Tags         roles
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Param        role    path string true "Role name"
// @Success      200 {object} response.Response
// @Router       /api/admin/users/{user_id}/roles/{role} [delete]
func (h *RoleHandler) RevokeRole(c *gin.Context) {
	// Verify admin
	middleware.RequireRole("admin")(c)

	userID := c.Param("user_id")
	roleName := c.Param("role")

	if err := h.roleService.RevokeRole(c.Request.Context(), userID, roleName); err != nil {
		h.log.Error("Failed to revoke role", zap.String("user_id", userID), zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	response.OK(c, "Role revoked", nil)
}

// GetUserRoles godoc
// @Summary      Get user roles
// @Description  Get all roles assigned to a user (admin only)
// @Tags         roles
// @Security     BearerAuth
// @Param        user_id path string true "User ID"
// @Success      200 {object} response.Response{data=[]string}
// @Router       /api/admin/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	// Verify admin
	middleware.RequireRole("admin")(c)

	userID := c.Param("user_id")
	roles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user roles", zap.String("user_id", userID), zap.Error(err))
		response.InternalServerError(c, "")
		return
	}

	response.OK(c, "User roles retrieved", roles)
}

// Helper function to convert domain role to HTTP response.
func (h *RoleHandler) mapRoleToResponse(r *domain.RoleEntity) RoleResponse {
	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
