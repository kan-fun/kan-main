package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}
	log.Println("------")
	log.Println(string(body))

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

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", mpAccessToken)

	bodyString := fmt.Sprintf(`{"expire_seconds": 604800, "action_name": "QR_SCENE", "action_info": {"scene": {"scene_id": %d}}}`, user.ID)
	resp, err := http.Post(url, "application/json", strings.NewReader(bodyString))
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	var wechatResp wechatQRRespStruct
	err = json.Unmarshal(body, &wechatResp)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	if wechatResp.ExpireSeconds == 0 {
		log.Println(body)
		c.String(403, "")
		return
	}

	c.String(200, wechatResp.Ticket)
}
