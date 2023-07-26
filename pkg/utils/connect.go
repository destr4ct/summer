package utils

import (
	"time"
)

var (
// ErrAttemptsExhausted =
)

func DoWithAttempts[v any](trial func() (v, error), times int) (res v, err error) {
	ticker := time.Tick(time.Second)

	for i := 0; i < times; i += 1 {
		res, err = trial()
		if err == nil {
			break
		}
		<-ticker
	}

	return
}
