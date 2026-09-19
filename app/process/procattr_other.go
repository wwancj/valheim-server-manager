//go:build !windows

package process

import "os/exec"

func configureProcess(cmd *exec.Cmd) {
	// No special process attributes needed on Linux/macOS
}