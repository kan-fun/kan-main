package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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
	db.DropTable(&model.ChannelWeChat{})
	db.DropTable(&model.Task{})
}

func dropAndMigrate() {
	dropDB()
	autoMigrate()
}

func testReq(method string, urlString string, data interface{}, commonParameter *core.CommonParameter, signature string) *httptest.ResponseRecorder {
	var body io.Reader

	if method == "post" {
		switch v := reflect.ValueOf(data); v.Kind() {
		case reflect.String:
			body = bytes.NewBuffer([]byte(v.String()))
		default:
			formData := data.(url.Values)
			body = strings.NewReader(formData.Encode())
		}
	}

	var req *http.Request
	if method == "post" {
		req, _ = http.NewRequest("POST", urlString, body)
	} else {
		req, _ = http.NewRequest("GET", urlString, nil)
	}

	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.String:
		req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	default:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

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

	w := testReq("post", "/signup", data, nil, "")

	if w.Code != 200 {
		panic("createUser Fail")
	}

	return w
}
