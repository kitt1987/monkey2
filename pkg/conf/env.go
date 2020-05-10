package conf

import (
	"os"
	"strconv"
	"time"
)

const (
	EnvCoffeeTimeUpperBound = "COFFEE_TIME"
	EnvNameLength           = "NAME_LENGTH"
	EnvWriteOnlyLength      = "LENGTH_WRITE_ONCE"
	EnvPercentageFileOP     = "PERCENTAGE_FILE_OPERATION"
)

func CoffeeTimeUpperBound() string {
	v := os.Getenv(EnvCoffeeTimeUpperBound)
	if len(v) > 0 {
		return v
	}

	return time.Minute.String()
}

func NameLength() int {
	return envInt(EnvNameLength)
}

func WriteOnceLengthUpperBound() int {
	return envInt(EnvWriteOnlyLength)
}

func PercentageFileOP() int {
	return envInt(EnvPercentageFileOP)
}

func envInt(key string) int {
	i, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	if err != nil {
		panic(os.Getenv(key))
	}

	return int(i)
}
