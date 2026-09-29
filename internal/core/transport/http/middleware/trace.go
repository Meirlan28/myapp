package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Meirlan28/myapp/internal/core/logger"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

// до этого не доебывайся это будет расширятся
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			rw := response.NewResponseWriter(w)

			before := time.Now()

			log.DebugContext(
				ctx,
				">>> incoming HTTP request",
				slog.String("http_method", r.Method),
				slog.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.DebugContext(
				ctx,
				"<<< done HTTP request",
				slog.Int("status_code", rw.StatusCode()),
				slog.Duration("latency", time.Since(before)),
			)
		})
	}
}
