package conf

import (
	"os"
	"time"
)

const (
	EnvIdleTimeUpperBound = "MONKEY_IDLE_UPPER_BOUND"
)

func MonkeyIdleTimeUpperBound() string {
	v := os.Getenv(EnvIdleTimeUpperBound)
	if len(v) > 0 {
		return v
	}

	return time.Minute.String()
}
