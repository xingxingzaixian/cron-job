package utils

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestExecShell(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var command string

	// 根据操作系统选择不同的测试命令
	if runtime.GOOS == "windows" {
		command = "echo Hello World"
	} else {
		command = "echo 'Hello World'"
	}

	output, err := ExecShell(ctx, command)
	if err != nil {
		t.Fatalf("ExecShell failed: %v", err)
	}

	if output == "" {
		t.Fatal("ExecShell returned empty output")
	}

	// 检查输出是否包含期望的内容
	if len(output) == 0 {
		t.Fatal("Output is empty")
	}

	t.Logf("Command: %s", command)
	t.Logf("Output: %s", output)
	t.Logf("OS: %s", runtime.GOOS)
}

func TestExecShellTimeout(t *testing.T) {
	// 创建一个很短的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var command string
	if runtime.GOOS == "windows" {
		// Windows下的长时间运行命令
		command = "ping -n 10 127.0.0.1"
	} else {
		// Unix下的长时间运行命令
		command = "sleep 5"
	}

	output, err := ExecShell(ctx, command)
	if err == nil {
		t.Fatal("Expected timeout error, but got none")
	}

	if err.Error() != "timeout killed" {
		t.Fatalf("Expected 'timeout killed' error, got: %v", err)
	}

	t.Logf("Timeout test passed. Output: %s, Error: %v", output, err)
}

func TestConvertEncodingWindows(t *testing.T) {
	// 这个测试只在Windows上运行
	if runtime.GOOS != "windows" {
		t.Skip("convertEncoding is Windows-specific")
	}

	// 测试UTF-8字符串（应该保持不变）
	input := "Hello World"
	output := convertEncoding(input)

	if output != input {
		t.Errorf("Expected %s, got %s", input, output)
	}
}

func TestCrossPlatformBehavior(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试跨平台的基本命令
	var commands []string
	if runtime.GOOS == "windows" {
		commands = []string{
			"echo test",
			"dir /b",
			"ver",
		}
	} else {
		commands = []string{
			"echo test",
			"ls",
			"uname",
		}
	}

	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			output, err := ExecShell(ctx, cmd)
			if err != nil {
				t.Errorf("Command '%s' failed: %v", cmd, err)
			}
			if output == "" {
				t.Errorf("Command '%s' returned empty output", cmd)
			}
			t.Logf("Command '%s' output: %s", cmd, output)
		})
	}
}
