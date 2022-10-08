package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"nglapi/global"
	"nglapi/jwt"
	"nglapi/models"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v9"
	"github.com/rs/xid"
)

type RefreshBody struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

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

		account.Id = xid.New().String()
		account.Slug = strings.ReplaceAll(account.InstagramUID, " ", "")

		err = global.RedisClient.Get(global.ContextConsume, account.Slug).Err()
		if err != nil && err == redis.Nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		} else if err != redis.Nil {
			rid := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(500)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Couldn't generate random id!"))
				return
			}
			account.Slug = strings.Join([]string{
				account.Slug,
				fmt.Sprint(rid),
			}, "")
		}

		data, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		go global.RedisClient.SetEx(global.ContextConsume, account.Slug, data, time.Hour*24)

		token, err := jwt.GetToken(account.Slug)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{
			"token": *token,
			"id": account.Id,
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

		body := &RefreshBody{}

		if err = json.Unmarshal(data, body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		parsed := jwt.ParseJWTToken(body.Token)
		if parsed == nil || parsed.Claims.(jwt.NglClaims).ExpiresAt.Unix() <= time.Now().Unix() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid token"))
			return
		} else {
			acc := &models.User{}

			data, err = global.RedisClient.Get(global.ContextConsume, parsed.Claims.(jwt.NglClaims).IgId).Bytes()
			if err != nil || err == redis.Nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Couldn't find your identity!"))
				return
			}

			if err = json.Unmarshal(data, acc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			} else if acc.Id != body.Id {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("You aren't able to refresh this user token!"))
				return
			}

			token, err := jwt.GetToken(parsed.Claims.(jwt.NglClaims).IgId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			if err = json.NewEncoder(w).Encode(map[string]string{
				"token": *token,
				"id": acc.Id,
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
	})

	return router
}
