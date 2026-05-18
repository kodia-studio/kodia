package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusExpired   OrderStatus = "expired"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusExpired PaymentStatus = "expired"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
)

type Order struct {
	ID              string
	OrderNumber     string
	UserID          string
	Status          OrderStatus
	Subtotal        int64
	DiscountAmount  int64
	Total           int64
	CouponID        *string
	CouponCode      string
	Notes           string
	Items           []OrderItem
	Payment         *Payment
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID              string
	OrderID         string
	ProductID       string
	VariantID       *string
	ProductName     string
	VariantName     string
	Price           int64
	Quantity        int
	DownloadToken   string
	DownloadCount   int
	MaxDownloads    int
	FileKey         string
	CreatedAt       time.Time
}

type Payment struct {
	ID                string
	OrderID           string
	MidtransOrderID   string
	SnapToken         string
	SnapURL           string
	Status            PaymentStatus
	PaymentMethod     string
	PaidAt            *time.Time
	MidtransResponse  string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Coupon struct {
	ID               string
	Code             string
	Description      string
	Type             DiscountType
	Value            int64
	MinOrderAmount   int64
	MaxDiscount      int64
	MaxUsage         int
	UsageCount       int
	StartAt          *time.Time
	ExpiresAt        *time.Time
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (c *Coupon) IsValid() bool {
	if !c.IsActive {
		return false
	}
	now := time.Now()
	if c.StartAt != nil && now.Before(*c.StartAt) {
		return false
	}
	if c.ExpiresAt != nil && now.After(*c.ExpiresAt) {
		return false
	}
	if c.MaxUsage > 0 && c.UsageCount >= c.MaxUsage {
		return false
	}
	return true
}

func (c *Coupon) CalculateDiscount(subtotal int64) int64 {
	if c.Type == DiscountTypePercentage {
		discount := (subtotal * c.Value) / 100
		if c.MaxDiscount > 0 && discount > c.MaxDiscount {
			discount = c.MaxDiscount
		}
		return discount
	}
	return c.Value
}
