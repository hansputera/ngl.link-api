package validation

import (
	"encoding/json"
	"io"
	"net/http"
	"nglapi/global"
	"nglapi/jwt"
	"nglapi/models"

	"github.com/go-redis/redis/v9"
)

type validationTokenBody struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

func ValidationToken(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	body := &validationTokenBody{}
	if err = json.Unmarshal(data, body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	if len(body.Token) < 10 || len(body.Id) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must fill the request body correctly!"))
		return
	}

	parsed := jwt.ParseJWTToken(body.Token)
	if parsed == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid token"))
		return
	} else {
		acc := &models.User{}

		data, err = global.RedisClient.Get(global.ContextConsume, parsed.Claims.(*jwt.NglClaims).IgId).Bytes()
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
			w.Write([]byte("You are not allowed to view this user identity!"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(acc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			return
		}
	}
}
