package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeChat(t *testing.T) {
	// dropAndMigrate()

	// const email = "h@h.com"
	// const password = "pwd123456"

	// createUser(email, password)

	// ✅ Success
	data := `<xml>
	<ToUserName><![CDATA[123456]]></ToUserName>
	<FromUserName><![CDATA[456789]]></FromUserName>
	<CreateTime>1348831860</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	<Content><![CDATA[this is a test]]></Content>
	<MsgId>1234567890123456</MsgId>
  </xml>
  `
	w := post("/wechat", data, nil, "")

	// ---
	assert.Equal(t, 200, w.Code)
	// println(w.Body.String())
	// assert.Equal(t, 119, len(w.Body.String()))
	//

	// // ❌ Failure
	// data = url.Values{"id": {"2"}}
	// w = post("/view-key", data, nil, "")

	// // ---
	// assert.Equal(t, 403, w.Code)
	// assert.Equal(t, "Not Found User", w.Body.String())
	// //
}
