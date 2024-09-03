package middleware

import (
	"myproject/pkg/config"
	"myproject/pkg/utils"
	"net/http"
)

type Middleware interface {
	AuthMiddleware(http.Handler) http.Handler
	LoggingMiddleware(http.Handler) http.Handler
}

type middleware struct {
	config   *config.Application
	jwtUtils *utils.JwtUtils
}

func NewMiddleware(appConfig *config.Application, jwtService *utils.JwtUtils) Middleware {
	return &middleware{
		config:   appConfig,
		jwtUtils: jwtService,
	}
}
