package payment

import (
	"context"
	"fmt"
)

type StripeDriver struct {
	ApiKey string
}

func NewStripeDriver(apiKey string) *StripeDriver {
	return &StripeDriver{ApiKey: apiKey}
}

func (d *StripeDriver) CreateIntent(ctx context.Context, amount int64, currency string, orderID string) (*PaymentIntent, error) {
	// In a real implementation, we would call Stripe API here
	fmt.Printf("[Stripe] Creating intent for Order %s: %d %s\n", orderID, amount, currency)
	
	return &PaymentIntent{
		ID:            "pi_stripe_" + orderID,
		Amount:        amount,
		Currency:      currency,
		Status:        StatusPending,
		CheckoutURL:   "https://checkout.stripe.com/pay/pi_stripe_" + orderID,
		TransactionID: "txn_stripe_" + orderID,
	}, nil
}

func (d *StripeDriver) Refund(ctx context.Context, transactionID string, amount int64) error {
	fmt.Printf("[Stripe] Refunding transaction %s: %d\n", transactionID, amount)
	return nil
}

func (d *StripeDriver) GetStatus(ctx context.Context, transactionID string) (PaymentStatus, error) {
	return StatusSucceeded, nil
}
