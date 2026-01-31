package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/logger"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/validator"
)

type NewsStorer interface {
	//Create News
	Create(model.NewsRecord) (model.NewsRecord, error)
	//Get All News
	FindAll() ([]model.NewsRecord, error)
	//Get News By Id
	FindById(uuid.UUID) (model.NewsRecord, error)
	//Update News By Id
	UpdateById(uuid.UUID, model.NewsRecord) (model.NewsRecord, error)
	//Delete News By Id
	DeleteById(uuid.UUID) error
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := logger.FromContext(r.Context())
		log.Info("postnews request recieved")
		var newsRequestBody model.NewsRecord
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("request decode failed, invalid request", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := validator.ValidateNewsRequest(newsRequestBody); err != nil {
			log.Error("validation error, invalid request", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		_, err := ns.Create(newsRequestBody)
		if err != nil {
			log.Error("failed adding news in db", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("getallnews request recieved")

		newsRecords, err := ns.FindAll()
		if err != nil {
			log.Error("failed finding all news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newsRecords = []model.NewsRecord{{
			Author:    "author",
			Title:     "test-title",
			Summary:   "test-summary",
			Content:   "test-content",
			Source:    "test-url",
			CreatedAt: "2026-01-30T18:35:43+05:30",
			Tags:      []string{"test-tag"},
		}, {Author: "124",
			Title:     "test-title",
			Summary:   "test-summary",
			Content:   "test-content",
			Source:    "test-url",
			CreatedAt: "2026-01-30T18:35:43+05:30",
			Tags:      []string{"test-tag"},
		},
		}
		if err := json.NewEncoder(w).Encode(newsRecords); err != nil {
			log.Error("failed encoding records", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("getnewsbyid request recieved")
		// id := r.PathValue("news_id")
		id := strings.TrimPrefix(r.URL.Path, "/news/")
		newsId, err := validator.ValidateNewsId(id)
		if err != nil {
			log.Error("invalid newsId", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newsRecord, err := ns.FindById(newsId)
		if err != nil {
			log.Error("failed finding newsrecord by id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newsRecord = model.NewsRecord{
			Author:    "124",
			Title:     "test-title",
			Summary:   "test-summary",
			Content:   "test-content",
			Source:    "test-url",
			CreatedAt: "2026-01-30T18:35:43+05:30",
			Tags:      []string{"test-tag"},
		}
		if err := json.NewEncoder(w).Encode(newsRecord); err != nil {
			log.Error("failed encoding response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("updatenewsbyid request recieved")
		var req model.NewsRecord
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("request decoding failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := validator.ValidateNewsRequest(req); err != nil {
			log.Error("invalid request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/news/")
		newsId, err := validator.ValidateNewsId(id)
		if err != nil {
			log.Error("invalid news id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := ns.UpdateById(newsId, req)
		if err != nil {
			log.Error("unable to update news by id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp = model.NewsRecord{
			Author:    "test-author",
			Title:     "test-title",
			Summary:   "test-summary",
			Content:   "test-content",
			Source:    "https://google.com",
			CreatedAt: "2026-01-30T18:35:43+05:30",
			Tags:      []string{"test-tag"},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("response encoding failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("deletenewsbyid request recieved")

		id := strings.TrimPrefix(r.URL.Path, "/news/")
		newsId, err := validator.ValidateNewsId(id)
		if err != nil {
			log.Error("invalid newsid", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := ns.DeleteById(newsId); err != nil {
			log.Error("failed to delete news by id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
