package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	r, _, err := generateCode("main@kan.fun")
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, len(r))
}
