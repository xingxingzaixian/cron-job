package utils

import (
	"os"
	"path/filepath"
)

// BaseDir 返回应用基础目录
// 优先使用环境变量 CRONJOB_BASE_DIR，未设置时使用当前工作目录。
// 从任意目录启动时设置该环境变量，可避免安装状态、配置和日志路径错乱。
func BaseDir() string {
	if dir := os.Getenv("CRONJOB_BASE_DIR"); dir != "" {
		return dir
	}
	return "."
}

// BasePath 拼接基础目录下的相对路径
func BasePath(name string) string {
	return filepath.Join(BaseDir(), name)
}
