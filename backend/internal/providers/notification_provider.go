package providers

import (
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/events/listeners"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/validation"
)

type NotificationProvider struct{}

func NewNotificationProvider() *NotificationProvider {
	return &NotificationProvider{}
}

func (p *NotificationProvider) Name() string {
	return "kodia:notification"
}

func (p *NotificationProvider) Register(app *kodia.App) error {
	// 1. Repository
	notifRepo := postgres.NewNotificationRepository(app.DB)

	// 2. Get dependencies from container
	broadcaster := kodia.MustResolve[ports.Broadcaster](app, "broadcaster_port")
	dispatcher := kodia.MustResolve[ports.EventDispatcher](app, "event_dispatcher")
	mailer := kodia.MustResolve[ports.Mailer](app, "mailer")
	userRepo := postgres.NewUserRepository(app.DB)

	// 3. Register listener with dependencies
	emailListener := &listeners.SendNotificationEmail{
		Mailer:   mailer,
		UserRepo: userRepo,
	}
	dispatcher.Register("NotificationCreated", emailListener)

	// 4. Service
	notifService := services.NewNotificationService(notifRepo, broadcaster, dispatcher, app.Log)
	app.Set("notification_service", notifService)

	// 5. Handler
	validate := validation.New()
	notifHandler := handlers.NewNotificationHandler(notifService, validate, app.Log)
	app.Set("notification_handler", notifHandler)

	return nil
}

func (p *NotificationProvider) Boot(app *kodia.App) error {
	if app.Router == nil {
		return nil
	}

	notifHandler := kodia.MustResolve[*handlers.NotificationHandler](app, "notification_handler")
	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")

	api := app.Router.Group("/api")
	notifs := api.Group("/notifications")
	notifs.Use(middleware.Auth(jwtManager))
	{
		notifs.GET("", notifHandler.List)
		notifs.GET("/unread-count", notifHandler.UnreadCount)
		notifs.PUT("/:id/read", notifHandler.MarkAsRead)
		notifs.PUT("/read-all", notifHandler.MarkAllAsRead)
		notifs.DELETE("/:id", notifHandler.Delete)
	}

	return nil
}
