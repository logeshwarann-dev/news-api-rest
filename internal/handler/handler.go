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
	Create(model.NewNewsRecord) (model.NewNewsRecord, error)
	//Get All News
	FindAll() ([]model.NewNewsRecord, error)
	//Get News By Id
	FindById(uuid.UUID) (model.NewNewsRecord, error)
	//Update News By Id
	UpdateById(uuid.UUID) (model.NewNewsRecord, error)
	//Delete News By Id
	DeleteById(uuid.UUID) (model.NewNewsRecord, error)
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := logger.FromContext(r.Context())
		log.Info("postnews request recieved")
		var newsRequestBody model.NewNewsRecord
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("request decode failed, invalid request", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := validator.ValidateNewNewsRequest(newsRequestBody); err != nil {
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
		newsRecords = []model.NewNewsRecord{{
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
		newsRecord = model.NewNewsRecord{
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

func UpdateNewsByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func DeleteNewsByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
