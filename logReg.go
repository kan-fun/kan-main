package main

import (
	"github.com/gin-gonic/gin"
)

func logReg(c *gin.Context) {
	println("------x-ca-deviceid")
	println(c.GetHeader("x-ca-deviceid"))
	c.JSON(200, gin.H{})
}
