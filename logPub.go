package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/kan-fun/kan-server-core/model"
)

// Time allowed to write a message to the peer.
const writeWait = 10 * time.Second

// Time allowed to read the next pong message from the peer.
const pongWait = 290 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
const pingPeriod = 260 * time.Second

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func sendPing(quit chan struct{}, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println(err)
			}
		case <-quit:
			return
		}
	}
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

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	quit := make(chan struct{})
	defer close(quit)

	go sendPing(quit, conn)

	_, topic, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	_, logTypeBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	logType, err := strconv.Atoi(string(logTypeBytes))
	if err != nil {
		log.Println(err)
		return
	}

	task := &model.Task{
		UserID: user.ID,
		Topic:  string(topic),
		Type:   uint8(logType),
	}

	if err := db.Create(task).Error; err != nil {
		log.Println(err)
		return
	}

	reversedID := reverse(strconv.FormatUint(uint64(task.ID), 10))

	for {
		_, contentBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				db.Model(&task).Update("status", 1)

				err = serviceGlobal.weChatNotify("oOCN8xCIjo5QXoDXokJO6Knib618", task.Topic, true)
				if err != nil {
					log.Println(err)
				}
			} else if websocket.IsCloseError(err, 4000) {
				// Todo: do sth if user want to get notify when exit code not 0
				db.Model(&task).Update("status", 2)

				err = serviceGlobal.weChatNotify("oOCN8xCIjo5QXoDXokJO6Knib618", task.Topic, false)
				if err != nil {
					log.Println(err)
				}
			} else {
				// Todo: do sth if user want to get notify when websocket disconnect abnormal
				db.Model(&task).Update("status", 3)

				log.Println(err.Error())

				err = serviceGlobal.weChatNotify("oOCN8xCIjo5QXoDXokJO6Knib618", task.Topic, false)
				if err != nil {
					log.Println(err)
				}
			}

			break
		}

		content := string(contentBytes)

		if err := serviceGlobal.logPut(reversedID, content); err != nil {
			log.Println(err)
		}
	}
}
