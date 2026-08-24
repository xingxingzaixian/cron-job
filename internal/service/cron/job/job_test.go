package job

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/lib/database"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
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

// panicHandler 模拟执行过程中发生未预期 panic 的处理器
type panicHandler struct{}

func (p *panicHandler) Run(_ *models.Task, _ uint) (string, error) {
	panic("unexpected panic")
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

func TestExecJobPanic(t *testing.T) {
	h := &panicHandler{}
	task := &models.Task{RetryTimes: 0}

	result := execJob(context.Background(), h, task, 1)
	if result.Err == nil {
		t.Fatal("期望 panic 后返回失败结果，而不是被记为成功")
	}
	if !strings.Contains(result.Err.Error(), "内部错误") {
		t.Fatalf("错误信息应包含内部错误提示, 实际: %v", result.Err)
	}
}

func TestTruncateLogOutput(t *testing.T) {
	short := "正常输出"
	if got := truncateLogOutput(short); got != short {
		t.Fatalf("短输出不应被截断, 实际 %q", got)
	}

	exact := strings.Repeat("a", maxLogOutputSize)
	if got := truncateLogOutput(exact); got != exact {
		t.Fatal("恰好等于上限的输出不应被截断")
	}

	long := strings.Repeat("b", maxLogOutputSize+100)
	got := truncateLogOutput(long)
	if !strings.HasPrefix(got, long[:maxLogOutputSize]) {
		t.Fatal("截断后应保留前 maxLogOutputSize 字节")
	}
	if !strings.Contains(got, "输出已截断") {
		t.Fatal("截断后应包含截断提示")
	}
	if !strings.Contains(got, "1048676") {
		t.Fatalf("截断提示应包含原始长度, 实际: %q", got)
	}
}

func setupJobTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	viper.Set("db.engine", "sqlite")
	viper.Set("db.path", filepath.Join(t.TempDir(), "test_job.db"))

	sqliteDB := &database.SQLiteDB{}
	db, err := sqliteDB.Create(&gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&models.Task{}, &models.TaskLog{}); err != nil {
		t.Fatalf("迁移数据表失败: %v", err)
	}
	return db
}

func TestTaskStatusLifecycle(t *testing.T) {
	oldDB := global.GormDB
	defer func() { global.GormDB = oldDB }()

	db := setupJobTestDB(t)
	global.GormDB = db

	task := &models.Task{
		Name:     "test",
		Spec:     "0 */5 * * * *",
		Protocol: global.TaskProtocolHttp,
		Command:  `{"url":"http://example.com","method":"GET"}`,
		Policy:   global.TaskPolicyMulti,
		Status:   global.TaskStatusEnabled,
	}
	if _, err := task.Create(); err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}

	start := time.Now()
	logID := beforeExecJob(task)
	if logID == 0 {
		t.Fatal("beforeExecJob 应创建日志")
	}

	var status global.TaskStatus
	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Select("status").Scan(&status).Error; err != nil {
		t.Fatal(err)
	}
	if status != global.TaskStatusRunning {
		t.Fatalf("执行中状态应为 Running, 实际 %d", status)
	}

	afterExecJob(task, global.TaskResult{Result: "ok", Err: nil}, logID, start)

	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Select("status").Scan(&status).Error; err != nil {
		t.Fatal(err)
	}
	if status != global.TaskStatusFinish {
		t.Fatalf("执行完成后状态应为 Finish, 实际 %d", status)
	}

	// 用户手动停止（禁用）后，任务结束不应覆盖禁用状态
	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Update("status", global.TaskStatusRunning).Error; err != nil {
		t.Fatal(err)
	}
	beforeExecJob(task)
	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Update("status", global.TaskStatusDisabled).Error; err != nil {
		t.Fatal(err)
	}
	afterExecJob(task, global.TaskResult{Result: "", Err: errors.New("boom")}, logID, start)

	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Select("status").Scan(&status).Error; err != nil {
		t.Fatal(err)
	}
	if status != global.TaskStatusDisabled {
		t.Fatalf("手动停止的状态不应被任务结束覆盖, 实际 %d", status)
	}
}
