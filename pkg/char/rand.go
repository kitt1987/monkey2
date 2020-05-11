package char

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/fs"
	"math/rand"
	"sort"
	"time"
)

const FullPercent = PercentageWithoutSign(100)

type PercentageWithoutSign uint

func (p PercentageWithoutSign) Validate() {
	if p > 100 {
		panic(p)
	}
}

type PercentageDistribution []PercentageWithoutSign

func (b PercentageDistribution) Len() int {
	return len(b)
}

func (b PercentageDistribution) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b PercentageDistribution) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b PercentageDistribution) Set(id int, percentage int) {
	b[id] = PercentageWithoutSign(percentage)
}

func (b PercentageDistribution) Complete() {
	base := PercentageWithoutSign(0)
	for i, v := range b {
		base += v
		b[i] = base
	}

	if base != FullPercent || !sort.IsSorted(b) {
		panic(b)
	}
}

func (b PercentageDistribution) RandomObject() (i int) {
	indicator := PercentageWithoutSign(randomN(100))
	var p PercentageWithoutSign
	for i, p = range b {
		if indicator < p {
			break
		}
	}

	return
}

func NewObjectBias() PercentageDistribution {
	return make([]PercentageWithoutSign, fs.TotalFSObject)
}

func NewFileOPBias() PercentageDistribution {
	return make([]PercentageWithoutSign, fs.TotalFSOP)
}

func NewDirOPBias() PercentageDistribution {
	return make([]PercentageWithoutSign, fs.TotalFSOP-1)
}

func randomFSOp(obBias, fileBias, dirBias PercentageDistribution) (fsObj fs.WorktreeObject, fsOP fs.WorktreeOP) {
	obBias.Complete()
	fileBias.Complete()
	dirBias.Complete()

	fsObj = fs.WorktreeObject(obBias.RandomObject())
	if fsObj == fs.File {
		fsOP = fs.WorktreeOP(fileBias.RandomObject())
	} else {
		fsOP = fs.WorktreeOP(dirBias.RandomObject())
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

const contentBytes = `"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()-=_+\][{}|;'":/.,<>?` + "\n"
const nameBytes = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.`

func randomName(size int) string {
	return randomString(nameBytes, size)
}

func randomText(size int) string {
	return randomString(contentBytes, size)
}

func randomString(set string, size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, size)
	for i := range b {
		b[i] = set[r.Intn(len(set))]
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
