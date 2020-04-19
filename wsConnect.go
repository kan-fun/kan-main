package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kan-fun/kan-server-core/model"
)

type wechatQRRespStruct struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds uint32 `json:"expire_seconds"`
	URL           string `json:"url"`
}

func wsConnect(c *gin.Context) {
	id, ok := c.GetPostForm("id")
	if !ok {
		c.String(403, "No ID")
		return
	}

	var user model.User
	db.Select("id").Where("id = ?", id).First(&user)
	if user.ID == 0 {
		c.String(403, "Not Found User")
		return
	}

	connectionID := c.GetHeader("Kan-Connectionid")
	if connectionID == "" {
		c.String(403, "No Kan-Connectionid")
		return
	}

	err := serviceGlobal.newWsSession(connectionID, int64(user.ID))
	if err != nil {
		c.String(403, "Fail to new ws session")
		return
	}

	c.Status(200)
}
