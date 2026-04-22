package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/dto"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"github.com/kodia-studio/kodia/pkg/response"
	"github.com/kodia-studio/kodia/pkg/validation"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	service  ports.NotificationService
	validate *validation.Validator
	log      *zap.Logger
}

func NewNotificationHandler(service ports.NotificationService, validate *validation.Validator, log *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		service:  service,
		validate: validate,
		log:      log,
	}
}

// List godoc
// @Summary      List notifications
// @Description  Get paginated notifications for the authenticated user
// @Tags         notifications
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        per_page query int false "Items per page (default: 10)"
// @Success      200 {object} response.Response{data=[]dto.NotificationResponse}
// @Failure      401 {object} response.Response
// @Security     Bearer
// @Router       /api/notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	params := pagination.FromContext(c)

	notifications, total, err := h.service.GetAll(c.Request.Context(), userID.(string), params)
	if err != nil {
		h.log.Error("Failed to fetch notifications", zap.String("userID", userID.(string)), zap.Error(err))
		response.InternalServerError(c, "Failed to fetch notifications")
		return
	}

	response.OK(c, "Notifications retrieved", map[string]interface{}{
		"notifications": dto.MapNotificationsToResponse(notifications),
		"total":         total,
		"page":          params.Page,
		"per_page":      params.PerPage,
	})
}

// UnreadCount godoc
// @Summary      Get unread notification count
// @Description  Get the number of unread notifications for the authenticated user
// @Tags         notifications
// @Produce      json
// @Success      200 {object} response.Response{data=dto.UnreadCountResponse}
// @Failure      401 {object} response.Response
// @Security     Bearer
// @Router       /api/notifications/unread-count [get]
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	count, err := h.service.CountUnread(c.Request.Context(), userID.(string))
	if err != nil {
		h.log.Error("Failed to count unread notifications", zap.String("userID", userID.(string)), zap.Error(err))
		response.InternalServerError(c, "Failed to count unread notifications")
		return
	}

	response.OK(c, "Unread count retrieved", dto.UnreadCountResponse{Count: count})
}

// MarkAsRead godoc
// @Summary      Mark notification as read
// @Description  Mark a specific notification as read
// @Tags         notifications
// @Produce      json
// @Param        id path string true "Notification ID"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     Bearer
// @Router       /api/notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Notification ID is required", nil)
		return
	}

	err := h.service.MarkAsRead(c.Request.Context(), id, userID.(string))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "Notification not found")
			return
		}
		h.log.Error("Failed to mark notification as read", zap.String("id", id), zap.Error(err))
		response.InternalServerError(c, "Failed to mark notification as read")
		return
	}

	response.OK(c, "Notification marked as read", nil)
}

// MarkAllAsRead godoc
// @Summary      Mark all notifications as read
// @Description  Mark all notifications for the user as read
// @Tags         notifications
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Security     Bearer
// @Router       /api/notifications/read-all [put]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	err := h.service.MarkAllAsRead(c.Request.Context(), userID.(string))
	if err != nil {
		h.log.Error("Failed to mark all notifications as read", zap.String("userID", userID.(string)), zap.Error(err))
		response.InternalServerError(c, "Failed to mark all notifications as read")
		return
	}

	response.OK(c, "All notifications marked as read", nil)
}

// Delete godoc
// @Summary      Delete notification
// @Description  Delete a specific notification
// @Tags         notifications
// @Produce      json
// @Param        id path string true "Notification ID"
// @Success      204 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     Bearer
// @Router       /api/notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Notification ID is required", nil)
		return
	}

	err := h.service.Delete(c.Request.Context(), id, userID.(string))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "Notification not found")
			return
		}
		h.log.Error("Failed to delete notification", zap.String("id", id), zap.Error(err))
		response.InternalServerError(c, "Failed to delete notification")
		return
	}

	c.Status(http.StatusNoContent)
}
