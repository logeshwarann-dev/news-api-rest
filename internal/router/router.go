package router

import (
	"net/http"

	"github.com/logeshwarann-dev/news-api-rest/internal/handler"
)

func New(ns handler.NewsStorer) *http.ServeMux {
	//Setup new server mux
	r := http.NewServeMux()

	//Create News
	r.HandleFunc("POST /news", handler.PostNews(ns))
	//Get all News
	r.HandleFunc("GET /news", handler.GetAllNews(ns))
	//Get News By Id
	r.HandleFunc("GET /news/{news_id}", handler.GetNewsByID(ns))
	//Update News By Id
	r.HandleFunc("PUT /news/{news_id}", handler.UpdateNewsByID(ns))
	//Delete News By Id
	r.HandleFunc("DELETE /news/{news_id}", handler.DeleteNewsByID())

	return r
}
