package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/transport/http/middleware"
)

type Server struct {
	cfg        Config
	mux        *http.ServeMux
	Logger     *slog.Logger
	middleware []middleware.Middleware
}

func New(cfg Config, logger *slog.Logger, middleware ...middleware.Middleware) *Server {
	mux := http.NewServeMux()

	return &Server{
		cfg,
		mux,
		logger,
		middleware,
	}
}

func (s *Server) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		handlers := router.Handlers()

		for path, handler := range handlers {
			s.mux.Handle(path, handler)
		}
	}
}

func (s *Server) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		path := route.Method + " " + route.Path
		handler := route.WithMiddleware()

		s.mux.Handle(path, handler)
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      mux,
		IdleTimeout:  s.cfg.IdleTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	// Буферизированный канал (размер 1), чтобы горутина не заблокировалась
	// при отправке ошибки, если основной поток уже ушёл в select.
	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.Logger.Info("start HTTP server", slog.String("addr", s.cfg.Addr))

		err := server.ListenAndServe()

		// userhttp.ErrServerClosed — нормальное завершение после Shutdown(), не ошибка.
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		// HTTP сервер завершился из-за ошибки (например, порт занят).
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		// Получен сигнал завершения — выполняем graceful shutdown.
		s.Logger.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			// Если graceful shutdown не успел — принудительно закрываем.
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.Logger.Warn("HTTP server stopped")
	}

	return nil
}
