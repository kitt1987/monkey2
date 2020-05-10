package conf

import (
	"os"
	"time"
)

const (
	EnvCoffeeTimeUpperBound = "COFFEE_TIME_UPPER_BOUND"
)

func CoffeeTimeUpperBound() string {
	v := os.Getenv(EnvCoffeeTimeUpperBound)
	if len(v) > 0 {
		return v
	}

	return time.Minute.String()
}
