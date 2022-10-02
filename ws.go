package main

import (
	"log"
	"net/http"
	"nglapi/jwt"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		token := r.Header.Get("X-Token")
		if len(token) < 5 {
			return false
		} else {
			return jwt.VerifyJWTToken(token)
		}
	},
}

func ReceiveWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer ws.Close()
}
