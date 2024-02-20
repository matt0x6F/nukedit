package reddit

import (
	"strconv"
	"time"
)

func evaluateLimit(remaining, secondsUntilReset string) error {
	if remaining == "0" {
		resetSeconds, err := strconv.Atoi(secondsUntilReset)
		if err != nil {
			return err
		}

		duration := time.Duration(resetSeconds) * time.Second

		// sleep until the rate limit is reset
		time.Sleep(duration)
	}

	return nil
}
