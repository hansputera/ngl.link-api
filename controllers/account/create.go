package account

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"nglapi/global"
	"nglapi/jwt"
	"nglapi/models"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/xid"
)

func AccountCreate(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
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
		"id":    account.Id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
