package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	core "github.com/kan-fun/kan-core"
	"github.com/kan-fun/kan-server-core/model"
)

var router *gin.Engine

func init() {
	setup(true)
	router = setupRouter()
	serviceGlobal = mockService{}
}

func dropDB() {
	db.DropTable(&model.User{})
	db.DropTable(&model.ChannelEmail{})
}

func dropAndMigrate() {
	dropDB()
	autoMigrate()
}

func post(url string, data url.Values, commonParameter *core.CommonParameter, signature string) *httptest.ResponseRecorder {
	body := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if commonParameter != nil {
		req.Header.Set("Kan-Key", commonParameter.AccessKey)
		req.Header.Set("Kan-Timestamp", commonParameter.Timestamp)
		req.Header.Set("Kan-Nonce", commonParameter.SignatureNonce)

		req.Header.Set("Kan-Signature", signature)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	return w
}

func createUser(email string, password string) *httptest.ResponseRecorder {
	raw, _, err := generateCode(email)
	if err != nil {
		panic(err)
	}

	data := url.Values{
		"email":      {email},
		"password":   {password},
		"code":       {raw},
		"code_hash":  {core.HashString(raw, secretKeyGlobal)},
		"channel_id": {email},
	}

	w := post("/signup", data, nil, "")

	if w.Code != 200 {
		panic("createUser Fail")
	}

	return w
}
