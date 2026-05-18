package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/transformers"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/response"
)

// CouponHandler handles coupon-related HTTP requests.
type CouponHandler struct {
	couponService *services.CouponService
}

// NewCouponHandler creates a new CouponHandler.
func NewCouponHandler(svc *services.CouponService) *CouponHandler {
	return &CouponHandler{couponService: svc}
}

type validateCouponRequest struct {
	Code     string `json:"code" binding:"required"`
	Subtotal int64  `json:"subtotal" binding:"required"`
}

// ValidateCoupon godoc
// @Summary Validate a coupon code
// @Tags coupons
// @Security BearerAuth
// @Accept json
// @Router /coupons/validate [post]
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	var req validateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.couponService.Validate(c.Request.Context(), req.Code, req.Subtotal)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Coupon valid", gin.H{
		"coupon":          transformers.CouponToResponse(result.Coupon),
		"discount_amount": result.DiscountAmount,
	})
}

// ListCoupons godoc
// @Summary List coupons (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/coupons [get]
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page, _ := parsePageParams(c)
	perPage := 50

	coupons, total, err := h.couponService.List(c.Request.Context(), page, perPage)
	if err != nil {
		response.InternalServerError(c, "Failed to list coupons")
		return
	}

	response.OKWithMeta(c, "Coupons retrieved", transformers.CouponsToResponse(coupons),
		response.NewMeta(page, perPage, total))
}

type createCouponRequest struct {
	Code           string  `json:"code" binding:"required"`
	Description    string  `json:"description"`
	Type           string  `json:"type" binding:"required"`
	Value          int64   `json:"value" binding:"required"`
	MinOrderAmount int64   `json:"min_order_amount"`
	MaxDiscount    int64   `json:"max_discount"`
	MaxUsage       int     `json:"max_usage"`
	StartAt        *string `json:"start_at"`
	ExpiresAt      *string `json:"expires_at"`
	IsActive       bool    `json:"is_active"`
}

// CreateCoupon godoc
// @Summary Create a coupon (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/coupons [post]
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req createCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	coupon := &domain.Coupon{
		Code:           req.Code,
		Description:    req.Description,
		Type:           domain.DiscountType(req.Type),
		Value:          req.Value,
		MinOrderAmount: req.MinOrderAmount,
		MaxDiscount:    req.MaxDiscount,
		MaxUsage:       req.MaxUsage,
		IsActive:       req.IsActive,
	}

	if req.StartAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.StartAt); err == nil {
			coupon.StartAt = &t
		}
	}
	if req.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			coupon.ExpiresAt = &t
		}
	}

	if err := h.couponService.Create(c.Request.Context(), coupon); err != nil {
		response.InternalServerError(c, "Failed to create coupon")
		return
	}

	response.Created(c, "Coupon created", transformers.CouponToResponse(coupon))
}

// UpdateCoupon godoc
// @Summary Update a coupon (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/coupons/{id} [patch]
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	fields := make(map[string]any)
	for k, v := range body {
		switch k {
		case "type":
			if s, ok := v.(string); ok {
				fields["type"] = domain.DiscountType(s)
			}
		case "value", "min_order_amount", "max_discount":
			if f, ok := v.(float64); ok {
				fields[k] = int64(f)
			}
		case "max_usage":
			if f, ok := v.(float64); ok {
				fields[k] = int(f)
			}
		case "start_at", "expires_at":
			if s, ok := v.(string); ok && s != "" {
				if t, err := time.Parse(time.RFC3339, s); err == nil {
					fields[k] = t
				}
			} else {
				fields[k] = nil
			}
		default:
			fields[k] = v
		}
	}

	coupon, err := h.couponService.Update(c.Request.Context(), c.Param("id"), fields)
	if err != nil {
		response.NotFound(c, "Coupon not found")
		return
	}

	response.OK(c, "Coupon updated", transformers.CouponToResponse(coupon))
}

// DeleteCoupon godoc
// @Summary Delete a coupon (admin)
// @Tags admin
// @Security BearerAuth
// @Router /admin/coupons/{id} [delete]
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	if err := h.couponService.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.NotFound(c, "Coupon not found")
		return
	}
	response.NoContent(c)
}
