package ports

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	FindByUserID(ctx context.Context, userID string, page, perPage int) ([]*domain.Order, int64, error)
	List(ctx context.Context, page, perPage int) ([]*domain.Order, int64, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error

	CreateItem(ctx context.Context, item *domain.OrderItem) error
	FindItemByDownloadToken(ctx context.Context, token string) (*domain.OrderItem, error)
	IncrementDownloadCount(ctx context.Context, itemID string) error

	CreatePayment(ctx context.Context, payment *domain.Payment) error
	FindPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
	FindPaymentByMidtransOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
}

type CouponRepository interface {
	FindByCode(ctx context.Context, code string) (*domain.Coupon, error)
	FindByID(ctx context.Context, id string) (*domain.Coupon, error)
	Create(ctx context.Context, coupon *domain.Coupon) error
	Update(ctx context.Context, coupon *domain.Coupon) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, perPage int) ([]*domain.Coupon, int64, error)
	IncrementUsage(ctx context.Context, id string) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []map[string]any, couponCode string) (*domain.Order, error)
	GetMyOrders(ctx context.Context, userID string, page, perPage int) ([]*domain.Order, int64, error)
	GetMyOrder(ctx context.Context, orderID, userID string) (*domain.Order, error)
	ListAllOrders(ctx context.Context, page, perPage int) ([]*domain.Order, int64, error)
	GetOrder(ctx context.Context, id string) (*domain.Order, error)

	InitPayment(ctx context.Context, orderID, userID, customerName, customerEmail string) (*domain.Payment, error)
	HandlePaymentWebhook(ctx context.Context, notification map[string]any) error
	GenerateDownloadURL(ctx context.Context, itemID, userID string) (string, error)
}

type CouponService interface {
	Validate(ctx context.Context, code string, subtotal int64) (*domain.Coupon, int64, error)
	List(ctx context.Context, page, perPage int) ([]*domain.Coupon, int64, error)
	Create(ctx context.Context, coupon *domain.Coupon) error
	Update(ctx context.Context, id string, fields map[string]any) (*domain.Coupon, error)
	Delete(ctx context.Context, id string) error
}
