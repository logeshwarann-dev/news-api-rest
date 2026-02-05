package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/logger"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
	"github.com/logeshwarann-dev/news-api-rest/internal/validator"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler.go -package=mockshandler

// NewsStorer represents the news store operations
type NewsStorer interface {
	//Create News
	Create(context.Context, news.Record) (news.Record, error)
	//Get All News
	FindAll(context.Context) ([]news.Record, error)
	//Get News By Id
	FindById(context.Context, uuid.UUID) (news.Record, error)
	//Update News By Id
	UpdateById(context.Context, uuid.UUID, news.Record) error
	//Delete News By Id
	DeleteById(context.Context, uuid.UUID) error
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("postnews request recieved")
		var newsRequestBody model.NewsRecord
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("request decode failed, invalid request", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var newsReq news.Record
		var err error
		if newsReq, err = validator.ValidateNewsRequest(newsRequestBody); err != nil {
			log.Error("validation error, invalid request", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		respRecord, err := ns.Create(ctx, newsReq)
		if err != nil {
			log.Error("failed adding news in db", "error", err.Error())
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.GetHttpStatus())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(respRecord); err != nil {
			log.Error("failed encoding records", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("getallnews request recieved")

		records, err := ns.FindAll(ctx)
		newsRecords := model.AllNewsRecords{
			NewsRecords: records,
		}
		if err != nil {
			log.Error("failed finding all news", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.GetHttpStatus())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
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
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("getnewsbyid request recieved")
		id := r.PathValue("news_id")
		newsId, err := validator.ValidateNewsId(id)
		if err != nil {
			log.Error("invalid newsId", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newsRecord, err := ns.FindById(ctx, newsId)
		if err != nil {
			log.Error("failed finding newsrecord by id", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.GetHttpStatus())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
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
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("updatenewsbyid request recieved")
		var req model.NewsRecord
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("request decoding failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var newsReq news.Record
		var err error
		if newsReq, err = validator.ValidateNewsRequest(req); err != nil {
			log.Error("invalid request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := r.PathValue("news_id")
		newsId, err := validator.ValidateNewsId(id)
		if err != nil {
			log.Error("invalid news id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = ns.UpdateById(ctx, newsId, newsReq)
		if err != nil {
			log.Error("unable to update news by id", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.GetHttpStatus())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("deletenewsbyid request recieved")

		id := r.PathValue("news_id")
		newsId, err := validator.ValidateNewsId(id)
		if err != nil {
			log.Error("invalid newsid", "error", err)
			var dbErr *news.CustomError
			if errors.As(err, &dbErr) {
				w.WriteHeader(dbErr.GetHttpStatus())
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := ns.DeleteById(ctx, newsId); err != nil {
			log.Error("failed to delete news by id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
