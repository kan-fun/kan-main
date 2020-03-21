package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewKey(t *testing.T) {
	dropAndMigrate()

	const email = "h@h.com"
	const password = "pwd123456"

	createUser(email, password)

	// ✅ Success
	data := url.Values{"id": {"1"}}
	w := post(data, "/view-key")

	// ---
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 119, len(w.Body.String()))
	//

	// ❌ Failure
	data = url.Values{"id": {"2"}}
	w = post(data, "/view-key")

	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "Not Found User", w.Body.String())
	//
}
