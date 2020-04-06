package main

import (
	"log"

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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recvMT")
		log.Print(mt)
		log.Printf("recv: %s", message)
	}
}
