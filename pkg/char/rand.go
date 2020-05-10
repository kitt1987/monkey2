package char

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/op"
	"math/rand"
	"time"
)

type percentageWithoutSign

type worktreeObjectBias map[op.WorktreeObject]int

func randomFSOp() (fsObj op.WorktreeObject, fsOP op.WorktreeOP) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fsObj = op.WorktreeObject(r.Intn(op.TotalFSObject))
	if fsObj == op.FSFile {
		fsOP = op.WorktreeOP(r.Intn(op.TotalFSOP))
	} else {
		fsOP = op.WorktreeOP(r.Intn(op.TotalFSOP - 1))
	}

	return
}

func randomCoffeeTime() time.Duration {
	du, err := time.ParseDuration(conf.CoffeeTimeUpperBound())
	if err != nil {
		panic(fmt.Sprintf("%s:%s", conf.CoffeeTimeUpperBound(), err))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(r.Intn(int(du.Seconds()))) * time.Second
}

func randomItem(c []string) string {
	return c[randomN(len(c))]
}

const letterBytes = `"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()-=_+\][{}|;'":/.,<>?`+"\n"

func randomText(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func randomN(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}

func randomN64(n int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(n)
}