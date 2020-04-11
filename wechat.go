package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func wechat(c *gin.Context) {
	echostr := c.Query("echostr")

	c.String(200, echostr)
}
