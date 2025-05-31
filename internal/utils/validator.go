package utils

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"
)

// ValidateTaskName 验证任务名称
func ValidateTaskName(name string) bool {
	if len(name) < 2 || len(name) > 32 {
		return false
	}
	// 只允许字母、数字、中文、下划线和短横线
	// Go的正则表达式不支持\u转义，直接使用中文字符范围
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_一-龯-]+$`, name)
	return matched
}

// ValidateCronSpec 验证Cron表达式
func ValidateCronSpec(spec string) bool {
	if spec == "" {
		return false
	}
	// 简单的Cron表达式验证，支持5位和6位格式
	parts := strings.Fields(spec)
	return len(parts) == 5 || len(parts) == 6
}

// ValidateURL 验证URL格式
func ValidateURL(url string) bool {
	if url == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, url)
	return matched
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return true // 邮箱可以为空
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) bool {
	if len(password) < 6 || len(password) > 128 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	// 至少包含大写字母、小写字母、数字中的两种
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasDigit {
		count++
	}

	return count >= 2
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}
	// 只允许字母、数字、下划线，且必须以字母开头
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, username)
	return matched
}

// ValidateTag 验证标签
func ValidateTag(tag string) bool {
	if tag == "" {
		return true // 标签可以为空
	}
	if len(tag) > 32 {
		return false
	}
	// 只允许字母、数字、中文、下划线和短横线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_一-龯-]+$`, tag)
	return matched
}

// ValidateTimeout 验证超时时间
func ValidateTimeout(timeout int) bool {
	return timeout >= 0 && timeout <= 3600 // 0-3600秒
}

// ValidateRetryTimes 验证重试次数
func ValidateRetryTimes(retryTimes int8) bool {
	return retryTimes >= 0 && retryTimes <= 10
}

// ValidateRetryInterval 验证重试间隔
func ValidateRetryInterval(interval int16) bool {
	return interval >= 0 && interval <= 3600 // 0-3600秒
}

// ValidatePageSize 验证分页大小
func ValidatePageSize(pageSize int) bool {
	return pageSize > 0 && pageSize <= 100
}

// ValidatePageNo 验证页码
func ValidatePageNo(pageNo int) bool {
	return pageNo > 0
}

// SanitizeString 清理字符串，移除危险字符
func SanitizeString(input string) string {
	// 移除控制字符
	result := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, input)

	// 移除前后空白字符
	return strings.TrimSpace(result)
}

// ValidateJSONString 验证JSON字符串格式
func ValidateJSONString(jsonStr string) bool {
	if jsonStr == "" {
		return true
	}
	// 使用json包进行严格验证
	jsonStr = strings.TrimSpace(jsonStr)
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}
