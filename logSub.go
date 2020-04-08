package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/kan-fun/kan-server-core/model"
)

var subUpgrader = websocket.Upgrader{}

func logSub(c *gin.Context) {
	// Check Signature and Get User(Contain id, secret_key)
	_, err := checkSignature(c, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	conn, err := subUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}
	defer conn.Close()

	_, idBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	idString := string(idBytes)

	id, err := strconv.Atoi(idString)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	var _log model.Log
	db.Select("id").Where("id = ?", id).First(&_log)
	if _log.ID == 0 {
		c.String(403, err.Error())
		return
	}

	reversedID := reverse(idString)

	var lastAutoID int64 = 0
	var contents []string

	for {
		contents, lastAutoID, err = serviceGlobal.logGetToEnd(reversedID, lastAutoID)
		if err != nil {
			c.String(403, err.Error())
			return
		}

		for _, content := range contents {
			conn.WriteMessage(websocket.TextMessage, []byte(content))
		}

		time.Sleep(2 * time.Second)
	}
}
