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

	req, err := http.NewRequest("GET", ts.URL+"/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Auth bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c")
	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	defer rs.Body.Close()

	body, _ := io.ReadAll(rs.Body)

	var user struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "steven hansen", user.Name)

}

func Test_handlerUserCreate(t *testing.T) {
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
