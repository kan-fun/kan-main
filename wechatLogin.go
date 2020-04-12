package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type wechatLoginRespStruct struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

func wechatLogin(c *gin.Context) {
	code, ok := c.GetQuery("code")
	if !ok {
		c.String(403, "Doesn't get code")
		return
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", mpAPPID, mpSECRET, code)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	var wechatResp wechatLoginRespStruct
	err = json.Unmarshal(body, &wechatResp)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	if wechatResp.ExpiresIn == 0 {
		c.String(403, "invalid code")
		return
	}

	// var cWeChat model.ChannelWeChat
	// db.Select("user_id").Where("mp_open_id = ?", wechatResp.OpenID).First(&cWeChat)

	// if cWeChat.UserID == 0 {
	// 	c.String(403, "Not bind WeChat on website:", wechatResp.OpenID)
	// 	return
	// }

	// token, err := generateIDToken(fmt.Sprint(cWeChat.UserID))
	// if err != nil {
	// 	log.Println(err)
	// 	c.String(403, "")
	// 	return
	// }

	c.JSON(200, wechatResp)
}
