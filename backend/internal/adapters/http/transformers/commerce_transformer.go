package transformers

import (
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
)

// ProductResponse is the JSON representation of a Product.
type ProductResponse struct {
	ID              string             `json:"id"`
	Slug            string             `json:"slug"`
	Name            string             `json:"name"`
	Tagline         string             `json:"tagline"`
	Description     string             `json:"description"`
	Type            string             `json:"type"`
	ServiceType     *string            `json:"service_type,omitempty"`
	CoverURL        string             `json:"cover_url"`
	Tags            []string           `json:"tags"`
	IsPublished     bool               `json:"is_published"`
	SortOrder       int                `json:"sort_order"`
	MetaTitle       string             `json:"meta_title,omitempty"`
	MetaDescription string             `json:"meta_description,omitempty"`
	OGImageURL      string             `json:"og_image_url,omitempty"`
	Variants        []VariantResponse  `json:"variants"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type VariantResponse struct {
	ID           string     `json:"id"`
	ProductID    string     `json:"product_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Price        int64      `json:"price"`
	PromoPrice   *int64     `json:"promo_price,omitempty"`
	PromoStartAt *time.Time `json:"promo_start_at,omitempty"`
	PromoEndAt   *time.Time `json:"promo_end_at,omitempty"`
	FileKey      string     `json:"file_key,omitempty"`
	IsActive     bool       `json:"is_active"`
	SortOrder    int        `json:"sort_order"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func ProductToResponse(p *domain.Product) *ProductResponse {
	if p == nil {
		return nil
	}
	variants := make([]VariantResponse, 0, len(p.Variants))
	for _, v := range p.Variants {
		variants = append(variants, variantToResponse(&v))
	}
	var st *string
	if p.ServiceType != nil {
		s := string(*p.ServiceType)
		st = &s
	}
	return &ProductResponse{
		ID:              p.ID,
		Slug:            p.Slug,
		Name:            p.Name,
		Tagline:         p.Tagline,
		Description:     p.Description,
		Type:            string(p.Type),
		ServiceType:     st,
		CoverURL:        p.CoverURL,
		Tags:            p.Tags,
		IsPublished:     p.IsPublished,
		SortOrder:       p.SortOrder,
		MetaTitle:       p.MetaTitle,
		MetaDescription: p.MetaDescription,
		OGImageURL:      p.OGImageURL,
		Variants:        variants,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func ProductsToResponse(products []*domain.Product) []*ProductResponse {
	result := make([]*ProductResponse, 0, len(products))
	for _, p := range products {
		result = append(result, ProductToResponse(p))
	}
	return result
}

func variantToResponse(v *domain.ProductVariant) VariantResponse {
	return VariantResponse{
		ID:           v.ID,
		ProductID:    v.ProductID,
		Name:         v.Name,
		Description:  v.Description,
		Price:        v.Price,
		PromoPrice:   v.PromoPrice,
		PromoStartAt: v.PromoStartAt,
		PromoEndAt:   v.PromoEndAt,
		FileKey:      v.FileKey,
		IsActive:     v.IsActive,
		SortOrder:    v.SortOrder,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

func VariantToResponse(v *domain.ProductVariant) *VariantResponse {
	r := variantToResponse(v)
	return &r
}

// CouponResponse is the JSON representation of a Coupon.
type CouponResponse struct {
	ID             string     `json:"id"`
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	Type           string     `json:"type"`
	Value          int64      `json:"value"`
	MinOrderAmount int64      `json:"min_order_amount"`
	MaxDiscount    int64      `json:"max_discount"`
	MaxUsage       int        `json:"max_usage"`
	UsageCount     int        `json:"usage_count"`
	StartAt        *time.Time `json:"start_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
}

func CouponToResponse(c *domain.Coupon) *CouponResponse {
	if c == nil {
		return nil
	}
	return &CouponResponse{
		ID:             c.ID,
		Code:           c.Code,
		Description:    c.Description,
		Type:           string(c.Type),
		Value:          c.Value,
		MinOrderAmount: c.MinOrderAmount,
		MaxDiscount:    c.MaxDiscount,
		MaxUsage:       c.MaxUsage,
		UsageCount:     c.UsageCount,
		StartAt:        c.StartAt,
		ExpiresAt:      c.ExpiresAt,
		IsActive:       c.IsActive,
		CreatedAt:      c.CreatedAt,
	}
}

func CouponsToResponse(coupons []*domain.Coupon) []*CouponResponse {
	result := make([]*CouponResponse, 0, len(coupons))
	for _, c := range coupons {
		result = append(result, CouponToResponse(c))
	}
	return result
}

// OrderResponse is the JSON representation of an Order.
type OrderResponse struct {
	ID             string          `json:"id"`
	OrderNumber    string          `json:"order_number"`
	UserID         string          `json:"user_id"`
	Status         string          `json:"status"`
	Subtotal       int64           `json:"subtotal"`
	DiscountAmount int64           `json:"discount_amount"`
	Total          int64           `json:"total"`
	CouponCode     string          `json:"coupon_code,omitempty"`
	Items          []OrderItemResp `json:"items"`
	Payment        *PaymentResp    `json:"payment,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type OrderItemResp struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	VariantID     *string   `json:"variant_id,omitempty"`
	ProductName   string    `json:"product_name"`
	VariantName   string    `json:"variant_name"`
	Price         int64     `json:"price"`
	Quantity      int       `json:"quantity"`
	DownloadToken string    `json:"download_token,omitempty"`
	DownloadCount int       `json:"download_count"`
	MaxDownloads  int       `json:"max_downloads"`
	CreatedAt     time.Time `json:"created_at"`
}

type PaymentResp struct {
	ID            string     `json:"id"`
	SnapToken     string     `json:"snap_token,omitempty"`
	SnapURL       string     `json:"snap_url,omitempty"`
	Status        string     `json:"status"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func OrderToResponse(o *domain.Order) *OrderResponse {
	if o == nil {
		return nil
	}
	items := make([]OrderItemResp, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, OrderItemResp{
			ID:            item.ID,
			ProductID:     item.ProductID,
			VariantID:     item.VariantID,
			ProductName:   item.ProductName,
			VariantName:   item.VariantName,
			Price:         item.Price,
			Quantity:      item.Quantity,
			DownloadToken: item.DownloadToken,
			DownloadCount: item.DownloadCount,
			MaxDownloads:  item.MaxDownloads,
			CreatedAt:     item.CreatedAt,
		})
	}
	var payment *PaymentResp
	if o.Payment != nil {
		payment = &PaymentResp{
			ID:            o.Payment.ID,
			SnapToken:     o.Payment.SnapToken,
			SnapURL:       o.Payment.SnapURL,
			Status:        string(o.Payment.Status),
			PaymentMethod: o.Payment.PaymentMethod,
			PaidAt:        o.Payment.PaidAt,
			CreatedAt:     o.Payment.CreatedAt,
		}
	}
	return &OrderResponse{
		ID:             o.ID,
		OrderNumber:    o.OrderNumber,
		UserID:         o.UserID,
		Status:         string(o.Status),
		Subtotal:       o.Subtotal,
		DiscountAmount: o.DiscountAmount,
		Total:          o.Total,
		CouponCode:     o.CouponCode,
		Items:          items,
		Payment:        payment,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}

func OrdersToResponse(orders []*domain.Order) []*OrderResponse {
	result := make([]*OrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, OrderToResponse(o))
	}
	return result
}
