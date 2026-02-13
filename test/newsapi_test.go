//go:build integration

package test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	NAMESPACE     = "news-service"
	SERVICENAME   = "news-api-service"
	CONTAINERPORT = 8080
)

func TestGetAllNews(t *testing.T) {
	testcases := []struct {
		name           string
		body           io.Reader
		expectedStatus int
		expectedErr    error
	}{
		{
			name:           "success_all_news",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			baseUrl := forward(t, NAMESPACE, SERVICENAME, CONTAINERPORT)
			endpoint := fmt.Sprintf("%s/news", baseUrl)
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, endpoint, tc.body)
			require.NoError(t, err)

			client := http.DefaultClient
			res, err := client.Do(req)
			assert.Equal(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}
