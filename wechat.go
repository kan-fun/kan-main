package main

import (
	"encoding/xml"

	"github.com/silenceper/wechat/message"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func wechatGet(c *gin.Context) {
	echostr := c.Query("echostr")

	c.String(200, echostr)
}

func parse(rawXMLMsgBytes []byte) (msg message.MixMessage, err error) {
	msg = message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, &msg)
	return
}

func wechatPost(c *gin.Context) {
	var req message.MixMessage
	if err := c.ShouldBindXML(&req); err != nil {
		c.String(403, err.Error())
		return
	}

	resp := message.NewText("hhhhhh")
	resp.SetToUserName(req.FromUserName)
	resp.SetFromUserName(req.ToUserName)
	resp.SetCreateTime(req.CreateTime)
	resp.SetMsgType("text")

	c.XML(200, resp)
}
