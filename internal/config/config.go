package config

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"scrapygo/internal/database"
)

type AppConfig struct {
	ProjectRoot  string
	DBConn       *sql.DB
	DB           *database.Queries
	Logger       *slog.Logger
	DotEnvFile   string
	FixturesPath string
}

var (
	projectRoot string
	dotEnv      string
)

func NewConfig(dotEnvFiles ...string) *AppConfig {
	projectRoot = getProjectRoot()

	if len(dotEnvFiles) > 0 {
		dotEnv = dotEnvFiles[0]
	} else {
		dotEnv = ".env"
	}

	dotEnvPath := filepath.Join(projectRoot, dotEnv)
	err := godotenv.Load(dotEnvPath)
	if err != nil {
		panic(err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	options := &slog.HandlerOptions{AddSource: true}
	logger := slog.New(slog.NewTextHandler(os.Stdout, options))

	fixturesPath := filepath.Join(projectRoot, "db", "testdata", "fixtures")

	cfg := AppConfig{
		ProjectRoot:  projectRoot,
		DotEnvFile:   dotEnv,
		Logger:       logger,
		DBConn:       db,
		DB:           database.New(db),
		FixturesPath: fixturesPath,
	}

	return &cfg
}

func NewConfigTest() *AppConfig {
	return NewConfig(".env.test")
}

func (cfg *AppConfig) TestFixtures() (*testfixtures.Loader, error) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(cfg.DBConn),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(cfg.FixturesPath),
	)

	if err != nil {
		return nil, err
	}

	return fixtures, nil
}

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../..")
}
