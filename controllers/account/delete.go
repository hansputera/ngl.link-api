package account

import (
	"encoding/json"
	"io"
	"net/http"
	"nglapi/global"
	"nglapi/models"

	"github.com/go-redis/redis/v9"
)

type accountDeleteBody struct {
	Slug string `json:"slug"`
	Id   string `json:"id"`
}

func AccountDelete(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	body := &accountDeleteBody{}
	if err = json.Unmarshal(data, body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data, err = global.RedisClient.Get(global.ContextConsume, body.Slug).Bytes()
	if err != nil && err != redis.Nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	acc := &models.User{}
	if err = json.Unmarshal(data, acc); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if body.Id != acc.Id {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You aren't able to delete this account!"))
		return
	}

	go global.RedisClient.Del(global.ContextConsume, acc.Slug)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
