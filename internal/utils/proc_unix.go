//go:build !windows
// +build !windows

package utils

import (
	"os/exec"
	"syscall"
)

// setProcAttr 设置Unix平台的进程属性
func setProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

// killProcessGroup Unix平台进程组终止函数
func killProcessGroup(pid int) {
	syscall.Kill(-pid, syscall.SIGKILL)
}
