package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/logeshwarann-dev/news-api-rest/internal/middleware"
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
	"github.com/logeshwarann-dev/news-api-rest/internal/postgres"
	"github.com/logeshwarann-dev/news-api-rest/internal/router"
	"golang.org/x/sync/errgroup"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	dbConn, err := postgres.NewDB(&postgres.Config{
		DbHost:   os.Getenv("DATABASE_HOST"),
		DbPort:   os.Getenv("DATABASE_PORT"),
		DbName:   os.Getenv("DATABASE_NAME"),
		UserName: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		SSLMode:  "disable",
	})

	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	ns := news.NewStore(dbConn)
	r := router.New(ns)
	wrappedRouter := middleware.AddLogger(log, middleware.LogRequest(r))
	log.Info("server running on port 8080")

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           wrappedRouter,
	}

	errGrp, errGrpCtx := errgroup.WithContext(context.Background())

	errGrp.Go(func() error {
		if err := server.ListenAndServe(); err != nil {
			log.Error("error in running server", "error", err)
			return fmt.Errorf("error in running server: %v", err)
		}
		return nil
	})

	errGrp.Go(func() error {

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		select {
		case sig := <-sigCh:
			log.Info("signal received", "signal", sig)
		case <-errGrpCtx.Done():
		}

		ctx, cancelFunc := context.WithTimeout(errGrpCtx, 5*time.Second)
		defer cancelFunc()

		log.Info("initiating graceful shutdown")

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %v", err)
		}
		return nil
	})

	if err := errGrp.Wait(); err != nil {
		log.Error("failed running errgroup", "error", err)
	}
}
