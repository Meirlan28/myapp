package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type HTTPServer struct {
	cfg    Config
	mux    *http.ServeMux
	server *http.Server
	Logger *slog.Logger
}

func New(cfg Config, logger *slog.Logger) *HTTPServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &HTTPServer{
		cfg,
		mux,
		server,
		logger,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		handlers := router.Handlers()

		for path, handler := range handlers {
			s.mux.Handle(path, handler)
		}
	}
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		path := route.Method + " " + route.Path
		handler := route.Handler

		s.mux.Handle(path, handler)
	}
}

func (s *HTTPServer) Start() {
	s.Logger.Info(
		"server started",
		"port", s.cfg.Addr,
	)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
