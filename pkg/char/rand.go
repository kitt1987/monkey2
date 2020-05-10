package char

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/op"
	"math/rand"
	"time"
)

const FullPercent = PercentageWithoutSign(100)
type PercentageWithoutSign uint

func (p PercentageWithoutSign) Validate() {
	if p > 100 {
		panic(p)
	}
}

type WorktreeObjectBias []PercentageWithoutSign

func NewObjectBias() WorktreeObjectBias {
	return make([]PercentageWithoutSign, op.TotalFSObject)
}

func (b WorktreeObjectBias) Set(ob op.WorktreeObject, percentage PercentageWithoutSign) {
	b[int(ob)] = percentage
}

func (b WorktreeObjectBias) Complete() {
	base := PercentageWithoutSign(0)
	for i, v := range b {
		base += v
		b[i] = base
	}

	if base != FullPercent {
		panic(b)
	}
}

func (b WorktreeObjectBias) RandomObject() op.WorktreeObject {
	indicator := randomN(100)
}

type WorktreeOPBias []PercentageWithoutSign

func NewFileOPBias() WorktreeOPBias {
	return make([]PercentageWithoutSign, op.TotalFSOP)
}

func NewDirOPBias() WorktreeOPBias {
	return make([]PercentageWithoutSign, op.TotalFSOP-1)
}

func (b WorktreeOPBias) Set(op op.WorktreeOP, percentage PercentageWithoutSign) {
	b[int(op)] = percentage
}

func (b WorktreeOPBias) Complete(ob op.WorktreeObject) {

}

func randomFSOp(obBias WorktreeObjectBias, opBias WorktreeOPBias) (fsObj op.WorktreeObject, fsOP op.WorktreeOP) {
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