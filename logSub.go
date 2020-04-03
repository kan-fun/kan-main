package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func logSub(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("123456"))
	if err != nil {
		println("Wrong websocket")
	}
}
