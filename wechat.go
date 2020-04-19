package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
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

	if req.MsgType == "event" {
		if req.Event == "subscribe" {
			userIDString := req.EventKey[8:]
			userID, err := strconv.ParseInt(userIDString, 10, 64)
			if err != nil {
				c.String(403, err.Error())
				return
			}

			connectionIDs, err := serviceGlobal.UserIDToConnectionIDs(userID)
			for _, connectionID := range connectionIDs {
				output, err := awsAPIGateway.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
					ConnectionId: &connectionID,
					Data:         []byte("subscribe"),
				})

				if err != nil {
					log.Println(err)
					log.Println(output)
				}
			}
		}
		c.Status(200)
		return
	}

	resp := message.NewText(`<a href="https://www.kan-fun.com/">kan</a>`)
	resp.SetToUserName(req.FromUserName)
	resp.SetFromUserName(req.ToUserName)
	resp.SetCreateTime(req.CreateTime)
	resp.SetMsgType("text")

	c.XML(200, resp)
}

type weChatQRClientReqStruct struct {
	ConnectionID string `json:"connectionId"`
}

func weChatQR(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	var clientReq weChatQRClientReqStruct
	err = json.Unmarshal(body, &clientReq)
	if err != nil {
		log.Println(err)
		c.String(403, "")
		return
	}

	connectionID := clientReq.ConnectionID

	if connectionID == "" {
		log.Println(body)
		c.String(403, "")
		return
	}

	userID, err := serviceGlobal.connectionIDToUserID(connectionID)
	if err != nil {
		log.Println(err)
		c.String(403, "Fail to get user id")
		return
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", mpAccessToken)

	bodyString := fmt.Sprintf(`{"expire_seconds": 604800, "action_name": "QR_SCENE", "action_info": {"scene": {"scene_id": %d}}}`, userID)
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

	output, err := awsAPIGateway.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionID,
		Data:         []byte(wechatResp.Ticket),
	})

	if err != nil {
		log.Println(err)
		log.Println(output)
	}

	c.Status(200)
}
