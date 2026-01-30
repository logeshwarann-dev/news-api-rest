package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/logeshwarann-dev/news-api-rest/internal/handler"
)

func Test_PostNews(t *testing.T) {
	testCases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			//Act
			handler.PostNews()(w, r)
			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected: %d, got: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_GetAllNews(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			// Act
			handler.GetAllNews()(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected: %d, got: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_GetNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			//Act
			handler.GetNewsByID()(w, r)

			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected: %d, got: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_UpdateNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	for _, tc := range testcases {
		//Arrange
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/", nil)
		//Act
		handler.UpdateNewsByID()(w, r)
		//Assert
		if w.Result().StatusCode != tc.expectedStatus {
			t.Errorf("expected: %d, got: %d", tc.expectedStatus, w.Result().StatusCode)
		}
	}
}

func Test_DeleteNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "not implemented",
			expectedStatus: http.StatusNotImplemented,
		},
	}

	for _, tc := range testcases {
		//Arrange
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/", nil)
		//Act
		handler.DeleteNewsByID()(w, r)
		//Assert
		if w.Result().StatusCode != tc.expectedStatus {
			t.Errorf("expected: %d, got: %d", tc.expectedStatus, w.Result().StatusCode)
		}
	}
}
