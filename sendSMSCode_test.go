package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMSEmailCode(t *testing.T) {
	// ✅ Success
	data := url.Values{"number": {"17080056600"}}
	w := testReq("post", "/send-sms-code", data, nil, "")

	// ---
	assert.Equal(t, 200, w.Code)
	//

	// ❌ Failure
	data = url.Values{}
	w = testReq("post", "/send-sms-code", data, nil, "")

	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "No Phone Number", w.Body.String())
	//
}
