package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	EnvCoffeeTimeUpperBound = "COFFEE_TIME"
	EnvNameLength           = "NAME_LENGTH"
	EnvWriteOnceLength      = "LENGTH_WRITE_ONCE"
	EnvPercentageFileOP     = "PERCENTAGE_FILE_OPERATION"
	EnvWorktree             = "WORKTREE"
	EnvSidecarStdFile       = "SIDECAR_STD_FILE"
)

func CoffeeTimeUpperBound() string {
	v := os.Getenv(EnvCoffeeTimeUpperBound)
	if len(v) > 0 {
		return v
	}

	return time.Minute.String()
}

func Worktree() string {
	wt := os.Getenv(EnvWorktree)
	if len(wt) == 0 {
		wt = filepath.Join(os.TempDir(), "monkey")
		fmt.Printf(`🚁 The workdir will be "%s"`+"\n", wt)
	}

	return wt
}

func SidecarStdFile() string {
	std := os.Getenv(EnvSidecarStdFile)
	if len(std) == 0 {
		fmt.Printf(`🚁 Stdout of the sidecar will be written to "%s"`+"\n", std)
	}

	return std
}

func NameLength() int {
	return envInt(EnvNameLength, 8, fmt.Sprintf(`🚁 Length of file/dir name would be "%d"`+"\n", 8))
}

func WriteOnceLengthUpperBound() int {
	return envInt(
		EnvWriteOnceLength, 2048,
		fmt.Sprintf(`🚁 Length of each file write would be "%d"`+"\n", 2048),
	)
}

func PercentageFileOP() int {
	return envInt(EnvPercentageFileOP)
}

func envInt(key string, def int, defHint string) int {
	if len(key) == 0 {
		fmt.Print(defHint)
		return def
	}

	i, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	if err != nil {
		panic(os.Getenv(key))
	}

	return int(i)
}
