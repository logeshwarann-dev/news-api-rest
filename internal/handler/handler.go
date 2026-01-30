package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/logeshwarann-dev/news-api-rest/internal/logger"
	"github.com/logeshwarann-dev/news-api-rest/internal/model"
	"github.com/logeshwarann-dev/news-api-rest/internal/validator"
)

type NewsStorer interface {
	//Create News
	Create(model.NewNewsRecord) (model.NewNewsRecord, error)
	//Get All News
	FindAll() (model.NewNewsRecord, error)
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
		log.Info("request recieved")
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

func GetAllNews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func GetNewsByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
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
