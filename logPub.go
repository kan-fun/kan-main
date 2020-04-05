package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func logPub(c *gin.Context) {
	// Check Signature and Get User(Contain id, secret_key)
	_, err := checkSignature(c, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("123456"))
	if err != nil {
		println("Wrong websocket")
	}
}
