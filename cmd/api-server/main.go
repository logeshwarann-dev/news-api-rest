package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/logeshwarann-dev/news-api-rest/internal/router"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	r := router.New()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
