// +build !linux

package side

import (
	"os/exec"
)

func setTermSig(_ *exec.Cmd) {
}
