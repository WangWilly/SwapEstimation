package utils

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInitLogging(t *testing.T) {
	Convey("Given the InitLogging function", t, func() {
		// Save original environment variables to restore them later
		origEnv := saveEnvironment()
		defer restoreEnvironment(origEnv)

		// Reset logger before each test
		resetLoggers := func() {
			logger = nil
			detailedLogger = nil
		}

		Convey("When called with default configuration", func() {
			resetLoggers()
			// Clear environment variables to use defaults
			os.Clearenv()

			InitLogging(context.Background())

			Convey("Then the loggers should be initialized", func() {
				So(logger, ShouldNotBeNil)
				So(detailedLogger, ShouldNotBeNil)
			})

			Convey("And the global log level should be info", func() {
				So(zerolog.GlobalLevel(), ShouldEqual, zerolog.InfoLevel)
			})
		})

		Convey("When called with debug level configuration", func() {
			resetLoggers()
			os.Clearenv()
			os.Setenv("LOG_LEVEL", "debug")

			InitLogging(context.Background())

			Convey("Then the global log level should be debug", func() {
				So(zerolog.GlobalLevel(), ShouldEqual, zerolog.DebugLevel)
			})
		})

		Convey("When called with console format", func() {
			resetLoggers()
			os.Clearenv()
			os.Setenv("LOG_FORMAT", "console")

			InitLogging(context.Background())

			Convey("Then the loggers should be initialized", func() {
				So(logger, ShouldNotBeNil)
				So(detailedLogger, ShouldNotBeNil)
			})
		})

		Convey("When called with custom project settings", func() {
			resetLoggers()
			os.Clearenv()
			os.Setenv("LOG_PROJECT_NAME", "test-project")
			os.Setenv("LOG_PROJECT_ID", "123")
			os.Setenv("LOG_PROJECT_ENV", "testing")
			os.Setenv("LOG_PROJECT_VER", "1.0.0")
			os.Setenv("LOG_PROJECT_HOST", "testhost")
			os.Setenv("LOG_PROJECT_PORT", "9090")

			InitLogging(context.Background())

			Convey("Then the loggers should be initialized with custom settings", func() {
				So(logger, ShouldNotBeNil)
				So(detailedLogger, ShouldNotBeNil)
				// Note: The detailed logger's fields can't be easily inspected directly
				// This is an implementation detail limitation
			})
		})

		Convey("When called with invalid log format", func() {
			resetLoggers()
			os.Clearenv()
			os.Setenv("LOG_FORMAT", "invalid")

			Convey("Then it should panic", func() {
				So(func() { InitLogging(context.Background()) }, ShouldPanic)
			})
		})

		Convey("When called multiple times", func() {
			resetLoggers()
			os.Clearenv()

			// First call
			InitLogging(context.Background())
			firstLogger := logger

			// Second call
			InitLogging(context.Background())
			secondLogger := logger

			Convey("Then it should not reinitialize the loggers", func() {
				So(secondLogger, ShouldPointTo, firstLogger)
			})
		})
	})
}

func TestGetLogger(t *testing.T) {
	Convey("Given the GetLogger function", t, func() {
		Convey("When logger is already initialized", func() {
			// Initialize the logger
			if logger == nil {
				InitLogging(context.Background())
			}

			l := GetLogger()

			Convey("Then it should return the initialized logger", func() {
				So(l, ShouldNotBeNil)
				So(l, ShouldPointTo, logger)
			})
		})

		Convey("When logger is not initialized", func() {
			// Reset logger
			logger = nil

			Convey("Then it should panic", func() {
				So(func() { GetLogger() }, ShouldPanic)
			})
		})
	})
}

func TestGetDetailedLogger(t *testing.T) {
	Convey("Given the GetDetailedLogger function", t, func() {
		// Save original loggers to restore them later
		origLogger, origDetailedLogger := logger, detailedLogger
		defer func() {
			logger, detailedLogger = origLogger, origDetailedLogger
		}()

		Convey("When logger is already initialized", func() {
			// Initialize the logger
			if logger == nil {
				InitLogging(context.Background())
			}

			l := GetDetailedLogger()

			Convey("Then it should return the initialized detailed logger", func() {
				So(l, ShouldNotBeNil)
				So(l, ShouldPointTo, detailedLogger)
			})
		})

		Convey("When logger is not initialized", func() {
			// Reset loggers
			logger, detailedLogger = nil, nil

			l := GetDetailedLogger()

			Convey("Then it should initialize the logger and return it", func() {
				So(l, ShouldNotBeNil)
				So(logger, ShouldNotBeNil)
				So(detailedLogger, ShouldNotBeNil)
			})
		})
	})
}

// Helper function to save environment variables
func saveEnvironment() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := make([]byte, len(e))
		copy(pair, e)
		for i := 0; i < len(e); i++ {
			if pair[i] == '=' {
				env[string(pair[:i])] = string(pair[i+1:])
				break
			}
		}
	}
	return env
}

// Helper function to restore environment variables
func restoreEnvironment(env map[string]string) {
	os.Clearenv()
	for k, v := range env {
		os.Setenv(k, v)
	}
}
