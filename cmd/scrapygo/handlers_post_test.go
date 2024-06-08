package scrapygo

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

var postList string = `
{"data":
  [{"id":"5d0ebc37-27dd-4c8a-b28b-10fd41939e9d",
    "feed_id":"5b024738-1068-4e44-b852-9fcfc80222e0",
    "title":"The-Zen-of-Proverbs",
    "description":"20-rules",
    "url":"https://wagslane.dev/posts/zen-of-proverbs/",
    "published_at":"2023-01-07T16:00:00-08:00",
    "created_at":"2024-06-05T21:09:41.466485-07:00",
    "updated_at":"2024-06-05T21:09:41.466485-07:00"
   }]
}
`

func Test_handlerPostList(t *testing.T) {
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
			name:     "Returns Posts for a User",
			urlPath:  "/v1/posts",
			apiKey:   "bd5bf06cf44212cd15cfcbab2ce4f738223b7bdac7b34c6c1a8873c379735f6c",
			wantCode: http.StatusOK,
			wantBody: strings.Join(strings.Fields(postList), ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+tt.urlPath, nil)
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
