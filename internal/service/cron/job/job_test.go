package job

import (
	"context"
	"cronJob/internal/models"
	"errors"
	"testing"
	"time"
)

// fakeHandler 模拟任务执行处理器
type fakeHandler struct {
	failCount int
	callCount int
	err       error
}

func (f *fakeHandler) Run(_ *models.Task, _ uint) (string, error) {
	f.callCount++
	if f.callCount <= f.failCount {
		return "", f.err
	}
	return "ok", nil
}

func TestRetryInterval(t *testing.T) {
	cases := []struct {
		name     string
		interval int16
		times    int8
		want     time.Duration
	}{
		{"配置的间隔优先", 10, 3, 10 * time.Second},
		{"未配置时按次数递增", 0, 2, 2 * time.Minute},
		{"未配置时封顶5分钟", 0, 10, 5 * time.Minute},
		{"未配置且次数为1", 0, 1, 1 * time.Minute},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryInterval(tc.interval, tc.times); got != tc.want {
				t.Fatalf("retryInterval(%d, %d) = %v, 期望 %v", tc.interval, tc.times, got, tc.want)
			}
		})
	}
}

func TestExecJobSuccess(t *testing.T) {
	h := &fakeHandler{}
	task := &models.Task{RetryTimes: 2}

	result := execJob(context.Background(), h, task, 1)
	if result.Err != nil {
		t.Fatalf("期望执行成功, 实际错误: %v", result.Err)
	}
	if result.Result != "ok" {
		t.Fatalf("期望输出 ok, 实际 %q", result.Result)
	}
	if result.RetryTimes != 0 {
		t.Fatalf("期望不重试, 实际重试 %d 次", result.RetryTimes)
	}
	if h.callCount != 1 {
		t.Fatalf("期望只执行1次, 实际 %d 次", h.callCount)
	}
}

func TestExecJobRetryThenSuccess(t *testing.T) {
	h := &fakeHandler{failCount: 1, err: errors.New("boom")}
	task := &models.Task{RetryTimes: 2, RetryInterval: 1} // 1秒重试间隔

	start := time.Now()
	result := execJob(context.Background(), h, task, 1)
	elapsed := time.Since(start)

	if result.Err != nil {
		t.Fatalf("期望重试后成功, 实际错误: %v", result.Err)
	}
	if h.callCount != 2 {
		t.Fatalf("期望执行2次, 实际 %d 次", h.callCount)
	}
	if result.RetryTimes != 1 {
		t.Fatalf("期望记录重试1次, 实际 %d", result.RetryTimes)
	}
	if elapsed < 900*time.Millisecond {
		t.Fatalf("期望至少等待一次重试间隔, 实际耗时 %v", elapsed)
	}
}

func TestExecJobCancelDuringRetry(t *testing.T) {
	h := &fakeHandler{failCount: 99, err: errors.New("boom")}
	task := &models.Task{RetryTimes: 3}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	result := execJob(ctx, h, task, 1)
	elapsed := time.Since(start)

	if result.Err == nil {
		t.Fatal("期望返回取消错误")
	}
	if elapsed > 2*time.Second {
		t.Fatalf("取消后应快速返回, 实际耗时 %v", elapsed)
	}
}

func TestExecJobExhaustRetries(t *testing.T) {
	h := &fakeHandler{failCount: 99, err: errors.New("boom")}
	task := &models.Task{RetryTimes: 0} // 只执行一次，不重试

	result := execJob(context.Background(), h, task, 1)
	if result.Err == nil {
		t.Fatal("期望返回失败错误")
	}
	if h.callCount != 1 {
		t.Fatalf("期望只执行1次, 实际 %d 次", h.callCount)
	}
}
