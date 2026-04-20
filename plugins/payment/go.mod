module github.com/kodia-studio/payment

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/kodia-studio/kodia v0.0.0
	go.uber.org/zap v1.27.1
)

replace github.com/kodia-studio/kodia => ../../backend
