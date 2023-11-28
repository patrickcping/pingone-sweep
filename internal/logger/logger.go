package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var once sync.Once
var logger zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {

		logLevelEnv := os.Getenv("P1_SWEEP_LOG")
		var logLevel zerolog.Level

		switch logLevelEnv {
		case "DEBUG":
			logLevel = zerolog.DebugLevel
		case "INFO":
			logLevel = zerolog.InfoLevel
		case "WARN":
			logLevel = zerolog.WarnLevel
		case "ERROR":
			logLevel = zerolog.ErrorLevel
		case "FATAL":
			logLevel = zerolog.FatalLevel
		case "PANIC":
			logLevel = zerolog.PanicLevel
		case "NOLEVEL":
			logLevel = zerolog.NoLevel
		default:
			logLevel = zerolog.Disabled
		}

		logPath := os.Getenv("P1_SWEEP_LOG_PATH")

		var output io.Writer
		if logPath != "" && logLevel != zerolog.Disabled {

			var err error
			output, err = os.Create(logPath)
			if err != nil {
				panic(err)
			}

		} else {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
		}

		logger = zerolog.New(output).
			Level(logLevel).
			With().
			Timestamp().
			Logger()
	})
	return logger
}
