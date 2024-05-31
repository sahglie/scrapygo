package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func Test_handlerUserCreate(t *testing.T) {
	reqBody := strings.NewReader(`{"name": "steven hansen"}`)
	rs, err := ts.Client().Post(ts.URL+"/v1/users", "Content-Type: application/json", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, rs.StatusCode)

	defer rs.Body.Close()

	//body, _ := io.ReadAll(rs.Body)
	//body = bytes.TrimSpace(body)

	//assert.Equal(t, `{"error":"internal server error"}`, string(body))
}
