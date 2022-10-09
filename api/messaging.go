package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MessagingRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		
	})

	return router
}
