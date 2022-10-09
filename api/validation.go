package api

import (
	"net/http"
	"nglapi/controllers/validation"

	"github.com/go-chi/chi/v5"
)

func ValidationRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Let me validate your request!"))
	})

	router.Post("/token", validation.ValidationToken)

	return router
}
