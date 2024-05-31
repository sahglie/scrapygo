package main

import (
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"scrapygo/internal/config"
	"testing"
)

var (
	app      application
	ts       *httptest.Server
	fixtures *testfixtures.Loader
)

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	app = application{AppConfig: config.NewTestConfig()}

	var err error
	fixtures, err = app.TestFixtures()
	if err != nil {
		panic(err)
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
