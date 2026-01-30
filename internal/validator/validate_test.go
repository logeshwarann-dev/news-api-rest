package validator_test

import (
	"testing"

	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/validator"
)

func Test_ValidateNewNewsRequest(t *testing.T) {
	testcases := []struct {
		name        string
		req         model.NewNewsRecord
		expectedErr bool
	}{
		{
			name: "return_error_empty_author",
			req: model.NewNewsRecord{
				Author:    "",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedErr: true,
		},
		{
			name: "return_error_empty_title",
			req: model.NewNewsRecord{
				Author:    "test-author",
				Title:     "",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedErr: true,
		},
		{
			name: "return_error_empty_summary",
			req: model.NewNewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedErr: true,
		},
		{
			name: "return_error_invalid_source",
			req: model.NewNewsRecord{
				Author:    "",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedErr: true,
		},
		{
			name: "return_error_empty_createdAt",
			req: model.NewNewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "test-time",
				Tags:      []string{"test-tag"},
			},
			expectedErr: true,
		},
		{
			name: "return_error_empty_tags",
			req: model.NewNewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{},
			},
			expectedErr: true,
		},
		{
			name: "validate",
			req: model.NewNewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "https://www.google.com",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateNewNewsRequest(tc.req)
			if tc.expectedErr && err == nil {
				t.Errorf("expected err: %v, got: %v", tc.expectedErr, err)
			} else if !tc.expectedErr && err != nil {
				t.Errorf("expected err: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}
