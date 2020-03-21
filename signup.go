package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	. "kan-server-core/model"
)

func signup(c *gin.Context) {
	email, ok := c.GetPostForm("email")
	if !ok {
		c.String(403, "No Email")
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.String(403, "No Password")
		return
	}

	accessKey, err := generateKey()
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	secretKey, err := generateKey()
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	user := &User{
		Email:     email,
		Password:  hashPassword(password),
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	if err := db.Create(user).Error; err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	cEmail := &ChannelEmail{
		UserID:  user.ID,
		Address: email,
		Count:   100,
	}

	if err := db.Create(cEmail).Error; err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	token, err := generateIDToken(fmt.Sprint(user.ID))
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
