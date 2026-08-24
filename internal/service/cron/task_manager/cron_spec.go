package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
)

// cronFieldRanges 6个字段各自的合法取值范围（秒 分 时 日 月 周）
var cronFieldRanges = [6][2]int{
	{0, 59},
	{0, 59},
	{0, 23},
	{1, 31},
	{1, 12},
	{0, 6},
}

var (
	monthNameMap = map[string]int{
		"january": 1, "jan": 1,
		"february": 2, "feb": 2,
		"march": 3, "mar": 3,
		"april": 4, "apr": 4,
		"may":  5,
		"june": 6, "jun": 6,
		"july": 7, "jul": 7,
		"august": 8, "aug": 8,
		"september": 9, "sep": 9,
		"october": 10, "oct": 10,
		"november": 11, "nov": 11,
		"december": 12, "dec": 12,
	}
	weekNameMap = map[string]int{
		"sunday": 0, "sun": 0,
		"monday": 1, "mon": 1,
		"tuesday": 2, "tue": 2,
		"wednesday": 3, "wed": 3,
		"thursday": 4, "thu": 4,
		"friday": 5, "fri": 5,
		"saturday": 6, "sat": 6,
	}
)

// ValidateCronSpec 校验cron表达式
// 使用与调度器完全相同的解析器，并额外校验字段取值范围。
// 说明：gcron解析器本身不校验字段范围（如 "61 * * * * *" 会被接受但永远不会触发），
// 且步长为0（如 "*/0 ..."）会导致其解析器死循环，这里必须提前拦截。
// 支持6字段（秒 分 时 日 月 周）及 @daily/@every 等预定义格式
func ValidateCronSpec(spec string) error {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return gerror.New("cron表达式不能为空")
	}

	// 预定义格式（@daily/@every等）直接交给gcron解析器校验
	if strings.HasPrefix(spec, "@") {
		return validateCronSpecWithGcron(spec)
	}

	fields := strings.Fields(spec)
	if len(fields) != 6 {
		return gerror.Newf("cron表达式需为6个字段（秒 分 时 日 月 周），当前为%d个字段", len(fields))
	}
	// 先做取值范围与步长检查，再交给gcron解析器
	if err := validateCronFields(fields); err != nil {
		return err
	}
	return validateCronSpecWithGcron(spec)
}

// validateCronSpecWithGcron 使用gcron的解析器校验（与调度器AddTask使用同一实现）
func validateCronSpecWithGcron(spec string) error {
	cron := gcron.New()
	defer cron.Stop()
	_, err := cron.Add(gctx.GetInitCtx(), spec, func(ctx context.Context) {}, "")
	return err
}

// validateCronFields 校验每个字段的取值与步长
func validateCronFields(fields []string) error {
	for i, field := range fields {
		// gcron 用 "#" 作为秒字段占位符（忽略秒，等价5字段语义）
		if i == 0 && field == "#" {
			continue
		}
		min, max := cronFieldRanges[i][0], cronFieldRanges[i][1]
		for _, item := range strings.Split(field, ",") {
			base := item
			// 步长
			if idx := strings.Index(item, "/"); idx >= 0 {
				step, err := strconv.Atoi(item[idx+1:])
				if err != nil || step <= 0 {
					return gerror.Newf("cron表达式第%d个字段步长无效: %q", i+1, item)
				}
				base = item[:idx]
			}
			if base == "*" || base == "?" {
				continue
			}
			// 范围 n-m（允许月/周名称）
			if idx := strings.Index(base, "-"); idx >= 0 {
				from, err1 := cronFieldValue(base[:idx], i)
				to, err2 := cronFieldValue(base[idx+1:], i)
				if err1 != nil || err2 != nil {
					return gerror.Newf("cron表达式第%d个字段无效: %q", i+1, item)
				}
				if from < min || to > max || from > to {
					return gerror.Newf("cron表达式第%d个字段取值范围错误: %q（合法范围 %d-%d）", i+1, item, min, max)
				}
				continue
			}
			// 单值
			val, err := cronFieldValue(base, i)
			if err != nil || val < min || val > max {
				return gerror.Newf("cron表达式第%d个字段取值范围错误: %q（合法范围 %d-%d）", i+1, item, min, max)
			}
		}
	}
	return nil
}

// cronFieldValue 将字段值解析为数字（支持月/周名称）
func cronFieldValue(value string, fieldIndex int) (int, error) {
	if n, err := strconv.Atoi(value); err == nil {
		return n, nil
	}
	if fieldIndex == 4 { // 月
		if n, ok := monthNameMap[strings.ToLower(value)]; ok {
			return n, nil
		}
	}
	if fieldIndex == 5 { // 周
		if n, ok := weekNameMap[strings.ToLower(value)]; ok {
			return n, nil
		}
	}
	return 0, fmt.Errorf("cron字段值无效: %q", value)
}
