package scrapygo

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func Test_handlerError(t *testing.T) {
	rs, err := ts.Client().Get(ts.URL + "/v1/error")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusInternalServerError, rs.StatusCode)

	defer rs.Body.Close()

	body, _ := io.ReadAll(rs.Body)
	assert.Equal(t, `{"error":"internal server error"}`, string(body))
}

func Test_handlerReadiness(t *testing.T) {
	rs, err := ts.Client().Get(ts.URL + "/v1/readiness")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	defer rs.Body.Close()

	body, _ := io.ReadAll(rs.Body)
	assert.Equal(t, `{"status":"ok"}`, string(body))
}
