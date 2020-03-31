package main

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sign "github.com/kan-fun/kan-core"
	. "github.com/kan-fun/kan-server-core/model"
)

func TestSendEmail(t *testing.T) {
	dropAndMigrate()

	const email = "h.tsai@hotmail.com"
	const password = "pwd123456"

	createUser(email, password)

	var user User
	db.Select("id, access_key, secret_key").Where("email = ?", email).First(&user)

	accessKey := user.AccessKey
	signatureNonce := "sn123"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	commonParameter := sign.CommonParameter{
		AccessKey:      accessKey,
		SignatureNonce: signatureNonce,
		Timestamp:      timestamp,
	}

	msg := "msg"
	topic := "X模型"

	specificParameter := make(map[string]string)

	specificParameter["msg"] = msg
	specificParameter["topic"] = topic

	credential, err := sign.NewCredential(user.AccessKey, user.SecretKey)
	assert.Equal(t, nil, err)

	signature := credential.Sign(commonParameter, specificParameter)

	// ✅ Success
	data := url.Values{
		"access_key":      {accessKey},
		"signature":       {signature},
		"signature_nonce": {signatureNonce},
		"timestamp":       {timestamp},
		"topic":           {topic},
		"msg":             {msg},
	}
	w := post(data, "/send-email")
	// ---
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 2, len(w.Body.String()))
	//

	// ❌ Failure for Rewrite by Middle
	data = url.Values{
		"access_key":      {accessKey},
		"signature":       {signature},
		"signature_nonce": {signatureNonce},
		"timestamp":       {timestamp},
		"topic":           {"篡改 topic"},
		"msg":             {msg},
	}
	w = post(data, "/send-email")
	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "Signature not Valid", w.Body.String())
	//

	// ❌ Failure for Count not Enough
	data = url.Values{
		"access_key":      {accessKey},
		"signature":       {signature},
		"signature_nonce": {signatureNonce},
		"timestamp":       {timestamp},
		"topic":           {topic},
		"msg":             {msg},
	}

	var cEmail ChannelEmail
	db.Model(&user).Related(&cEmail)

	cEmail.Count = 0
	db.Save(&cEmail)

	w = post(data, "/send-email")
	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "Email Count not Enough", w.Body.String())
	//

	// ❌ Failure for No Topic
	data = url.Values{
		"access_key":      {accessKey},
		"signature":       {signature},
		"signature_nonce": {signatureNonce},
		"timestamp":       {timestamp},
		"msg":             {msg},
	}

	cEmail.Count = 10
	db.Save(&cEmail)

	w = post(data, "/send-email")
	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "No topic", w.Body.String())
	//
}
