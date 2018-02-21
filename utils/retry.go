package utils

import (
	"time"

	"github.com/rs/zerolog"
)

// Retry helper method to for connect retries in a sane way
func Retry(attempts int, sleep time.Duration, factor float32, logger zerolog.Logger, toRetry func() error) (err error) {
	for i := 0; ; i++ {
		err = toRetry()
		if err == nil {
			return nil
		}

		if i >= (attempts - 1) {
			break
		}

		logger.Error().Err(err).Msgf("retry number %d in %s", i+1, sleep)
		time.Sleep(sleep)
		sleep = time.Duration(float32(sleep) * factor) // increase time to sleep by factor
	}
	logger.Error().Msgf("failed to complete in %d retries", attempts)

	return err
}
