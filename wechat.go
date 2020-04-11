package main

import (
	"github.com/silenceper/wechat/message"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func wechatGet(c *gin.Context) {
	echostr := c.Query("echostr")

	c.String(200, echostr)
}

func wechatPost(c *gin.Context) {
	c.XML(200, message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText("hhhhhh"),
	})
}
