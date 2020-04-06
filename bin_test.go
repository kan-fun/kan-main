package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBin(t *testing.T) {
	// ✅ Success
	req, _ := http.NewRequest("GET", "/bin", nil)

	q := req.URL.Query()
	q.Add("platform", "linux")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ---
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
	//

	// ❌ Failure for No Platform
	req, _ = http.NewRequest("GET", "/bin", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// ---
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "No Platform", w.Body.String())
	//
}
