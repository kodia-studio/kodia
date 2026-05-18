package midtrans

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Client struct {
	snapClient *snap.Client
	serverKey  string
}

type Customer struct {
	Name  string
	Email string
	Phone string
}

type CreateTransactionRequest struct {
	OrderID  string
	Amount   int64
	Customer Customer
	Items    []ItemDetail
}

type ItemDetail struct {
	ID       string
	Name     string
	Price    int64
	Quantity int32
}

func New(serverKey, clientKey string, isProduction bool) *Client {
	snapClient := &snap.Client{}

	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	snapClient.New(serverKey, env)

	return &Client{
		snapClient: snapClient,
		serverKey:  serverKey,
	}
}

func (c *Client) CreateTransaction(ctx context.Context, order interface{}, customer Customer) (string, string, error) {
	// This is a placeholder implementation
	// In production, convert order to proper request format
	return "", "", fmt.Errorf("not implemented")
}

func (c *Client) VerifyWebhookSignature(orderID string, statusCode, grossAmount, signatureKey string) bool {
	hash := sha512.Sum512([]byte(orderID + statusCode + grossAmount + c.serverKey))
	expectedSignature := hex.EncodeToString(hash[:])
	return expectedSignature == signatureKey
}
