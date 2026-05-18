package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/http/transformers"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/response"
)

// OrderHandler handles order-related HTTP requests.
type OrderHandler struct {
	orderService *services.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(svc *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: svc}
}

type createOrderRequest struct {
	Items []struct {
		ProductID string `json:"product_id"`
		VariantID string `json:"variant_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items" binding:"required"`
	CouponCode string `json:"coupon_code"`
}

// CreateOrder godoc
// @Summary Create a new order
// @Tags orders
// @Security BearerAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	userID := middleware.GetUserID(c)

	orderReq := &services.CreateOrderRequest{
		UserID:     userID,
		CouponCode: req.CouponCode,
	}
	for _, item := range req.Items {
		qty := item.Quantity
		if qty < 1 {
			qty = 1
		}
		orderReq.Items = append(orderReq.Items, struct {
			ProductID string `json:"product_id"`
			VariantID string `json:"variant_id"`
			Quantity  int    `json:"quantity"`
		}{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Quantity:  qty,
		})
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), orderReq)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Order created", transformers.OrderToResponse(order))
}

// GetMyOrders godoc
// @Summary List authenticated user's orders
// @Tags orders
// @Security BearerAuth
// @Router /orders [get]
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	page, perPage := parsePageParams(c)
	userID := middleware.GetUserID(c)

	orders, total, err := h.orderService.GetMyOrders(c.Request.Context(), userID, page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to list orders")
		return
	}

	response.OKWithMeta(c, "Orders retrieved", transformers.OrdersToResponse(orders),
		response.NewMeta(page, perPage, total))
}

// GetMyOrder godoc
// @Summary Get a specific order for the authenticated user
// @Tags orders
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetMyOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	order, err := h.orderService.GetMyOrder(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Order not found")
		return
	}
	response.OK(c, "Order retrieved", transformers.OrderToResponse(order))
}

type initPaymentRequest struct {
	CustomerName  string `json:"customer_name" binding:"required"`
	CustomerEmail string `json:"customer_email" binding:"required"`
	CustomerPhone string `json:"customer_phone"`
}

// InitPayment godoc
// @Summary Initialize Midtrans payment for an order
// @Tags orders
// @Security BearerAuth
// @Router /orders/{id}/pay [post]
func (h *OrderHandler) InitPayment(c *gin.Context) {
	var req initPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.orderService.InitPayment(c.Request.Context(), &services.InitPaymentRequest{
		OrderID:       c.Param("id"),
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
	})
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Payment initialized", gin.H{
		"snap_token": result.SnapToken,
		"snap_url":   result.SnapURL,
	})
}

// GenerateDownloadURL godoc
// @Summary Generate a temporary download URL for an order item
// @Tags orders
// @Security BearerAuth
// @Param orderId path string true "Order ID"
// @Param itemId path string true "Order Item ID"
// @Router /orders/{orderId}/items/{itemId}/download [get]
func (h *OrderHandler) GenerateDownloadURL(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID := c.Param("id")
	itemID := c.Param("itemId")

	info, err := h.orderService.GenerateDownloadURL(c.Request.Context(), userID, orderID, itemID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Download URL generated", gin.H{
		"download_url": info.DownloadURL,
		"expires_in":   300,
	})
}

// GetAdminOrders godoc
// @Summary List all orders (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/orders [get]
func (h *OrderHandler) GetAdminOrders(c *gin.Context) {
	page, perPage := parsePageParams(c)

	orders, total, err := h.orderService.GetAdminOrders(c.Request.Context(), page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to list orders")
		return
	}

	response.OKWithMeta(c, "Orders retrieved", transformers.OrdersToResponse(orders),
		response.NewMeta(page, perPage, total))
}

// GetAdminOrder godoc
// @Summary Get a specific order (admin)
// @Tags admin
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Router /admin/orders/{id} [get]
func (h *OrderHandler) GetAdminOrder(c *gin.Context) {
	order, err := h.orderService.GetAdminOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.NotFound(c, "Order not found")
		return
	}
	response.OK(c, "Order retrieved", transformers.OrderToResponse(order))
}

// DashboardAnalytics godoc
// @Summary Get dashboard analytics (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/analytics/dashboard [get]
func (h *OrderHandler) DashboardAnalytics(c *gin.Context) {
	orders, _, err := h.orderService.GetAdminOrders(c.Request.Context(), 1, 1000)
	if err != nil {
		response.InternalServerError(c, "Failed to get analytics")
		return
	}

	var totalRevenue, pendingRevenue int64
	var totalOrders, pendingOrders, paidOrders, expiredOrders, cancelledOrders int

	for _, o := range orders {
		totalOrders++
		switch o.Status {
		case "paid":
			paidOrders++
			totalRevenue += o.Total
		case "pending":
			pendingOrders++
			pendingRevenue += o.Total
		case "expired":
			expiredOrders++
		case "cancelled":
			cancelledOrders++
		}
	}

	var avgOrderValue int64
	if paidOrders > 0 {
		avgOrderValue = totalRevenue / int64(paidOrders)
	}

	response.OK(c, "Analytics retrieved", gin.H{
		"revenue_stats": gin.H{
			"total_revenue":       totalRevenue,
			"total_orders":        paidOrders,
			"average_order_value": avgOrderValue,
			"pending_revenue":     pendingRevenue,
			"currency":            "IDR",
		},
		"order_stats": gin.H{
			"total":     totalOrders,
			"pending":   pendingOrders,
			"paid":      paidOrders,
			"expired":   expiredOrders,
			"cancelled": cancelledOrders,
		},
		"monthly_sales": []gin.H{},
		"top_products":  []gin.H{},
	})
}

// parsePageParams extracts page and per_page query params with sensible defaults.
func parsePageParams(c *gin.Context) (page, perPage int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "15"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 200 {
		perPage = 15
	}
	return
}
