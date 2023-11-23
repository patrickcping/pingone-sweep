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
		logLevel := zerolog.DebugLevel

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		logger = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger()
	})
	return logger
}
