package providers

import (
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/midtrans"
)

// CommerceProvider wires all marketplace functionality (products, orders, coupons, payments).
type CommerceProvider struct{}

func NewCommerceProvider() *CommerceProvider {
	return &CommerceProvider{}
}

func (p *CommerceProvider) Name() string {
	return "kodia:commerce"
}

func (p *CommerceProvider) Register(app *kodia.App) error {
	// Midtrans client
	midtransClient := midtrans.New(
		app.Config.Midtrans.ServerKey,
		app.Config.Midtrans.ClientKey,
		app.Config.Midtrans.IsProduction,
	)

	// Wrap the infra storage provider to satisfy services.StorageProvider
	rawStorage, ok := app.Get("storage")
	var storageAdapter services.StorageProvider
	if ok && rawStorage != nil {
		if sp, ok := rawStorage.(StorageProvider); ok {
			storageAdapter = &storageProviderAdapter{sp}
		}
	}
	if storageAdapter == nil {
		storageAdapter = &noopStorageProvider{}
	}

	// Repositories
	productRepo := postgres.NewProductRepository(app.DB)
	couponRepo := postgres.NewCouponRepository(app.DB)
	orderRepo := postgres.NewOrderRepository(app.DB)
	productDocRepo := postgres.NewProductDocRepository(app.DB)

	// Services
	productSvc := services.NewProductService(productRepo, storageAdapter)
	couponSvc := services.NewCouponService(couponRepo)
	orderSvc := services.NewOrderService(orderRepo, couponRepo, productRepo, couponSvc, midtransClient, storageAdapter)
	productDocSvc := services.NewProductDocService(productDocRepo)

	// Handlers
	productHandler := handlers.NewProductHandler(productSvc)
	couponHandler := handlers.NewCouponHandler(couponSvc)
	orderHandler := handlers.NewOrderHandler(orderSvc)
	paymentHandler := handlers.NewPaymentHandler(orderSvc)
	productDocHandler := handlers.NewProductDocHandler(productDocSvc, productSvc)

	app.Set("product_handler", productHandler)
	app.Set("coupon_handler", couponHandler)
	app.Set("order_handler", orderHandler)
	app.Set("payment_handler", paymentHandler)
	app.Set("product_doc_handler", productDocHandler)

	return nil
}

func (p *CommerceProvider) Boot(app *kodia.App) error {
	if app.Router == nil {
		return nil
	}

	productHandler := kodia.MustResolve[*handlers.ProductHandler](app, "product_handler")
	couponHandler := kodia.MustResolve[*handlers.CouponHandler](app, "coupon_handler")
	orderHandler := kodia.MustResolve[*handlers.OrderHandler](app, "order_handler")
	paymentHandler := kodia.MustResolve[*handlers.PaymentHandler](app, "payment_handler")
	productDocHandler := kodia.MustResolve[*handlers.ProductDocHandler](app, "product_doc_handler")
	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

	api := app.Router.Group("/api")

	// Public product routes
	api.GET("/products", productHandler.ListPublic)
	api.GET("/products/:slug", productHandler.GetBySlug)

	// Public product documentation routes
	api.GET("/products/:slug/docs", productDocHandler.ListDocs)
	api.GET("/products/:slug/docs/:section/:docSlug", productDocHandler.GetDoc)
	api.GET("/products/:slug/guide", productDocHandler.ListDocs)
	api.GET("/products/:slug/guide/:section/:docSlug", productDocHandler.GetDoc)

	// Public payment webhook (no auth — verified via Midtrans signature)
	api.POST("/payments/webhook", paymentHandler.HandleWebhook)

	// Authenticated user routes
	auth := api.Group("")
	auth.Use(middleware.Auth(jwtManager))
	{
		registerOrderRoutes(auth.Group("/orders"), orderHandler)
		auth.POST("/coupons/validate", couponHandler.ValidateCoupon)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.Auth(jwtManager))
	admin.Use(middleware.RequireRole("admin"))
	{
		productsGroup := admin.Group("/products")
		registerAdminProductRoutes(productsGroup, productHandler)
		registerAdminProductDocRoutes(productsGroup, productDocHandler)
		registerAdminCouponRoutes(admin.Group("/coupons"), couponHandler)
		registerAdminOrderRoutes(admin.Group("/orders"), orderHandler)
		registerAnalyticsRoutes(admin.Group("/analytics"), orderHandler)
	}

	return nil
}

func registerOrderRoutes(g *gin.RouterGroup, h *handlers.OrderHandler) {
	g.POST("", h.CreateOrder)
	g.GET("", h.GetMyOrders)
	g.GET("/:id", h.GetMyOrder)
	g.POST("/:id/pay", h.InitPayment)
	g.GET("/:id/items/:itemId/download", h.GenerateDownloadURL)
}

func registerAdminProductRoutes(g *gin.RouterGroup, h *handlers.ProductHandler) {
	g.GET("", h.ListAdmin)
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.PATCH("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/cover", h.UploadCover)
	g.POST("/:id/variants", h.AddVariant)
	g.PATCH("/:id/variants/:variantId", h.UpdateVariant)
	g.DELETE("/:id/variants/:variantId", h.DeleteVariant)
	g.POST("/:id/variants/:variantId/upload", h.UploadVariantFile)
}

func registerAdminCouponRoutes(g *gin.RouterGroup, h *handlers.CouponHandler) {
	g.GET("", h.ListCoupons)
	g.POST("", h.CreateCoupon)
	g.PATCH("/:id", h.UpdateCoupon)
	g.DELETE("/:id", h.DeleteCoupon)
}

func registerAdminOrderRoutes(g *gin.RouterGroup, h *handlers.OrderHandler) {
	g.GET("", h.GetAdminOrders)
	g.GET("/:id", h.GetAdminOrder)
}

func registerAnalyticsRoutes(g *gin.RouterGroup, h *handlers.OrderHandler) {
	g.GET("/dashboard", h.DashboardAnalytics)
}

func registerAdminProductDocRoutes(g *gin.RouterGroup, h *handlers.ProductDocHandler) {
	g.GET("/:id/documentation", h.ListAdminDocs)
	g.POST("/:id/documentation", h.CreateDoc)
	g.PUT("/:id/documentation/:docId", h.UpdateDoc)
	g.DELETE("/:id/documentation/:docId", h.DeleteDoc)
}

// storageProviderAdapter bridges providers.StorageProvider → services.StorageProvider.
type storageProviderAdapter struct {
	inner StorageProvider
}

func (a *storageProviderAdapter) UploadStream(ctx context.Context, key string, src io.Reader, size int64) error {
	_, err := a.inner.Store(ctx, &FileUpload{
		Content: src,
		Path:    key,
		Size:    size,
	})
	return err
}

func (a *storageProviderAdapter) GetURL(ctx context.Context, key string) (string, error) {
	return a.inner.GetURL(ctx, key)
}

func (a *storageProviderAdapter) GetTemporaryURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return a.inner.GetTemporaryURL(ctx, key, expiration)
}

// noopStorageProvider is used when no storage backend is configured.
type noopStorageProvider struct{}

func (n *noopStorageProvider) UploadStream(_ context.Context, _ string, _ io.Reader, _ int64) error {
	return nil
}
func (n *noopStorageProvider) GetURL(_ context.Context, key string) (string, error) {
	return key, nil
}
func (n *noopStorageProvider) GetTemporaryURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return key, nil
}
