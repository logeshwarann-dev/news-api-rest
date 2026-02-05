package handler_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/logeshwarann-dev/news-api-rest/internal/handler"
	mockshandler "github.com/logeshwarann-dev/news-api-rest/internal/handler/mocks"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
)

func Test_PostNews(t *testing.T) {
	testCases := []struct {
		name           string
		request        io.Reader
		setup          func(tb testing.TB) handler.NewsStorer
		expectedStatus int
	}{

		{
			name:    "incorrect_request_body",
			request: strings.NewReader(`{`),
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
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
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
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
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().Create(gomock.Any(), gomock.Any()).Return(news.Record{}, errors.New("db error"))
				return mh
			},
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
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().Create(gomock.Any(), gomock.Any()).Return(news.Record{}, nil)
				return mh
			},
			expectedStatus: http.StatusCreated,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//Arrange
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/news", tc.request)
			//Act
			handler.PostNews(tc.setup(t))(w, r)
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
		setup          func(tb testing.TB) handler.NewsStorer
		expectedStatus int
	}{
		{
			name: "db_error",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().FindAll(gomock.Any()).Return([]news.Record{}, errors.New("db error"))
				return mh
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "return_success",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().FindAll(gomock.Any()).Return([]news.Record{}, nil)
				return mh
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
			handler.GetAllNews(tc.setup(t))(w, r)

			var actualResp model.AllNewsRecords
			if len(w.Body.Bytes()) != 0 {
				if err := json.Unmarshal(w.Body.Bytes(), &actualResp); err != nil {
					t.Errorf("response unmarshalling failed: %v", err)
					return
				}
			}
			//Assert
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
		setup          func(tb testing.TB) handler.NewsStorer
		expectedStatus int
	}{
		{
			name:   "invalid_newsId",
			newsId: "#$21245",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "db_error",
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(news.Record{}, errors.New("db error"))
				return mh
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "return_success",
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(news.Record{}, nil)
				return mh
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
			handler.GetNewsByID(tc.setup(t))(w, r)

			var actualResp model.NewsRecord
			if len(w.Body.Bytes()) != 0 {
				if err := json.Unmarshal(w.Body.Bytes(), &actualResp); err != nil {
					t.Errorf("response unmarshalling failed: %v", err)
					return
				}
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
		setup          func(tb testing.TB) handler.NewsStorer
		expectedStatus int
	}{
		{
			name:    "incorrect_request_body",
			request: strings.NewReader(`{`),
			newsId:  "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
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
			newsId: "$%123",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
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
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
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
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().UpdateById(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				return mh
			},
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
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().UpdateById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return mh
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
			handler.UpdateNewsByID(tc.setup(t))(w, r)
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
		setup          func(tb testing.TB) handler.NewsStorer
		expectedStatus int
	}{
		{
			name:   "invalid_newsId",
			newsId: "$@123",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				return mockshandler.NewMockNewsStorer(gomock.NewController(t))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "db_error",
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().DeleteById(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
				return mh
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "delete_success",
			newsId: "c2f92052-348f-4372-b4bc-43dbbc88445a",
			setup: func(tb testing.TB) handler.NewsStorer {
				tb.Helper()
				mh := mockshandler.NewMockNewsStorer(gomock.NewController(t))
				mh.EXPECT().DeleteById(gomock.Any(), gomock.Any()).Return(nil)
				return mh
			},
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
			handler.DeleteNewsByID(tc.setup(t))(w, r)
			//Assert
			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("expected status: %d, got status: %d", tc.expectedStatus, w.Result().StatusCode)
			}
		})
	}
}
