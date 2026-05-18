package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

type CouponService struct {
	repo ports.CouponRepository
}

func NewCouponService(repo ports.CouponRepository) *CouponService {
	return &CouponService{repo: repo}
}

type ValidateCouponResult struct {
	Coupon       *domain.Coupon `json:"coupon"`
	DiscountAmount int64        `json:"discount_amount"`
}

func (s *CouponService) Validate(ctx context.Context, code string, subtotal int64) (*ValidateCouponResult, error) {
	coupon, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("coupon not found")
	}

	if !coupon.IsActive {
		return nil, fmt.Errorf("coupon is not active")
	}

	now := time.Now()

	if coupon.StartAt != nil && now.Before(*coupon.StartAt) {
		return nil, fmt.Errorf("coupon is not yet active")
	}

	if coupon.ExpiresAt != nil && now.After(*coupon.ExpiresAt) {
		return nil, fmt.Errorf("coupon has expired")
	}

	if coupon.MaxUsage > 0 && coupon.UsageCount >= coupon.MaxUsage {
		return nil, fmt.Errorf("coupon usage limit exceeded")
	}

	if subtotal < coupon.MinOrderAmount {
		return nil, fmt.Errorf("order amount is below minimum required for this coupon")
	}

	// Calculate discount
	var discountAmount int64
	if coupon.Type == domain.DiscountTypePercentage {
		discountAmount = (subtotal * coupon.Value) / 100
	} else {
		discountAmount = coupon.Value
	}

	// Apply max discount cap if set
	if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
		discountAmount = coupon.MaxDiscount
	}

	return &ValidateCouponResult{
		Coupon:         coupon,
		DiscountAmount: discountAmount,
	}, nil
}

func (s *CouponService) Create(ctx context.Context, coupon *domain.Coupon) error {
	return s.repo.Create(ctx, coupon)
}

func (s *CouponService) Update(ctx context.Context, id string, fields map[string]any) (*domain.Coupon, error) {
	coupon, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for key, value := range fields {
		switch key {
		case "code":
			if v, ok := value.(string); ok {
				coupon.Code = v
			}
		case "description":
			if v, ok := value.(string); ok {
				coupon.Description = v
			}
		case "type":
			if v, ok := value.(domain.DiscountType); ok {
				coupon.Type = v
			}
		case "value":
			if v, ok := value.(int64); ok {
				coupon.Value = v
			}
		case "min_order_amount":
			if v, ok := value.(int64); ok {
				coupon.MinOrderAmount = v
			}
		case "max_discount":
			if v, ok := value.(int64); ok {
				coupon.MaxDiscount = v
			}
		case "max_usage":
			if v, ok := value.(int); ok {
				coupon.MaxUsage = v
			}
		case "is_active":
			if v, ok := value.(bool); ok {
				coupon.IsActive = v
			}
		case "start_at":
			if value == nil {
				coupon.StartAt = nil
			} else if t, ok := value.(time.Time); ok {
				coupon.StartAt = &t
			}
		case "expires_at":
			if value == nil {
				coupon.ExpiresAt = nil
			} else if t, ok := value.(time.Time); ok {
				coupon.ExpiresAt = &t
			}
		}
	}

	if err := s.repo.Update(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

func (s *CouponService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CouponService) List(ctx context.Context, page, perPage int) ([]*domain.Coupon, int64, error) {
	return s.repo.List(ctx, page, perPage)
}

func (s *CouponService) GetByID(ctx context.Context, id string) (*domain.Coupon, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CouponService) IncrementUsage(ctx context.Context, couponID string) error {
	coupon, err := s.repo.FindByID(ctx, couponID)
	if err != nil {
		return err
	}

	coupon.UsageCount++
	return s.repo.Update(ctx, coupon)
}
