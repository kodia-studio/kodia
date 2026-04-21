module github.com/kodia-studio/authsocial

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/kodia-studio/kodia v0.0.0
	go.uber.org/zap v1.27.1
	github.com/google/uuid v1.6.0
	golang.org/x/oauth2 v0.36.0
)

replace github.com/kodia-studio/kodia => ../../backend
