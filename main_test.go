package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"kan-core"
	. "kan-server-core/model"
)

var router *gin.Engine

func init() {
	setup(true)
	router = setupRouter()
	service_global = MockService{}
}

func dropDB() {
	db.DropTable(&User{})
	db.DropTable(&ChannelEmail{})
}

func dropAndMigrate() {
	dropDB()
	autoMigrate()
}

func post(data url.Values, url string) *httptest.ResponseRecorder {
	body := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		"email": {email},
		"password": {password},
		"code": {raw},
		"code_hash": {sign.HashString(raw, secretKey_global)},
		"channel_id": {email},
	}

	w := post(data, "/signup")

	if w.Code != 200 {
		panic("createUser Fail")
	}

	return w
}
