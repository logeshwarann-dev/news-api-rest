package logger_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/logeshwarann-dev/news-api-rest/internal/logger"
)

func Test_CtxWithLogger(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    context.Context
		logger *slog.Logger
		exists bool
	}{
		{
			name:   "return_context_without_logger",
			ctx:    context.Background(),
			exists: false,
		},
		{
			name:   "return_context_with_existing_logger",
			ctx:    context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			exists: true,
		},
		{
			name:   "return_context_with_new_logger",
			ctx:    context.Background(),
			logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			exists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := logger.CtxWithLogger(tc.ctx, tc.logger)
			_, ok := ctx.Value(logger.CtxKey{}).(*slog.Logger)
			if ok != tc.exists {
				t.Errorf("expected: %v, got: %v", tc.exists, ok)
			}
		})
	}
}

func Test_FromContext(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name:     "return_existing_logger",
			ctx:      context.WithValue(context.TODO(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			expected: true,
		},
		{
			name:     "return_new_logger",
			ctx:      context.Background(),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log := logger.FromContext(tc.ctx)
			if tc.expected && log == nil {
				t.Errorf("expected: %v, got: %v", tc.expected, log)
			}
		})
	}
}
