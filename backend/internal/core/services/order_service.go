package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/midtrans"
)

type OrderService struct {
	orderRepo       ports.OrderRepository
	couponRepo      ports.CouponRepository
	productRepo     ports.ProductRepository
	couponService   *CouponService
	midtransClient  *midtrans.Client
	storageProvider StorageProvider
}

func NewOrderService(
	orderRepo ports.OrderRepository,
	couponRepo ports.CouponRepository,
	productRepo ports.ProductRepository,
	couponService *CouponService,
	midtransClient *midtrans.Client,
	storageProvider StorageProvider,
) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		couponRepo:      couponRepo,
		productRepo:     productRepo,
		couponService:   couponService,
		midtransClient:  midtransClient,
		storageProvider: storageProvider,
	}
}

type CreateOrderRequest struct {
	UserID string
	Items  []struct {
		ProductID string `json:"product_id"`
		VariantID string `json:"variant_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
	CouponCode string `json:"coupon_code"`
}

func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	order := &domain.Order{
		ID:         generateOrderID(),
		OrderNumber: generateOrderNumber(),
		UserID:     req.UserID,
		Status:     domain.OrderStatusPending,
		Items:      []domain.OrderItem{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	var subtotal int64

	// Process each item
	for _, itemReq := range req.Items {
		product, err := s.productRepo.FindByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", itemReq.ProductID)
		}

		var variant *domain.ProductVariant
		for i := range product.Variants {
			if product.Variants[i].ID == itemReq.VariantID {
				variant = &product.Variants[i]
				break
			}
		}

		if variant == nil {
			return nil, fmt.Errorf("variant %s not found", itemReq.VariantID)
		}

		if !variant.IsActive {
			return nil, fmt.Errorf("variant %s is not available", itemReq.VariantID)
		}

		// Use effective price (promo if active, else regular)
		price := variant.EffectivePrice()

		orderItem := domain.OrderItem{
			ID:            generateItemID(),
			OrderID:       order.ID,
			ProductID:     itemReq.ProductID,
			VariantID:     &itemReq.VariantID,
			ProductName:   product.Name,
			VariantName:   variant.Name,
			Price:         price,
			Quantity:      itemReq.Quantity,
			FileKey:       variant.FileKey,
			MaxDownloads:  5,
			DownloadCount: 0,
			CreatedAt:     time.Now(),
		}

		order.Items = append(order.Items, orderItem)
		subtotal += price * int64(itemReq.Quantity)
	}

	order.Subtotal = subtotal
	order.Total = subtotal

	// Apply coupon if provided
	if req.CouponCode != "" {
		coupon, err := s.couponRepo.FindByCode(ctx, req.CouponCode)
		if err == nil && coupon != nil {
			// Validate coupon
			result, err := s.couponService.Validate(ctx, req.CouponCode, subtotal)
			if err == nil && result != nil {
				order.CouponCode = req.CouponCode
				order.DiscountAmount = result.DiscountAmount
				order.Total = subtotal - result.DiscountAmount
			}
		}
	}

	// Save order to database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetMyOrders(ctx context.Context, userID string, page, perPage int) ([]*domain.Order, int64, error) {
	return s.orderRepo.FindByUserID(ctx, userID, page, perPage)
}

func (s *OrderService) GetMyOrder(ctx context.Context, userID, orderID string) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

type InitPaymentRequest struct {
	OrderID      string `json:"order_id"`
	CustomerName string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type InitPaymentResponse struct {
	SnapToken string `json:"snap_token"`
	SnapURL   string `json:"snap_url"`
}

func (s *OrderService) InitPayment(ctx context.Context, req *InitPaymentRequest) (*InitPaymentResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	if order.Status != domain.OrderStatusPending {
		return nil, fmt.Errorf("order is not pending")
	}

	// Create Midtrans transaction
	snapToken, snapURL, err := s.midtransClient.CreateTransaction(
		ctx,
		order,
		midtrans.Customer{
			Name:  req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Create and save payment record
	payment := &domain.Payment{
		ID:              generatePaymentID(),
		OrderID:         order.ID,
		MidtransOrderID: order.OrderNumber,
		SnapToken:       snapToken,
		SnapURL:         snapURL,
		Status:          domain.PaymentStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.orderRepo.UpdatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return &InitPaymentResponse{
		SnapToken: snapToken,
		SnapURL:   snapURL,
	}, nil
}

func (s *OrderService) HandlePaymentWebhook(ctx context.Context, notification map[string]interface{}) error {
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return fmt.Errorf("invalid order_id in notification")
	}

	transactionStatus, ok := notification["transaction_status"].(string)
	if !ok {
		return fmt.Errorf("invalid transaction_status in notification")
	}

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	payment, err := s.orderRepo.FindPaymentByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("payment not found")
	}

	// Update payment status and response
	payment.MidtransResponse = fmt.Sprintf("%v", notification)
	payment.UpdatedAt = time.Now()

	switch transactionStatus {
	case "capture", "settlement":
		payment.Status = domain.PaymentStatusPaid
		payment.PaidAt = &payment.UpdatedAt
		order.Status = domain.OrderStatusPaid

		// Generate download tokens for each item
		for i := range order.Items {
			order.Items[i].DownloadToken = generateDownloadToken()
		}

	case "pending":
		payment.Status = domain.PaymentStatusPending

	case "deny", "cancel", "expire":
		payment.Status = domain.PaymentStatusFailed
		order.Status = domain.OrderStatusCancelled

	default:
		return fmt.Errorf("unknown transaction status: %s", transactionStatus)
	}

	// Update order and payment
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	if err := s.orderRepo.UpdatePayment(ctx, payment); err != nil {
		return err
	}

	// Increment coupon usage if applied
	if order.CouponCode != "" {
		_ = s.couponService.IncrementUsage(ctx, order.CouponCode)
	}

	return nil
}

type DownloadInfo struct {
	OrderItemID string
	FileKey     string
	DownloadURL string
}

func (s *OrderService) GenerateDownloadURL(ctx context.Context, userID, orderID, itemID string) (*DownloadInfo, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	if order.Status != domain.OrderStatusPaid {
		return nil, fmt.Errorf("order is not paid")
	}

	var orderItem *domain.OrderItem
	for i := range order.Items {
		if order.Items[i].ID == itemID {
			orderItem = &order.Items[i]
			break
		}
	}

	if orderItem == nil {
		return nil, fmt.Errorf("item not found in order")
	}

	if orderItem.FileKey == "" {
		return nil, fmt.Errorf("item has no downloadable file")
	}

	if orderItem.MaxDownloads > 0 && orderItem.DownloadCount >= orderItem.MaxDownloads {
		return nil, fmt.Errorf("download limit exceeded")
	}

	// Get temporary download URL from storage
	downloadURL, err := s.storageProvider.GetTemporaryURL(ctx, orderItem.FileKey, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download url: %w", err)
	}

	// Increment download count
	orderItem.DownloadCount++
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return &DownloadInfo{
		OrderItemID: orderItem.ID,
		FileKey:     orderItem.FileKey,
		DownloadURL: downloadURL,
	}, nil
}

func (s *OrderService) GetAdminOrders(ctx context.Context, page, perPage int) ([]*domain.Order, int64, error) {
	return s.orderRepo.List(ctx, page, perPage)
}

func (s *OrderService) GetAdminOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

// Helper functions for generating IDs and tokens
func generateOrderID() string {
	return "ORD-" + generateRandomString(12)
}

func generateOrderNumber() string {
	return "KD-" + time.Now().Format("20060102") + "-" + generateRandomString(6)
}

func generateItemID() string {
	return "ITEM-" + generateRandomString(12)
}

func generatePaymentID() string {
	return "PAY-" + generateRandomString(12)
}

func generateDownloadToken() string {
	return generateRandomString(32)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}
