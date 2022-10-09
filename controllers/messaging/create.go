package messaging

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nglapi/global"
	"nglapi/models"
	"time"

	"github.com/go-redis/redis/v9"
)

type messageBody struct {
	Destination string `json:"target"`
	Message     string `json:"text"`
}

func MessagingCreate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	body := &messageBody{}
	if err = json.Unmarshal(data, body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// use `.Bytes()` because I want to reuse the `data` variable.
	data, err = global.RedisClient.Get(global.ContextConsume, body.Destination).Bytes()
	if err != nil && err != redis.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	acc := &models.User{}

	if err = json.Unmarshal(data, acc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// push to notif server
	go global.RedisClient.Publish(global.ContextConsume, "new_message", map[string]string{
		"id":       acc.Id,
		"msg":      body.Message,
		"time_now": time.Now().UTC().String(),
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Sent to %s", acc.Slug)))
}
