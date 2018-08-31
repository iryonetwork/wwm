package errorChecker

import (
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
)

func LogError(err error) {
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Err(err).Logger()
		logger = logger.With().Bytes("Stack trace:", debug.Stack()).Logger()
		logger.Error().Msg("ErrorChecker: Error occurred")
	}
}

func FatalError(err error) {
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Err(err).Logger()
		logger = logger.With().Bytes("Stack trace:", debug.Stack()).Logger()
		logger.Fatal().Msg("ErrorChecker: Error occurred")
	}
}

type fataler interface {
	Fatalf(string, ...interface{})
}

func FatalTesting(t fataler, err error) {
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
