package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "kan-server-core/model"
)

func TestSignup(t *testing.T) {
	dropAndMigrate()

	const email = "h@h.com"
	const password = "pwd123456"

	w := createUser(email, password)
	assert.Equal(t, 200, w.Code)

	var user User
	db.Take(&user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, hashPassword(password), user.Password)
}
