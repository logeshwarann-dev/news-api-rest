package middleware

import (
	"log/slog"
	"net/http"

	"github.com/logeshwarann-dev/news-api-rest/internal/logger"
)

func AddLogger(slogger *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loggerCtx := logger.CtxWithLogger(r.Context(), slogger)
		r = r.Clone(loggerCtx)
		next.ServeHTTP(w, r)
	}
}

func LogRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slogger := logger.FromContext(r.Context())
		slogger.Info("request", "path", r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
