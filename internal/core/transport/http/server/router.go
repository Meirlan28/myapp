package server

import "net/http"

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
}

func NewAPIVersionRouter(
	apiVersion APIVersion,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) AddRoutes(routes ...Route) {
	r.routes = append(r.routes, routes...)
}

func (r *APIVersionRouter) Handlers() map[string]http.Handler {
	handlers := make(map[string]http.Handler, len(r.routes))

	for _, route := range r.routes {
		// Формируем полный паттерн: "GET /api/v1/tasks", "POST /api/v1/users" и т.д.
		pattern := route.Method + " /api/" + string(r.apiVersion) + route.Path
		handler := route.Handler

		handlers[pattern] = handler
	}

	return handlers
}
