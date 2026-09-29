package server

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	ApiVersion1 = APIVersion("v1")
	ApiVersion2 = APIVersion("v2")
	ApiVersion3 = APIVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	routes     []Route
	middleware []middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion APIVersion,
	middleware ...middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) AddRoutes(routes ...Route) {
	r.routes = append(r.routes, routes...)
}

func (r *APIVersionRouter) Handlers() map[string]http.Handler {
	handlers := make(map[string]http.Handler, len(r.routes))

	for _, route := range r.routes {
		pattern := route.Method + " /api/" + string(r.apiVersion) + route.Path
		handler := middleware.ChainMiddleware(
			route.WithMiddleware(),
			r.middleware...,
		)

		handlers[pattern] = handler
	}

	return handlers
}
