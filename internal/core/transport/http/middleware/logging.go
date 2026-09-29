package middleware

import (
	"log/slog"
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/logger"
)

func Logger(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				"request_id", requestID,
				"url", r.URL.String(),
			)

			ctx := logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
