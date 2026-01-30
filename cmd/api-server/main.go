package main

import (
	"log"
	"net/http"

	"github.com/logeshwarann-dev/news-api-rest/internal/router"
)

func main() {
	log.Print("Starting Server")
	r := router.New()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("unable to start server:", err)
	}
}
