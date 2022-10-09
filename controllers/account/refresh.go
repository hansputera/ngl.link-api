package account

import (
	"encoding/json"
	"io"
	"net/http"
	"nglapi/global"
	"nglapi/jwt"
	"nglapi/models"
	"time"

	"github.com/go-redis/redis/v9"
)

type refreshBody struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

func AccountRefresh(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	body := &refreshBody{}

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
			"id":    acc.Id,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
