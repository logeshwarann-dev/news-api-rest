package news_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/logeshwarann-dev/news-api-rest/internal/postgres"
	"github.com/testcontainers/testcontainers-go"
	pgtc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
)

func createTestContainer(ctx context.Context) (ctr *pgtc.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("working dir: %w", err)
	}
	sqlScripts := wd + "/testdata/sql/store.sql"

	ctr, err = pgtc.Run(
		ctx,
		"postgres:16-apline",
		pgtc.WithInitScripts(sqlScripts),
		pgtc.WithDatabase("postgres"),
		pgtc.WithUsername("postgres"),
		pgtc.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return ctr, fmt.Errorf("run container: %w", err)
	}
	return ctr, nil
}

type DBCleanup func(ctx context.Context) error

func createTestDB(ctx context.Context) (*bun.DB, DBCleanup, error) {
	pgCtr, err := createTestContainer(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("create test container: %w", err)
	}

	port, err := pgCtr.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("map port: %w", err)
	}

	db, err := postgres.NewDB(&postgres.Config{
		DbHost:   "localhost",
		DbPort:   port.Port(),
		DbName:   "postgres",
		UserName: "postgres",
		Password: "postgres",
		Debug:    true,
		SSLMode:  "disable",
	})

	cleanup := func(ctx context.Context) error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close db: %w", err)
		}

		if err := pgCtr.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate container: %w", err)
		}
		return nil
	}

	return db, cleanup, nil
}
