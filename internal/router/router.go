package router

import (
	"net/http"

	"github.com/logeshwarann-dev/news-api-rest/internal/handler"
)

func New() *http.ServeMux {
	//Setup new server mux
	r := http.NewServeMux()

	//Create News
	r.HandleFunc("POST /news", handler.PostNews())
	//Get all News
	r.HandleFunc("GET /news", handler.GetAllNews())
	//Get News By Id
	r.HandleFunc("GET /news/{news_id}", handler.GetNewsByID())
	//Update News By Id
	r.HandleFunc("PUT /news/{news_id}", handler.UpdateNewsByID())
	//Delete News By Id
	r.HandleFunc("DELETE /news/{news_id}", handler.DeleteNewsByID())

	return r
}
