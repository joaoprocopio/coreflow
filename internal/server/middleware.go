package server

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(h http.Handler) http.Handler

func loggerMiddleware(h http.Handler, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Info(
			"request received",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.String("ip", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.Duration("duration", duration),
		)
	}
}
