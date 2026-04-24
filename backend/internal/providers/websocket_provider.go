package providers

import (
	"github.com/kodia-studio/kodia/internal/adapters/websocket"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/kodia"
)

type WebSocketProvider struct{}

func NewWebSocketProvider() *WebSocketProvider {
	return &WebSocketProvider{}
}

func (p *WebSocketProvider) Name() string {
	return "kodia:websocket"
}

func (p *WebSocketProvider) Register(app *kodia.App) error {
	hub := websocket.NewHub()
	go hub.Run()
	app.Set("ws_hub", hub)

	jwtManager := kodia.MustResolve[*jwt.Manager](app, "jwt_manager")
	wsHandler := websocket.NewHandler(hub, jwtManager, app.Log)
	app.Set("ws_handler", wsHandler)

	return nil
}

func (p *WebSocketProvider) Boot(app *kodia.App) error {
	if app.Router != nil {
		wsHandler := kodia.MustResolve[*websocket.Handler](app, "ws_handler")
		api := app.Router.Group("/api/ws")
		{
			api.GET("", wsHandler.ServeWS)
			api.GET("/room/:room", wsHandler.ServeRoom)
			api.GET("/status", wsHandler.GetStatus)
		}
	}
	return nil
}
