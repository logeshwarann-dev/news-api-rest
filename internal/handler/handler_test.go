package handler_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/handler"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/store"
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
			r := httptest.NewRequest(http.MethodPost, "/news", tc.request)
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
		response       model.AllNewsRecords
		mockStore      mockNewsStore
		expectedStatus int
	}{
		{
			name:           "db_error",
			mockStore:      mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "return_success",
			mockStore: mockNewsStore{
				news: []store.News{
					{
						Author:  "author",
						Title:   "test-title",
						Summary: "test-summary",
						Content: "test-content",
						Source:  "test-url",
						CreatedAt: func() time.Time {
							val, _ := time.Parse(time.RFC3339, "2026-01-30T18:35:43+05:30")
							return val
						}(),
						Tags: []string{"test-tag"},
					},
					{
						Author:  "124",
						Title:   "test-title",
						Summary: "test-summary",
						Content: "test-content",
						Source:  "test-url",
						CreatedAt: func() time.Time {
							val, _ := time.Parse(time.RFC3339, "2026-01-30T18:35:43+05:30")
							return val
						}(),
						Tags: []string{"test-tag"},
					},
				},
			},
			response: model.AllNewsRecords{
				NewsRecords: []store.News{
					{
						Author:  "author",
						Title:   "test-title",
						Summary: "test-summary",
						Content: "test-content",
						Source:  "test-url",
						CreatedAt: func() time.Time {
							val, _ := time.Parse(time.RFC3339, "2026-01-30T18:35:43+05:30")
							return val
						}(),
						Tags: []string{"test-tag"},
					},
					{
						Author:  "124",
						Title:   "test-title",
						Summary: "test-summary",
						Content: "test-content",
						Source:  "test-url",
						CreatedAt: func() time.Time {
							val, _ := time.Parse(time.RFC3339, "2026-01-30T18:35:43+05:30")
							return val
						}(),
						Tags: []string{"test-tag"},
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/news", nil)

			// Act
			handler.GetAllNews(tc.mockStore)(w, r)

			var actualResp model.AllNewsRecords
			if len(w.Body.Bytes()) != 0 {
				if err := json.Unmarshal(w.Body.Bytes(), &actualResp); err != nil {
					t.Errorf("response unmarshalling failed: %v", err)
					return
				}
			}
			isEqual := reflect.DeepEqual(tc.response, actualResp)
			//Assert
			if !isEqual {
				t.Errorf("expected resp: %v, got resp: %v", tc.response, actualResp)
			}
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status: %d, got status: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_GetNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		newsId         string
		response       model.NewsRecord
		mockStore      mockNewsStore
		expectedStatus int
	}{
		{
			name:           "invalid_newsId",
			newsId:         "#$21245",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "db_error",
			newsId:         "c2f92052-348f-4372-b4bc-43dbbc88445a",
			mockStore:      mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "return_success",
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			mockStore: mockNewsStore{
				news: []store.News{
					{
						Id: func() uuid.UUID {
							id, _ := uuid.Parse("c2f92052-348f-4372-b4bc-43dbbc88445a")
							return id
						}(),
						Author:  "124",
						Title:   "test-title",
						Summary: "test-summary",
						Content: "test-content",
						Source:  "test-url",
						CreatedAt: func() time.Time {
							val, _ := time.Parse(time.RFC3339, "2026-01-30T18:35:43+05:30")
							return val
						}(),
						Tags: []string{"test-tag"},
					},
				},
			},
			response: model.NewsRecord{
				Author:    "124",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "test-url",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/news/", nil)
			r.SetPathValue("news_id", tc.newsId)
			//Act
			handler.GetNewsByID(tc.mockStore)(w, r)

			var actualResp model.NewsRecord
			if len(w.Body.Bytes()) != 0 {
				if err := json.Unmarshal(w.Body.Bytes(), &actualResp); err != nil {
					t.Errorf("response unmarshalling failed: %v", err)
					return
				}
			}
			isEqual := reflect.DeepEqual(tc.response, actualResp)
			if !isEqual {
				t.Errorf("expected resp: %v, got resp: %v", tc.response, actualResp)
			}
			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status: %d, got status: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_UpdateNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		request        io.Reader
		newsId         string
		response       model.NewsRecord
		mockStore      mockNewsStore
		expectedStatus int
	}{
		{
			name:           "incorrect_request_body",
			request:        strings.NewReader(`{`),
			newsId:         "c2f92052-348f-4372-b4bc-43dbbc88445a",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "incorrect_news_id",
			request: strings.NewReader(`{
			"id": "c2f92052-348f-4372-b4bc-43dbbc88445a",
			"author": "test-author",
			"title": "test-title",
			"summary": "test-summary",
			"content": "test-content",
			"source": "https://google.com",
			"createdAt": "2026-01-30T18:35:43+05:30",
			"tags": ["test-tag"]
			}`),
			newsId:         "$%123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_request",
			request: strings.NewReader(`{
			"id": "c2f92052-348f-4372-b4bc-43dbbc88445a",
			"author": "",
			"title": "",
			"summary": "",
			"content": "test-content",
			"source": "https://google.com",
			"createdAt": "1234",
			"tags": ["test-tag"]
			}`),
			newsId:         "c2f92052-348f-4372-b4bc-43dbbc88445a",
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
			"source": "https://google.com",
			"createdAt": "2026-01-30T18:35:43+05:30",
			"tags": ["test-tag"]
			}`),
			newsId:         "c2f92052-348f-4372-b4bc-43dbbc88445a",
			mockStore:      mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "update_success",
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
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			response: model.NewsRecord{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				Content:   "test-content",
				Source:    "https://google.com",
				CreatedAt: "2026-01-30T18:35:43+05:30",
				Tags:      []string{"test-tag"},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/news/", tc.request)
			r.SetPathValue("news_id", tc.newsId)
			//Act
			handler.UpdateNewsByID(tc.mockStore)(w, r)
			var actualResp model.NewsRecord
			if len(w.Body.Bytes()) != 0 {
				if err := json.Unmarshal(w.Body.Bytes(), &actualResp); err != nil {
					t.Errorf("response unmarshalling failed: %v", err)
					return
				}
			}
			isEqual := reflect.DeepEqual(tc.response, actualResp)
			if !isEqual {
				t.Errorf("expected resp: %v, got resp: %v", tc.response, actualResp)
			}
			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status: %d, got status: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

func Test_DeleteNewsByID(t *testing.T) {
	testcases := []struct {
		name           string
		newsId         string
		mockStore      mockNewsStore
		expectedStatus int
	}{
		{
			name:           "invalid_newsId",
			newsId:         "$@123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "db_error",
			newsId:         "c2f92052-348f-4372-b4bc-43dbbc88445a",
			mockStore:      mockNewsStore{errState: true},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "delete_success",
			newsId:         "c2f92052-348f-4372-b4bc-43dbbc88445a",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/news/", nil)
			r.SetPathValue("news_id", tc.newsId)
			//Act
			handler.DeleteNewsByID(tc.mockStore)(w, r)
			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status: %d, got status: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}

type mockNewsStore struct {
	errState bool
	news     []store.News
}

func (mns mockNewsStore) Create(newsRecord store.News) (store.News, error) {
	if mns.errState {
		return newsRecord, errors.New("failed to create news")
	}
	return newsRecord, nil
}

func (mns mockNewsStore) FindAll() (newsRecords []store.News, err error) {
	if mns.errState {
		return mns.news, errors.New("failed to find news")
	}
	return mns.news, nil
}

func (mns mockNewsStore) FindById(id uuid.UUID) (newsRecord store.News, err error) {
	if mns.errState {
		return newsRecord, errors.New("failed to find news by id")
	}
	for _, n := range mns.news {
		if n.Id == id {
			newsRecord = n
			break
		}
	}
	return newsRecord, nil
}

func (mns mockNewsStore) UpdateById(id uuid.UUID, newsRecord store.News) (store.News, error) {
	if mns.errState {
		return newsRecord, errors.New("failed to update news")
	}
	return newsRecord, nil
}

func (mns mockNewsStore) DeleteById(id uuid.UUID) (err error) {
	if mns.errState {
		return errors.New("failed to delete news")
	}
	return nil
}
