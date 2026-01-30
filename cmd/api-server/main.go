package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/logeshwarann-dev/news-api-rest/internal/middleware"
	"github.com/logeshwarann-dev/news-api-rest/internal/router"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	r := router.New(nil)
	wrappedRouter := middleware.AddLogger(log, middleware.LogRequest(r))
	log.Info("server running on port 8080")
	if err := http.ListenAndServe(":8080", wrappedRouter); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
