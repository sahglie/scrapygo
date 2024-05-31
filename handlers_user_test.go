package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_handlerUserList(t *testing.T) {
	prepareTestDatabase()

	reqBody := strings.NewReader(`{"name": "steven hansen"}`)
	rs, err := ts.Client().Post(ts.URL+"/v1/users", "Content-Type: application/json", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, rs.StatusCode)

	defer rs.Body.Close()

	user := userParams{}
	body, _ := io.ReadAll(rs.Body)

	err = json.Unmarshal(body, &user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "steven Hansen", user.Name)

}

func Test_handlerUserCreate2(t *testing.T) {
	prepareTestDatabase()

	tests := []struct {
		name     string
		urlPath  string
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid User",
			urlPath:  "/v1/users",
			reqBody:  `{"name": "new user"}`,
			wantCode: http.StatusCreated,
			wantBody: "",
		},
		{
			name:     "Duplicate User",
			urlPath:  "/v1/users",
			reqBody:  `{"name": "steven hansen"}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"name: already taken"}`,
		},
		{
			name:     "Invalid User",
			urlPath:  "/v1/users",
			reqBody:  `{"name": ""}`,
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"name: can't be blank"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := strings.NewReader(tt.reqBody)

			rs, err := ts.Client().Post(ts.URL+tt.urlPath, "Content-Type: application/json", reqBody)
			if err != nil {
				t.Fatal(err)
			}

			defer rs.Body.Close()

			assert.Equal(t, tt.wantCode, rs.StatusCode)

			if tt.wantBody != "" {
				body, err := io.ReadAll(rs.Body)
				if err != nil {
					t.Fatal(err)
				}
				assert.Contains(t, tt.wantBody, string(body))
			}
		})
	}

}
