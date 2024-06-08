package scrapygo

import (
	"github.com/go-testfixtures/testfixtures/v3"
	"net/http/httptest"
	"os"
	"scrapygo/internal/config"
	"testing"
)

var (
	app      Application
	ts       *httptest.Server
	fixtures *testfixtures.Loader
)

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	app = Application{AppConfig: config.NewTestConfig()}

	var err error
	fixtures, err = app.TestFixtures()
	if err != nil {
		panic(err)
	}

	ts = httptest.NewServer(app.Routes())
	defer ts.Close()

	code := m.Run()
	os.Exit(code)
}
