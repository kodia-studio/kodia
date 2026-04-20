package payment

import (
	"context"
)

// PaymentStatus represents the state of a payment.
type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusSucceeded PaymentStatus = "succeeded"
	StatusFailed    PaymentStatus = "failed"
	StatusRefunded  PaymentStatus = "refunded"
)

// PaymentIntent represents a request to create a payment.
type PaymentIntent struct {
	ID            string        `json:"id"`
	Amount        int64         `json:"amount"` // in smallest unit (e.g., cents)
	Currency      string        `json:"currency"`
	Status        PaymentStatus `json:"status"`
	CheckoutURL   string        `json:"checkout_url,omitempty"`
	TransactionID string        `json:"transaction_id,omitempty"`
}

// PaymentProvider defines the common interface for all payment drivers.
type PaymentProvider interface {
	// CreateIntent creates a new payment session/intent.
	CreateIntent(ctx context.Context, amount int64, currency string, orderID string) (*PaymentIntent, error)

	// Refund processes a refund for a specific transaction.
	Refund(ctx context.Context, transactionID string, amount int64) error

	// GetStatus retrieves the current status of a payment.
	GetStatus(ctx context.Context, transactionID string) (PaymentStatus, error)
}
