package news_test

import (
	"context"
	"fmt"
	"net/http"
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
			name:    "missing_author",
			context: context.Background(),
			record: news.Record{
				Title:   "Breaking NEWS",
				Summary: "A brief summary of news",
				Content: "Batman is a hero with super powers who saves people from devil",
				Source:  "https://www.google.com",
				Tags:    []string{"marvel", "test-1"},
			},
			expectedErr:    "not-null",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "create_news_success",
			context: context.Background(),
			record: news.Record{
				Author:  "Batman",
				Title:   "Breaking NEWS",
				Summary: "A brief summary of news",
				Content: "Batman is a hero with super powers who saves people from devil",
				Source:  "https://www.google.com",
				Tags:    []string{"marvel", "test-1"},
			},
			expectedStatus: http.StatusCreated,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := news.NewStore(db)
			got, err := store.Create(tc.context, tc.record)
			if err != nil {
				assert.Error(t, err)
				var storeErr *news.CustomError
				assert.ErrorAs(t, err, &storeErr)
				assert.Equal(t, tc.expectedStatus, storeErr.GetHttpStatus())
			} else {
				assert.NoError(t, err)
				assertOnNews(t, tc.record, got)
			}
		})
	}
}

func TestStore_FindById(t *testing.T) {
	testcases := []struct {
		name           string
		ctx            context.Context
		newsId         uuid.UUID
		expectedRecord news.Record
		expectedErr    string
		expectedStatus int
	}{
		{
			name:   "return_valid_news_record",
			ctx:    context.Background(),
			newsId: uuid.MustParse("67a7f4b7-2261-4121-a578-bf9da06aa0f3"),
			expectedRecord: news.Record{
				Author:  "Batman",
				Title:   "Breaking NEWS",
				Summary: "A brief summary of news",
				Content: "Batman is a hero with super powers who saves people from devil",
				Source:  "https://www.google.com",
				Tags:    []string{"marvel", "sci-fi"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "return_not_found_error",
			ctx:            context.Background(),
			newsId:         uuid.MustParse("79a7f4b7-2261-4121-a578-bf9da06aa0f3"),
			expectedErr:    "no rows",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := news.NewStore(db)
			got, err := s.FindById(tc.ctx, tc.newsId)
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
				var storeErr *news.CustomError
				assert.ErrorAs(t, err, &storeErr)
				assert.Equal(t, tc.expectedStatus, storeErr.GetHttpStatus())
			} else {
				assert.NoError(t, err)
				assertOnNews(t, tc.expectedRecord, got)
			}
		})
	}

}

func TestStore_FindAll(t *testing.T) {
	testcases := []struct {
		name            string
		ctx             context.Context
		expectedRecords []news.Record
		expectedErr     string
		expectedStatus  int
	}{
		{
			name: "return_all_records",
			ctx:  context.Background(),
			expectedRecords: []news.Record{
				{
					Author:  "Batman",
					Title:   "Breaking NEWS",
					Summary: "A brief summary of news",
					Content: "Batman is a hero with super powers who saves people from devil",
					Source:  "https://www.google.com",
					Tags:    []string{"marvel", "sci-fi"},
				},
				{
					Author:  "Superman",
					Title:   "Breaking NEWS",
					Summary: "A brief summary of news",
					Content: "Superman is a hero with super powers who saves people from devil",
					Source:  "https://www.google.com",
					Tags:    []string{"marvel", "sci-fi"},
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := news.NewStore(db)
			got, err := s.FindAll(tc.ctx)
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
				var storeErr *news.CustomError
				assert.ErrorAs(t, err, &storeErr)
				assert.Equal(t, tc.expectedStatus, storeErr.GetHttpStatus())
			} else {
				assert.NoError(t, err)
				for i, record := range tc.expectedRecords {
					assertOnNews(t, record, got[i])
				}
			}
		})
	}
}

func assertOnNews(tb testing.TB, expected, got news.Record) {
	tb.Helper()
	assert.Equal(tb, expected.Author, got.Author)
	assert.Equal(tb, expected.Title, got.Title)
	assert.Equal(tb, expected.Summary, got.Summary)
	assert.Equal(tb, expected.Content, got.Content)
	assert.Equal(tb, expected.Source, got.Source)
	assert.Equal(tb, expected.Tags, got.Tags)
	assert.NotEqual(tb, time.Time{}, got.CreateAt)
	assert.NotEqual(tb, time.Time{}, got.UpdatedAt)
	assert.Equal(tb, time.Time{}, got.DeletedAt)
}

func createTestContainer(ctx context.Context) (ctr *pgtc.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("working dir: %w", err)
	}
	sqlScripts := wd + "/testdata/sql/store.sql"

	ctr, err = pgtc.Run(
		ctx,
		"postgres:16-alpine",
		pgtc.WithInitScripts(sqlScripts),
		pgtc.WithDatabase("postgres"),
		pgtc.WithUsername("postgres"),
		pgtc.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").
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
