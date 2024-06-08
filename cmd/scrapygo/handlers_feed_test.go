package scrapygo

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_handlerFeedCreate(t *testing.T) {
	prepareTestDatabase()

	tests := []struct {
		name     string
		urlPath  string
		reqBody  string
		apiKey   string
		wantCode int
		wantBody string
	}{
		{
			name:     "Creates a Feed",
			urlPath:  "/v1/feeds",
			reqBody:  `{"name": "howdy", "url": "https://howdy.io/index.xml"}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusCreated,
			wantBody: "",
		},
		{
			name:     "Returns error when field 'url' is missing",
			urlPath:  "/v1/feeds",
			reqBody:  `{"name": "howdy", "url": ""}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"url: can't be blank"}`,
		},
		{
			name:     "Returns error when field 'name' is missing",
			urlPath:  "/v1/feeds",
			reqBody:  `{"name": "", "url": "https://howdy.io/index.xml"}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"name: can't be blank"}`,
		},
		{
			name:     "Returns error when field 'url' already exists",
			urlPath:  "/v1/feeds",
			reqBody:  `{"name": "The Boot.dev Blog", "url": "https://blog.boot.dev/index.xml"}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"url: already taken"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", ts.URL+tt.urlPath, strings.NewReader(tt.reqBody))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Auth "+tt.apiKey)
			rs, err := ts.Client().Do(req)

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
