package news_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
	"github.com/logeshwarann-dev/news-api-rest/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pgtc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
)

var db *bun.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	pgdb, cleanup, err := createTestDB(ctx)
	if err != nil {
		panic(err)
	}
	db = pgdb
	code := m.Run()
	if err := cleanup(ctx); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestStore_Create(t *testing.T) {
	testcases := []struct {
		name           string
		context        context.Context
		record         news.Record
		expectedErr    string
		expectedStatus int
	}{
		{
			name:    "return_success",
			context: context.Background(),
			record: news.Record{
				Id:        uuid.New(),
				Author:    "Batman",
				Title:     "Breaking NEWS",
				Summary:   "A brief summary of news",
				Content:   "Batman is a hero with super powers who saves people from devil",
				Source:    "https://www.google.com",
				Tags:      []string{"marvel", "test-1"},
				CreateAt:  time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedStatus: 200,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := news.NewStore(db)
			_, err := store.Create(tc.context, tc.record)
			if err != nil {
				assert.Contains(t, err, tc.expectedErr)
			}
		})
	}
}

func createTestContainer(ctx context.Context) (ctr *pgtc.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("working dir: %w", err)
	}
	sqlScripts := wd + "/testdata/sql/store.sql"

	ctr, err = pgtc.Run(
		ctx,
		"sha256:d7ead3a9d3fe4f2906d95e00edc8b36b65749a00a49c47285a35d2be95e876dd",
		pgtc.WithInitScripts(sqlScripts),
		pgtc.WithDatabase("postgres"),
		pgtc.WithUsername("postgres"),
		pgtc.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
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
