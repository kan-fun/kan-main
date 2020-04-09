package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/kan-fun/kan-server-core/model"
)

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

var pubUpgrader = websocket.Upgrader{}

func logPub(c *gin.Context) {
	// Check Signature and Get User(Contain id, secret_key)
	user, err := checkSignature(c, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}

	conn, err := pubUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(403, err.Error())
		return
	}
	defer conn.Close()

	_, topic, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	_log := &model.Log{
		UserID: user.ID,
		Topic:  string(topic),
	}

	if err := db.Create(_log).Error; err != nil {
		log.Println(err)
		return
	}

	reversedID := reverse(strconv.FormatUint(uint64(_log.ID), 10))

	for {
		_, contentBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read err:", err)
			break
		}

		content := string(contentBytes)

		if err := serviceGlobal.logPut(reversedID, content); err != nil {
			log.Println(err)
		}
	}
}
