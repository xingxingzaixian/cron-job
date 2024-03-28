package utils

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"os/exec"
	"strconv"
	"syscall"
)

// 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("cmd", "/C", command)
	// 隐藏cmd窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	var resultChan chan Result = make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
			cmd.Process.Kill()
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return ConvertEncoding(result.output), result.err
	}
}

func ConvertEncoding(outputGBK string) string {
	// windows平台编码为gbk，需转换为utf8才能入库
	outputUTF8, err := gcharset.ToUTF8("GBK", outputGBK)
	if err != nil {
		return outputGBK
	}
	return outputUTF8
}
