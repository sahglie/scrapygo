package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"scrapygo/internal/config"
	"testing"
)

var (
	app application
	ts  *httptest.Server
)

func TestMain(m *testing.M) {
	app = application{
		AppConfig: config.NewConfigTest(),
	}

	ts = httptest.NewServer(app.routes())
	defer ts.Close()

	code := m.Run()
	os.Exit(code)
}

func Test_handlerError(t *testing.T) {
	rs, err := ts.Client().Get(ts.URL + "/v1/error")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusInternalServerError, rs.StatusCode)

	defer rs.Body.Close()

	body, _ := io.ReadAll(rs.Body)
	body = bytes.TrimSpace(body)

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
	body = bytes.TrimSpace(body)

	assert.Equal(t, `{"status":"ok"}`, string(body))
}
