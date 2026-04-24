package providers

import (
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/sse"
	ws "github.com/kodia-studio/kodia/internal/adapters/websocket"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/internal/infrastructure/broadcasting"
	notifinfra "github.com/kodia-studio/kodia/internal/infrastructure/notification"
	"github.com/kodia-studio/kodia/internal/infrastructure/notification/channels"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
)

// RealtimeProvider wires SSE, EventBroadcaster, and NotificationManager into the app.
type RealtimeProvider struct{}

func NewRealtimeProvider() *RealtimeProvider { return &RealtimeProvider{} }

func (p *RealtimeProvider) Name() string { return "kodia:realtime" }

func (p *RealtimeProvider) Register(app *kodia.App) error {
	cfg := app.Config

	// 1. SSE Manager
	sseManager := sse.NewManager()
	app.Set("sse_manager", sseManager)

	// 2. SSE Handler
	sseHandler := sse.NewHandler(sseManager, app.Log)
	app.Set("sse_handler", sseHandler)

	// 3. Event Broadcaster (WS Hub + SSE)
	hub := kodia.MustResolve[*ws.Hub](app, "ws_hub")
	eventBroadcaster := broadcasting.NewEventBroadcaster(hub, sseManager, app.Log)
	app.Set("event_broadcaster", eventBroadcaster)
	app.Set("broadcaster_port", ports.Broadcaster(eventBroadcaster))

	// 4. Notification Manager with all registered channel drivers
	notifManager := notifinfra.NewManager(app.Log)

	// Email — always registered, wraps existing Mailer port
	mailer := kodia.MustResolve[ports.Mailer](app, "mailer")
	notifManager.Register(channels.NewEmailChannel(mailer))

	// WebSocket in-app channel — always registered
	notifManager.Register(channels.NewWebSocketChannel(hub))

	// SMS via Twilio — optional, only when credentials are configured
	if cfg.Notification.TwilioAccountSID != "" {
		notifManager.Register(channels.NewSMSChannel(
			cfg.Notification.TwilioAccountSID,
			cfg.Notification.TwilioAuthToken,
			cfg.Notification.TwilioFromNumber,
		))
	}

	// Slack via Incoming Webhook — optional
	if cfg.Notification.SlackWebhookURL != "" {
		notifManager.Register(channels.NewSlackChannel(cfg.Notification.SlackWebhookURL))
	}

	// Push via Firebase FCM — optional
	if cfg.Notification.FCMServerKey != "" {
		notifManager.Register(channels.NewPushChannel(cfg.Notification.FCMServerKey))
	}

	app.Set("notification_manager", notifManager)
	app.Set("notification_manager_port", ports.NotificationManager(notifManager))

	return nil
}

func (p *RealtimeProvider) Boot(app *kodia.App) error {
	if app.Router == nil {
		return nil
	}

	sseHandler := kodia.MustResolve[*sse.Handler](app, "sse_handler")
	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

	sseGroup := app.Router.Group("/api/v1/sse")
	{
		// Public channel — no auth required
		sseGroup.GET("/:channel", sseHandler.ServePublic)
		// Private user stream — requires valid JWT
		sseGroup.GET("/user", middleware.Auth(jwtManager), sseHandler.ServeUser)
		// Metrics
		sseGroup.GET("/status", sseHandler.Status)
	}

	return nil
}
