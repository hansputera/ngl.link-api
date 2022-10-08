package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func APIRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	router.Mount("/account", AccountRouter())
	router.Mount("/validation", ValidationRouter())

	return router
}
