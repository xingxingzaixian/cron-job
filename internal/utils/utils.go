package utils

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	math_rand "math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/encoding/gcharset"
	"go.uber.org/zap"
)

// Result 用于ExecShell函数的结果传递
type Result struct {
	output string
	err    error
}

// 生成长度为length的随机字符串
func RandString(length int64) string {
	sources := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	sourceLength := len(sources)
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(sourceLength)])
	}

	return string(result)
}

// 生成32位MD5摘要
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))

	return hex.EncodeToString(m.Sum(nil))
}

// 生成0-max之间随机数
func RandNumber(max int) int {
	r := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))

	return r.Intn(max)
}

// 判断文件是否存在及是否有权限访问
func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}

	return true
}

// 跨平台执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	var cmd *exec.Cmd

	// 根据操作系统选择不同的命令执行方式
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	// 设置平台特定的进程属性
	setProcAttr(cmd)

	// 使用缓冲通道，避免超时返回后命令 goroutine 阻塞在发送上造成泄漏
	resultChan := make(chan Result, 1)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()

	select {
	case <-ctx.Done():
		if cmd.Process != nil && cmd.Process.Pid > 0 {
			killProcess(cmd.Process.Pid)
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		output := result.output
		// Windows平台需要进行编码转换
		if runtime.GOOS == "windows" {
			output = convertEncoding(output)
		}
		return output, result.err
	}
}

// killProcess 跨平台进程终止函数
func killProcess(pid int) {
	if runtime.GOOS == "windows" {
		// Windows使用taskkill命令
		exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid)).Run()
	} else {
		// Unix系统使用进程组终止
		killProcessGroup(pid)
	}
}

// convertEncoding Windows平台编码转换函数
func convertEncoding(outputGBK string) string {
	if runtime.GOOS != "windows" {
		return outputGBK
	}

	// windows平台编码为gbk，需转换为utf8才能入库
	outputUTF8, err := gcharset.ToUTF8("GBK", outputGBK)
	if err != nil {
		return outputGBK
	}
	return outputUTF8
}

// GenerateRandomJWTSecret 生成随机的JWT密钥
func GenerateRandomJWTSecret() string {
	// 生成64字节的随机数据
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备方案
		zap.S().Warnf("生成随机JWT密钥失败，使用时间戳作为后备方案: %v", err)
		return fmt.Sprintf("jwt_secret_fallback_%d", time.Now().UnixNano())
	}

	// 将字节转换为十六进制字符串
	return hex.EncodeToString(bytes)
}
