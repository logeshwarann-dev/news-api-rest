package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/logeshwarann-dev/news-api-rest/internal/middleware"
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
	"github.com/logeshwarann-dev/news-api-rest/internal/postgres"
	"github.com/logeshwarann-dev/news-api-rest/internal/router"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	dbConn, err := postgres.NewDB(&postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	ns := news.NewStore(dbConn)
	r := router.New(ns)
	wrappedRouter := middleware.AddLogger(log, middleware.LogRequest(r))
	log.Info("server running on port 8080")
	if err := http.ListenAndServe(":8080", wrappedRouter); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
