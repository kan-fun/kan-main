package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func sendSMSCode(c *gin.Context) {
	number, ok := c.GetPostForm("number")
	if !ok {
		c.String(403, "No Phone Number")
		return
	}

	raw, token, err := generateCode()
	if err != nil {
		c.String(403, "Fail to Generate Code")
		return
	}

	err = service_global.sms(number, raw)
	if err != nil {
		log.Println(err)
		c.String(403, "Fail to Send SMS Code")
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
