package char

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"math/rand"
	"time"
)

func randomFSOp() (fsObj FileSystemObject, fsOP FileSystemOP) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fsObj = FileSystemObject(r.Intn(TotalFSObject))
	if fsObj == FSFile {
		fsOP = FileSystemOP(r.Intn(TotalFSOP))
	} else {
		fsOP = FileSystemOP(r.Intn(TotalFSOP - 1))
	}

	return
}

func randomIdleTime() time.Duration {
	du, err := time.ParseDuration(conf.MonkeyIdleTimeUpperBound())
	if err != nil {
		panic(fmt.Sprintf("%s:%s", conf.MonkeyIdleTimeUpperBound(), err))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(r.Intn(int(du.Seconds()))) * time.Second
}