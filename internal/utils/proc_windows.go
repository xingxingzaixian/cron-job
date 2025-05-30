//go:build windows
// +build windows

package utils

import (
	"os/exec"
	"syscall"
)

// setProcAttr 设置Windows平台的进程属性
func setProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}

// killProcessGroup Windows平台进程终止函数（空实现，因为Windows在killProcess中已处理）
func killProcessGroup(pid int) {
	// Windows平台在killProcess函数中已经使用taskkill处理了进程终止
	// 这里提供空实现以保持接口一致性
}
