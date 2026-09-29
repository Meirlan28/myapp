package server

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
