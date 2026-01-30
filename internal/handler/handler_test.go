package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/handler"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
)

func Test_PostNews(t *testing.T) {
	testCases := []struct {
		name           string
		request        io.Reader
		mockStore      mockNewsStore
		expectedStatus int
	}{
		{
			name:           "incorrect_request_body",
			request:        strings.NewReader(`{`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request",
			request: strings.NewReader(`
			{
			"id": "c2f92052-348f-4372-b4bc-43dbbc88445a",
			"author": "test-author",
			"title": "test-title",
			"summary": "test-summary",
			"content": "test-content",
			"source": "https://www.google.com",
			"createdAt": "",
			"tags": ["test-tag"]
			}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db_error",
			request: strings.NewReader(`
			{
			"id": "c2f92052-348f-4372-b4bc-43dbbc88445a",
			"author": "test-author",
			"title": "test-title",
			"summary": "test-summary",
			"content": "test-content",
			"source": "https://www.google.com",
			"createdAt": "2026-01-30T18:35:43+05:30",
			"tags": ["test-tag"]
			}`),
			mockStore:      mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "valid_request",
			request: strings.NewReader(`
			{
			"id": "c2f92052-348f-4372-b4bc-43dbbc88445a",
			"author": "test-author",
			"title": "test-title",
			"summary": "test-summary",
			"content": "test-content",
			"source": "https://google.com",
			"createdAt": "2026-01-30T18:35:43+05:30",
			"tags": ["test-tag"]
			}`),
			expectedStatus: http.StatusCreated,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", tc.request)
			//Act
			handler.PostNews(tc.mockStore)(w, r)
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

type mockNewsStore struct {
	errState bool
}

func (mns mockNewsStore) Create(newsRecord model.NewNewsRecord) (model.NewNewsRecord, error) {
	if mns.errState {
		return newsRecord, errors.New("failed to create news")
	}
	return newsRecord, nil
}

func (mns mockNewsStore) FindAll() (newsRecord model.NewNewsRecord, err error) {
	if mns.errState {
		return newsRecord, errors.New("failed to find news")
	}
	return newsRecord, nil
}

func (mns mockNewsStore) FindById(id uuid.UUID) (newsRecord model.NewNewsRecord, err error) {
	if mns.errState {
		return newsRecord, errors.New("failed to find news by id")
	}
	return newsRecord, nil
}

func (mns mockNewsStore) UpdateById(id uuid.UUID) (newsRecord model.NewNewsRecord, err error) {
	if mns.errState {
		return newsRecord, errors.New("failed to update news")
	}
	return newsRecord, nil
}

func (mns mockNewsStore) DeleteById(id uuid.UUID) (newsRecord model.NewNewsRecord, err error) {
	if mns.errState {
		return newsRecord, errors.New("failed to delete news")
	}
	return newsRecord, nil
}
