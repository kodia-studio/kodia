module github.com/kodia-studio/pulse

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/kodia-studio/kodia v0.0.0
	go.uber.org/zap v1.27.1
	github.com/gorilla/websocket v1.5.0
)

replace github.com/kodia-studio/kodia => ../../backend
