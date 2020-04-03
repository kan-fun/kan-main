package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogReg(t *testing.T) {
	dropAndMigrate()

	const email = "h@h.com"
	const password = "pwd123456"

	createUser(email, password)

	// ✅ Success
	data := url.Values{}

	w := post("/log/reg", data, nil, "")

	// ---
	assert.Equal(t, 200, w.Code)
	//
}
