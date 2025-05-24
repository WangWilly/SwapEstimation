package utils

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type LogConfig struct {
	Level   string `env:"LOG_LEVEL,default=info"`
	Format  string `env:"LOG_FORMAT,default=json"`
	DevMode bool   `env:"LOG_DEV_MODE,default=false"`

	ProjectName string `env:"LOG_PROJECT_NAME,default=swap-estimation"`
	ProjectID   string `env:"LOG_PROJECT_ID,default=0"`
	ProjectEnv  string `env:"LOG_PROJECT_ENV,default=development"`
	ProjectVer  string `env:"LOG_PROJECT_VER,default=0.0.1"`
	ProjectHost string `env:"LOG_PROJECT_HOST,default=localhost"`
	ProjectPort string `env:"LOG_PROJECT_PORT,default=8080"`
}

////////////////////////////////////////////////////////////////////////////////

var logger *zerolog.Logger
var detailedLogger *zerolog.Logger

func InitLogging(ctx context.Context) {
	if logger != nil {
		return
	}

	////////////////////////////////////////////////////////////////////////////

	// Load environment variables
	cfg := LogConfig{}
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Level == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	////////////////////////////////////////////////////////////////////////////

	var buildLogger zerolog.Logger

	// Set the output format
	switch cfg.Format {
	case "json":
		buildLogger = zerolog.New(os.Stdout).
			With().Timestamp().Logger()
	case "console":
		buildLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			With().Timestamp().Logger()
	default:
		panic("invalid log format")
	}

	// Set the project name
	logger = &buildLogger
	detailed := buildLogger.With().
		Str("project_name", cfg.ProjectName).
		Str("project_id", cfg.ProjectID).
		Str("project_env", cfg.ProjectEnv).
		Str("project_ver", cfg.ProjectVer).
		Str("project_host", cfg.ProjectHost).
		Str("project_port", cfg.ProjectPort).
		Logger()
	detailedLogger = &detailed
}

func GetLogger() *zerolog.Logger {
	if logger == nil {
		panic("logger not initialized")
	}
	return logger
}

func GetDetailedLogger() *zerolog.Logger {
	// if calling this function for the test, initialize the logger
	if logger == nil {
		InitLogging(context.Background())
	}

	if detailedLogger == nil {
		panic("logger not initialized")
	}
	return detailedLogger
}
