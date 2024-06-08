package scrapygo

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_handlerFeedFollowCreate(t *testing.T) {
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
			name:     "Creates a Feed Follow",
			urlPath:  "/v1/feed_follows",
			reqBody:  `{"feed_id": "27418a9b-bbb4-4dbc-b4f1-ad95322e3895"}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusCreated,
			wantBody: "",
		},
		{
			name:     "Returns error when field 'feed_id' is missing",
			urlPath:  "/v1/feed_follows",
			reqBody:  `{"feed_id": ""}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"failed to decode json"}`,
		},
		{
			name:     "Returns error when user is already following the feed",
			urlPath:  "/v1/feed_follows",
			reqBody:  `{"feed_id": "5b024738-1068-4e44-b852-9fcfc80222e0"}`,
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: `{"error":"user is already following that feed"}`,
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
