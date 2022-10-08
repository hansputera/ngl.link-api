package main

import (
	"log"
	"net/http"
	"nglapi/api"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartWeb() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte(""))
	})

	router.Mount("/api", api.APIRouter())

	port := os.Getenv("PORT")
	if len(port) < 4 {
		port = "3000"
	}

	log.Println("Listening to", port)
	if err := http.ListenAndServe(strings.Join([]string{":", port}, ""), router); err != nil {
		log.Fatal(err)
	}
}
