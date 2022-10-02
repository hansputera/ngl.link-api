package api

import (
	"encoding/json"
	"io"
	"net/http"
	"nglapi/global"
	"nglapi/jwt"
	"nglapi/models"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/xid"
)

func AccountRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var account models.User

		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else if len(account.InstagramUID) < 5 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("You need to enter instagram username correctly!"))
			return
		}

		// TODO: random slug if user is already registered

		account.Id = xid.New().String()
		account.Slug = strings.ReplaceAll(account.InstagramUID, " ", "")

		data, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		go global.RedisClient.SetEx(global.ContextConsume, account.Slug, data, time.Hour*24)

		token, err := jwt.GetToken(account.Id, account.InstagramUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{
			"Token": *token,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})

	router.Patch("/refresh", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		parsed := jwt.ParseJWTToken(string(data))
		if parsed == nil || parsed.Claims.(jwt.NglClaims).ExpiresAt.Unix() <= time.Now().Unix() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid token"))
			return
		} else {
			token, err := jwt.GetToken(parsed.Claims.(jwt.NglClaims).UserId, parsed.Claims.(jwt.NglClaims).IgId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			if err = json.NewEncoder(w).Encode(map[string]string{
				"Token": *token,
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
	})

	return router
}
