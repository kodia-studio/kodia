package payment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"go.uber.org/zap"
)

type PaymentServiceProvider struct{}

func NewServiceProvider() *PaymentServiceProvider {
	return &PaymentServiceProvider{}
}

func (p *PaymentServiceProvider) Name() string {
	return "kodia:payment"
}

func (p *PaymentServiceProvider) Register(app *kodia.App) error {
	// In a real framework, we would read this from app.Config
	// For PoC, let's assume we use 'stripe' by default
	driverName := "stripe" 
	
	var driver PaymentProvider
	switch driverName {
	case "stripe":
		driver = NewStripeDriver("sk_test_12345")
	case "midtrans":
		driver = NewMidtransDriver("mt_server_key_12345")
	default:
		return fmt.Errorf("unsupported payment driver: %s", driverName)
	}

	app.Set("payment", driver)
	app.Log.Info("Payment system initialized", 
		zap.String("driver", driverName),
	)

	return nil
}

func (p *PaymentServiceProvider) Boot(app *kodia.App) error {
	// Automatic Webhook Registration
	if app.Router != nil {
		api := app.Router.Group("/api/payments")
		{
			api.POST("/webhook", func(c *gin.Context) {
				// Handle generic webhook logic
				app.Log.Info("Payment webhook received")
				c.JSON(200, gin.H{"status": "received"})
			})
		}
		app.Log.Info("Payment webhook routes registered at /api/payments/webhook")
	}
	return nil
}
