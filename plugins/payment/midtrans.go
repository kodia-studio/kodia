package payment

import (
	"context"
	"fmt"
)

type MidtransDriver struct {
	ServerKey string
}

func NewMidtransDriver(serverKey string) *MidtransDriver {
	return &MidtransDriver{ServerKey: serverKey}
}

func (d *MidtransDriver) CreateIntent(ctx context.Context, amount int64, currency string, orderID string) (*PaymentIntent, error) {
	fmt.Printf("[Midtrans] Creating Snap Token for Order %s: %d %s\n", orderID, amount, currency)
	
	return &PaymentIntent{
		ID:            "mt_midtrans_" + orderID,
		Amount:        amount,
		Currency:      currency,
		Status:        StatusPending,
		CheckoutURL:   "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + orderID,
		TransactionID: "txn_midtrans_" + orderID,
	}, nil
}

func (d *MidtransDriver) Refund(ctx context.Context, transactionID string, amount int64) error {
	fmt.Printf("[Midtrans] Refunding through Midtrans: %s (%d)\n", transactionID, amount)
	return nil
}

func (d *MidtransDriver) GetStatus(ctx context.Context, transactionID string) (PaymentStatus, error) {
	return StatusPending, nil
}
