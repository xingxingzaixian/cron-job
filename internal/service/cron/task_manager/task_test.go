package task

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/lib/database"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func newTestTask(id uint, policy global.TaskPolicy) *models.Task {
	return &models.Task{
		Model:    gorm.Model{ID: id},
		Name:     "test",
		Spec:     "0 */5 * * * *",
		Protocol: global.TaskProtocolHttp,
		Command:  `{"url":"http://example.com","method":"GET"}`,
		Policy:   policy,
		Count:    3,
	}
}

func TestAddTaskOnceAlreadyExecuted(t *testing.T) {
	tm := NewTaskManager()
	tm.cron = gcron.New()

	task := newTestTask(1, global.TaskPolicyOnce)
	task.ExecutedTimes = 1
	if err := tm.AddTask(task); err == nil {
		t.Fatal("单次任务已执行过时应返回错误")
	}
	if tm.cron.Search("1") != nil {
		t.Fatal("已执行过的单次任务不应注册到调度器")
	}
}

func TestAddTaskTimesExhausted(t *testing.T) {
	tm := NewTaskManager()
	tm.cron = gcron.New()

	task := newTestTask(2, global.TaskPolicyTimes)
	task.Count = 5
	task.ExecutedTimes = 5
	if err := tm.AddTask(task); err == nil {
		t.Fatal("次数用尽的多次任务应返回错误")
	}
	if tm.cron.Search("2") != nil {
		t.Fatal("次数用尽的多次任务不应注册到调度器")
	}
}

func TestAddTaskTimesRemaining(t *testing.T) {
	tm := NewTaskManager()
	tm.cron = gcron.New()
	tm.cron.Start()
	defer tm.cron.Stop()

	task := newTestTask(3, global.TaskPolicyTimes)
	task.Count = 5
	task.ExecutedTimes = 2
	if err := tm.AddTask(task); err != nil {
		t.Fatalf("剩余次数>0 时应注册成功: %v", err)
	}
	if tm.cron.Search("3") == nil {
		t.Fatal("多次任务应注册到调度器")
	}
	tm.RemoveTask(task)
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	viper.Set("db.engine", "sqlite")
	viper.Set("db.path", filepath.Join(t.TempDir(), "test_task_manager.db"))

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

func TestRecordScheduledRun(t *testing.T) {
	oldDB := global.GormDB
	defer func() { global.GormDB = oldDB }()

	db := setupTestDB(t)
	global.GormDB = db

	task := &models.Task{
		Name:     "test",
		Spec:     "0 */5 * * * *",
		Protocol: global.TaskProtocolHttp,
		Command:  `{"url":"http://example.com","method":"GET"}`,
		Policy:   global.TaskPolicyTimes,
		Count:    3,
	}
	if _, err := task.Create(); err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}

	// 正常执行（最近日志为完成）后累计次数
	now := time.Now()
	log := &models.TaskLog{
		TaskId:    task.ID,
		TaskName:  "test",
		Protocol:  global.TaskProtocolHttp,
		Status:    global.TaskStatusFinish,
		StartTime: &now,
	}
	if _, err := log.Create(); err != nil {
		t.Fatalf("创建日志失败: %v", err)
	}

	tm := NewTaskManager()
	tm.recordScheduledRun(task)

	var executed int
	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Select("executed_times").Scan(&executed).Error; err != nil {
		t.Fatalf("查询执行次数失败: %v", err)
	}
	if executed != 1 {
		t.Fatalf("期望 executed_times=1, 实际 %d", executed)
	}

	// 依赖未满足（最近日志为取消）不累计次数
	cancelLog := &models.TaskLog{
		TaskId:    task.ID,
		TaskName:  "test",
		Protocol:  global.TaskProtocolHttp,
		Status:    global.TaskStatusCancel,
		StartTime: &now,
	}
	if _, err := cancelLog.Create(); err != nil {
		t.Fatalf("创建取消日志失败: %v", err)
	}
	tm.recordScheduledRun(task)

	if err := db.Model(&models.Task{}).Where("id = ?", task.ID).Select("executed_times").Scan(&executed).Error; err != nil {
		t.Fatalf("查询执行次数失败: %v", err)
	}
	if executed != 1 {
		t.Fatalf("取消执行不应累计次数, 实际 %d", executed)
	}
}

func TestValidateCronSpec(t *testing.T) {
	cases := []struct {
		spec string
		ok   bool
	}{
		{"0 */5 * * * *", true}, // 6字段标准格式
		{"0 0 3 * * *", true},
		{"0 0 * * 1-5 *", true},     // 范围
		{"0 0 1 * JAN-MAR *", true}, // 月份名称范围
		{"# 0 0 * * *", true},       // "#"忽略秒
		{"@daily", true},            // 预定义格式
		{"@hourly", true},
		{"@every 1m", true},
		{"*/5 * * * *", false},    // 5字段，gcron不支持
		{"61 * * * * *", false},   // 秒超范围
		{"0 0 25 * * *", false},   // 小时超范围
		{"0 0 * 0 * *", false},    // 日=0 超范围
		{"0 0 * 32 * *", false},   // 日超范围
		{"0 0 * * 13 *", false},   // 月超范围
		{"0 0 * * * 7", false},    // 周=7（gcron范围0-6）
		{"*/0 * * * * *", false},  // 步长为0（防挂死解析器）
		{"1-5/2 * * * * *", true}, // 范围+步长
		{"garbage", false},
		{"", false},
	}

	for _, tc := range cases {
		err := ValidateCronSpec(tc.spec)
		if tc.ok && err != nil {
			t.Fatalf("spec %q 应通过校验, 实际错误: %v", tc.spec, err)
		}
		if !tc.ok && err == nil {
			t.Fatalf("spec %q 应校验失败", tc.spec)
		}
	}
}
