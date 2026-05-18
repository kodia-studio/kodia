package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"gorm.io/gorm"
)

type gormOrder struct {
	ID             string         `gorm:"column:id;primaryKey"`
	OrderNumber    string         `gorm:"column:order_number;uniqueIndex"`
	UserID         string         `gorm:"column:user_id;index;not null"`
	Status         string         `gorm:"column:status;not null"`
	Subtotal       int64          `gorm:"column:subtotal;not null"`
	DiscountAmount int64          `gorm:"column:discount_amount;not null"`
	Total          int64          `gorm:"column:total;not null"`
	CouponID       *string        `gorm:"column:coupon_id"`
	CouponCode     string         `gorm:"column:coupon_code"`
	Notes          string         `gorm:"column:notes;type:text"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormOrder) TableName() string { return "orders" }

type gormOrderItem struct {
	ID            string    `gorm:"column:id;primaryKey"`
	OrderID       string    `gorm:"column:order_id;index;not null"`
	ProductID     string    `gorm:"column:product_id;not null"`
	VariantID     *string   `gorm:"column:variant_id"`
	ProductName   string    `gorm:"column:product_name;not null"`
	VariantName   string    `gorm:"column:variant_name"`
	Price         int64     `gorm:"column:price;not null"`
	Quantity      int       `gorm:"column:quantity;not null"`
	DownloadToken string    `gorm:"column:download_token"`
	DownloadCount int       `gorm:"column:download_count;default:0"`
	MaxDownloads  int       `gorm:"column:max_downloads;default:0"`
	FileKey       string    `gorm:"column:file_key"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (gormOrderItem) TableName() string { return "order_items" }

type gormPayment struct {
	ID               string         `gorm:"column:id;primaryKey"`
	OrderID          string         `gorm:"column:order_id;uniqueIndex;not null"`
	MidtransOrderID  string         `gorm:"column:midtrans_order_id;index"`
	SnapToken        string         `gorm:"column:snap_token"`
	SnapURL          string         `gorm:"column:snap_url"`
	Status           string         `gorm:"column:status;not null"`
	PaymentMethod    string         `gorm:"column:payment_method"`
	PaidAt           *time.Time     `gorm:"column:paid_at"`
	MidtransResponse string `gorm:"column:midtrans_response;type:text"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormPayment) TableName() string { return "payments" }

type gormCoupon struct {
	ID              string         `gorm:"column:id;primaryKey"`
	Code            string         `gorm:"column:code;uniqueIndex;not null"`
	Description     string         `gorm:"column:description;type:text"`
	Type            string         `gorm:"column:type;not null"`
	Value           int64          `gorm:"column:value;not null"`
	MinOrderAmount  int64          `gorm:"column:min_order_amount;default:0"`
	MaxDiscount     int64          `gorm:"column:max_discount;default:0"`
	MaxUsage        int            `gorm:"column:max_usage;default:0"`
	UsageCount      int            `gorm:"column:usage_count;default:0"`
	StartAt         *time.Time     `gorm:"column:start_at"`
	ExpiresAt       *time.Time     `gorm:"column:expires_at"`
	IsActive        bool           `gorm:"column:is_active;not null;default:true"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormCoupon) TableName() string { return "coupons" }

func (g *gormOrder) toDomain() *domain.Order {
	return &domain.Order{
		ID:             g.ID,
		OrderNumber:    g.OrderNumber,
		UserID:         g.UserID,
		Status:         domain.OrderStatus(g.Status),
		Subtotal:       g.Subtotal,
		DiscountAmount: g.DiscountAmount,
		Total:          g.Total,
		CouponID:       g.CouponID,
		CouponCode:     g.CouponCode,
		Notes:          g.Notes,
		CreatedAt:      g.CreatedAt,
		UpdatedAt:      g.UpdatedAt,
	}
}

func fromDomainOrder(o *domain.Order) *gormOrder {
	return &gormOrder{
		ID:             o.ID,
		OrderNumber:    o.OrderNumber,
		UserID:         o.UserID,
		Status:         string(o.Status),
		Subtotal:       o.Subtotal,
		DiscountAmount: o.DiscountAmount,
		Total:          o.Total,
		CouponID:       o.CouponID,
		CouponCode:     o.CouponCode,
		Notes:          o.Notes,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}

func (g *gormOrderItem) toDomain() *domain.OrderItem {
	return &domain.OrderItem{
		ID:            g.ID,
		OrderID:       g.OrderID,
		ProductID:     g.ProductID,
		VariantID:     g.VariantID,
		ProductName:   g.ProductName,
		VariantName:   g.VariantName,
		Price:         g.Price,
		Quantity:      g.Quantity,
		DownloadToken: g.DownloadToken,
		DownloadCount: g.DownloadCount,
		MaxDownloads:  g.MaxDownloads,
		FileKey:       g.FileKey,
		CreatedAt:     g.CreatedAt,
	}
}

func fromDomainOrderItem(oi *domain.OrderItem) *gormOrderItem {
	return &gormOrderItem{
		ID:            oi.ID,
		OrderID:       oi.OrderID,
		ProductID:     oi.ProductID,
		VariantID:     oi.VariantID,
		ProductName:   oi.ProductName,
		VariantName:   oi.VariantName,
		Price:         oi.Price,
		Quantity:      oi.Quantity,
		DownloadToken: oi.DownloadToken,
		DownloadCount: oi.DownloadCount,
		MaxDownloads:  oi.MaxDownloads,
		FileKey:       oi.FileKey,
		CreatedAt:     oi.CreatedAt,
	}
}

func (g *gormPayment) toDomain() *domain.Payment {
	var resp string
	if len(g.MidtransResponse) > 0 {
		resp = string(g.MidtransResponse)
	}
	return &domain.Payment{
		ID:               g.ID,
		OrderID:          g.OrderID,
		MidtransOrderID:  g.MidtransOrderID,
		SnapToken:        g.SnapToken,
		SnapURL:          g.SnapURL,
		Status:           domain.PaymentStatus(g.Status),
		PaymentMethod:    g.PaymentMethod,
		PaidAt:           g.PaidAt,
		MidtransResponse: resp,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
	}
}

func fromDomainPayment(p *domain.Payment) *gormPayment {
	gp := &gormPayment{
		ID:              p.ID,
		OrderID:         p.OrderID,
		MidtransOrderID: p.MidtransOrderID,
		SnapToken:       p.SnapToken,
		SnapURL:         p.SnapURL,
		Status:          string(p.Status),
		PaymentMethod:   p.PaymentMethod,
		PaidAt:          p.PaidAt,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
	if p.MidtransResponse != "" {
		gp.MidtransResponse = p.MidtransResponse
	}
	return gp
}

func (g *gormCoupon) toDomain() *domain.Coupon {
	return &domain.Coupon{
		ID:              g.ID,
		Code:            g.Code,
		Description:     g.Description,
		Type:            domain.DiscountType(g.Type),
		Value:           g.Value,
		MinOrderAmount:  g.MinOrderAmount,
		MaxDiscount:     g.MaxDiscount,
		MaxUsage:        g.MaxUsage,
		UsageCount:      g.UsageCount,
		StartAt:         g.StartAt,
		ExpiresAt:       g.ExpiresAt,
		IsActive:        g.IsActive,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
}

func fromDomainCoupon(c *domain.Coupon) *gormCoupon {
	return &gormCoupon{
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
		UpdatedAt:      c.UpdatedAt,
	}
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(fromDomainOrder(order)).Error
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	var g gormOrder
	err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	order := g.toDomain()
	var items []gormOrderItem
	r.db.WithContext(ctx).Find(&items, "order_id = ?", id)
	for i := range items {
		order.Items = append(order.Items, *items[i].toDomain())
	}
	var payment gormPayment
	if err := r.db.WithContext(ctx).First(&payment, "order_id = ?", id).Error; err == nil {
		p := payment.toDomain()
		order.Payment = p
	}
	return order, nil
}

func (r *OrderRepository) FindByUserID(ctx context.Context, userID string, page, perPage int) ([]*domain.Order, int64, error) {
	var gorms []gormOrder
	var total int64
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if err := query.Model(&gormOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&gorms).Error; err != nil {
		return nil, 0, err
	}
	orders := make([]*domain.Order, len(gorms))
	for i, g := range gorms {
		o := g.toDomain()
		var items []gormOrderItem
		r.db.WithContext(ctx).Find(&items, "order_id = ?", g.ID)
		for j := range items {
			o.Items = append(o.Items, *items[j].toDomain())
		}
		orders[i] = o
	}
	return orders, total, nil
}

func (r *OrderRepository) List(ctx context.Context, page, perPage int) ([]*domain.Order, int64, error) {
	var gorms []gormOrder
	var total int64
	if err := r.db.WithContext(ctx).Model(&gormOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(perPage).Find(&gorms).Error; err != nil {
		return nil, 0, err
	}
	orders := make([]*domain.Order, len(gorms))
	for i, g := range gorms {
		o := g.toDomain()
		orders[i] = o
	}
	return orders, total, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Save(fromDomainOrder(order)).Error
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&gormOrder{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *OrderRepository) CreateItem(ctx context.Context, item *domain.OrderItem) error {
	return r.db.WithContext(ctx).Create(fromDomainOrderItem(item)).Error
}

func (r *OrderRepository) FindItemByDownloadToken(ctx context.Context, token string) (*domain.OrderItem, error) {
	var g gormOrderItem
	err := r.db.WithContext(ctx).First(&g, "download_token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.toDomain(), nil
}

func (r *OrderRepository) IncrementDownloadCount(ctx context.Context, itemID string) error {
	return r.db.WithContext(ctx).Model(&gormOrderItem{}).Where("id = ?", itemID).Update("download_count", gorm.Expr("download_count + ?", 1)).Error
}

func (r *OrderRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(fromDomainPayment(payment)).Error
}

func (r *OrderRepository) FindPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var g gormPayment
	err := r.db.WithContext(ctx).First(&g, "order_id = ?", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.toDomain(), nil
}

func (r *OrderRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Save(fromDomainPayment(payment)).Error
}

func (r *OrderRepository) FindPaymentByMidtransOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var g gormPayment
	err := r.db.WithContext(ctx).First(&g, "midtrans_order_id = ?", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.toDomain(), nil
}

type CouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

func (r *CouponRepository) FindByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var g gormCoupon
	err := r.db.WithContext(ctx).First(&g, "UPPER(code) = UPPER(?)", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.toDomain(), nil
}

func (r *CouponRepository) Create(ctx context.Context, coupon *domain.Coupon) error {
	return r.db.WithContext(ctx).Create(fromDomainCoupon(coupon)).Error
}

func (r *CouponRepository) Update(ctx context.Context, coupon *domain.Coupon) error {
	return r.db.WithContext(ctx).Save(fromDomainCoupon(coupon)).Error
}

func (r *CouponRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gormCoupon{}, "id = ?", id).Error
}

func (r *CouponRepository) List(ctx context.Context, page, perPage int) ([]*domain.Coupon, int64, error) {
	var gorms []gormCoupon
	var total int64
	if err := r.db.WithContext(ctx).Model(&gormCoupon{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(perPage).Find(&gorms).Error; err != nil {
		return nil, 0, err
	}
	coupons := make([]*domain.Coupon, len(gorms))
	for i := range gorms {
		coupons[i] = gorms[i].toDomain()
	}
	return coupons, total, nil
}

func (r *CouponRepository) FindByID(ctx context.Context, id string) (*domain.Coupon, error) {
	var g gormCoupon
	err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.toDomain(), nil
}

func (r *CouponRepository) IncrementUsage(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&gormCoupon{}).Where("id = ?", id).Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}
