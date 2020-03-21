package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	. "github.com/kan-fun/kan-server-core/model"
)

func viewKey(c *gin.Context) {
	id, ok := c.GetPostForm("id")
	if !ok {
		c.String(403, "No ID")
		return
	}

	var user User
	db.Select("id, access_key, secret_key").Where("id = ?", id).First(&user)
	if user.ID == 0 {
		c.String(403, "Not Found User")
		return
	}

	c.JSON(200, gin.H{
		"AccessKey": user.AccessKey,
		"SecretKey": user.SecretKey,
	})
}
