package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/response"
)

// PaymentHandler handles Midtrans payment webhook notifications.
type PaymentHandler struct {
	orderService *services.OrderService
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler(svc *services.OrderService) *PaymentHandler {
	return &PaymentHandler{orderService: svc}
}

// HandleWebhook godoc
// @Summary Midtrans payment webhook receiver
// @Tags payments
// @Accept json
// @Router /payments/webhook [post]
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		response.BadRequest(c, "Invalid notification payload", nil)
		return
	}

	if err := h.orderService.HandlePaymentWebhook(c.Request.Context(), notification); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Webhook processed", nil)
}
