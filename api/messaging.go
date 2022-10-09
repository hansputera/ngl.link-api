package api

import (
	"net/http"
	"nglapi/controllers/messaging"

	"github.com/go-chi/chi/v5"
)

func MessagingRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", messaging.MessagingCreate)

	return router
}
