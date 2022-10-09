package api

import (
	"net/http"
	"nglapi/controllers/account"

	"github.com/go-chi/chi/v5"
)

func AccountRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", account.AccountCreate)
	router.Patch("/refresh", account.AccountRefresh)

	return router
}
