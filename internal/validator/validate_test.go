package validator_test

import (
	"testing"
	"time"

	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/store"
	"github.com/logeshwarann-dev/news-api-rest/internal/validator"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateNewNewsRequest(t *testing.T) {
	type expectations struct {
		err string
		req store.News
	}
	testcases := []struct {
		name     string
		req      model.NewsRecord
		expected expectations
	}{
		{
			name: "return_error_empty_author",
			req: model.NewsRecord{
				Author:    "",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{err: "author is empty"},
		},
		{
			name: "return_error_empty_title",
			req: model.NewsRecord{
				Author:    "test-author",
				Title:     "",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{err: "title is empty"},
		},
		{
			name: "return_error_empty_summary",
			req: model.NewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{err: "summary is empty"},
		},
		{
			name: "return_error_empty_content",
			req: model.NewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{err: "content is empty"},
		},
		{
			name: "return_error_invalid_source",
			req: model.NewsRecord{
				Author:    "",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{err: "source is empty"},
		},
		{
			name: "return_error_empty_createdAt",
			req: model.NewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "test-time",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{err: "invalid createdAt time"},
		},
		{
			name: "return_error_empty_tags",
			req: model.NewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{},
			},
			expected: expectations{err: "tags are empty"},
		},
		{
			name: "validate",
			req: model.NewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "https://www.google.com",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expected: expectations{req: store.News{
				Author:  "test-author",
				Title:   "test-title",
				Summary: "test-summary",
				Content: "test-content",
				Source:  "https://www.google.com",
				CreatedAt: func() time.Time {
					val, _ := time.Parse(time.RFC3339, "2026-01-30T18:35:43+05:30")
					return val
				}(),
				Tags: []string{"test-tag"},
			}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			validReq, err := validator.ValidateNewsRequest(tc.req)
			if tc.expected.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.req, validReq)
			}

		})
	}
}
