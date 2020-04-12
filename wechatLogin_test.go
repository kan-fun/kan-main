package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kan-fun/kan-server-core/model"
)

func TestWeChatLogin(t *testing.T) {
	dropAndMigrate()

	const email = "h.tsai@hotmail.com"
	const password = "pwd123456"

	createUser(email, password)

	var user model.User
	db.Select("id, access_key, secret_key").Where("email = ?", email).First(&user)

	cWeChat := &model.ChannelWeChat{
		UserID:   user.ID,
		MPOpenID: "MPOpenID123",
	}
	if err := db.Create(cWeChat).Error; err != nil {
		panic(err)
	}

	// ❌ Failure for No Code
	w := testReq("get", "/wechat-login", nil, nil, "")
	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "Doesn't get code", w.Body.String())
	//

	// ❌ Failure for No Code
	w = testReq("get", "/wechat-login?code=123", nil, nil, "")
	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "invalid code", w.Body.String())
	//
}
