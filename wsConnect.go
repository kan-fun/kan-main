package main

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func wsConnect(c *gin.Context) {
	for k, vals := range c.Request.Header {
		log.Printf("%s", k)
		for _, v := range vals {
			log.Printf("\t%s", v)
		}
	}

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	c.JSON(502, gin.H{
		"token": "token",
	})
}
