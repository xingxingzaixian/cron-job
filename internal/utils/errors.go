package utils

import (
	"fmt"
	"runtime"
)

// AppError 应用程序错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	File    string `json:"file,omitempty"`
	Line    int    `json:"line,omitempty"`
}

func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewAppError 创建新的应用程序错误
func NewAppError(code int, message string, detail ...string) *AppError {
	_, file, line, _ := runtime.Caller(1)
	
	err := &AppError{
		Code:    code,
		Message: message,
		File:    file,
		Line:    line,
	}
	
	if len(detail) > 0 {
		err.Detail = detail[0]
	}
	
	return err
}

// WrapError 包装现有错误
func WrapError(code int, message string, err error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  err.Error(),
		File:    file,
		Line:    line,
	}
}

// 预定义错误代码
const (
	ErrCodeInvalidParam     = 1001
	ErrCodeDatabaseError    = 1002
	ErrCodeTaskNotFound     = 1003
	ErrCodeTaskCreateFailed = 1004
	ErrCodeTaskUpdateFailed = 1005
	ErrCodeTaskDeleteFailed = 1006
	ErrCodeTaskExecFailed   = 1007
	ErrCodeUserNotFound     = 1008
	ErrCodeUserCreateFailed = 1009
	ErrCodeUserUpdateFailed = 1010
	ErrCodeUserDeleteFailed = 1011
	ErrCodeAuthFailed       = 1012
	ErrCodePermissionDenied = 1013
	ErrCodeTokenInvalid     = 1014
	ErrCodeTokenExpired     = 1015
)

// 预定义错误消息
var ErrorMessages = map[int]string{
	ErrCodeInvalidParam:     "参数无效",
	ErrCodeDatabaseError:    "数据库操作失败",
	ErrCodeTaskNotFound:     "任务不存在",
	ErrCodeTaskCreateFailed: "任务创建失败",
	ErrCodeTaskUpdateFailed: "任务更新失败",
	ErrCodeTaskDeleteFailed: "任务删除失败",
	ErrCodeTaskExecFailed:   "任务执行失败",
	ErrCodeUserNotFound:     "用户不存在",
	ErrCodeUserCreateFailed: "用户创建失败",
	ErrCodeUserUpdateFailed: "用户更新失败",
	ErrCodeUserDeleteFailed: "用户删除失败",
	ErrCodeAuthFailed:       "认证失败",
	ErrCodePermissionDenied: "权限不足",
	ErrCodeTokenInvalid:     "令牌无效",
	ErrCodeTokenExpired:     "令牌已过期",
}

// NewErrorByCode 根据错误代码创建错误
func NewErrorByCode(code int, detail ...string) *AppError {
	message, exists := ErrorMessages[code]
	if !exists {
		message = "未知错误"
	}
	return NewAppError(code, message, detail...)
}
