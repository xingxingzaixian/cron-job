package task

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/job"
	"cronJob/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var TaskManager = NewTaskManager()

func NewTaskManager() *tTaskManager {
	return &tTaskManager{
		once:    sync.Once{},
		running: make(map[uint]bool),
	}
}

type tTaskManager struct {
	once      sync.Once
	cron      *gcron.Cron
	runningMu sync.Mutex
	running   map[uint]bool
}

func CronServerRun() {
	TaskManager.startALL()
}

func CronServerStop() {
	TaskManager.WaitAndExit()
}

func (t *tTaskManager) startALL() {
	t.once.Do(func() {
		task := &models.Task{}
		// 获取任务列表
		taskList := task.GetActiveTasks()

		// 初始化定时任务执行对象
		t.cron = gcron.New()
		t.cron.Start()

		taskNum := 0
		for _, item := range taskList {
			if err := t.AddTask(&item); err != nil {
				zap.S().Errorf("启动时注册任务失败: taskId=%d, name=%s, error=%v", item.ID, item.Name, err)
				continue
			}
			taskNum++
		}

		zap.S().Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
	})
}

func (t *tTaskManager) RunTask(taskModel *models.Task) {
	// 同一任务同时只允许一个手动/依赖触发的执行实例，避免反复触发堆积goroutine
	if !t.tryAcquire(taskModel.ID) {
		zap.S().Warnf("任务正在执行中，跳过重复触发: taskId=%d", taskModel.ID)
		return
	}

	jobFunc := t.wrapJob(taskModel)
	if jobFunc == nil {
		t.release(taskModel.ID)
		return
	}

	// 单例策略下，若任务已有正在执行的实例（如调度触发的长任务），跳过手动/依赖触发
	if taskModel.Policy == global.TaskPolicySingle {
		var runningLog models.TaskLog
		if err := global.GormDB.Where("task_id = ?", taskModel.ID).Order("id DESC").First(&runningLog).Error; err == nil && runningLog.Status == global.TaskStatusRunning {
			t.release(taskModel.ID)
			zap.S().Infof("任务已在运行中，跳过手动执行: taskId=%d", taskModel.ID)
			return
		}
	}

	g.Go(gctx.GetInitCtx(), func(ctx context.Context) {
		defer t.release(taskModel.ID)
		jobFunc(ctx)
	}, nil)
}

// tryAcquire 尝试占用任务的执行名额，成功返回 true；已占用则返回 false
func (t *tTaskManager) tryAcquire(taskID uint) bool {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	if t.running[taskID] {
		return false
	}
	t.running[taskID] = true
	return true
}

// release 释放任务的执行名额
func (t *tTaskManager) release(taskID uint) {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	delete(t.running, taskID)
}

// wrapJob 包装任务执行函数，任务执行完成后自动检查并触发依赖该任务的下游任务
func (t *tTaskManager) wrapJob(taskModel *models.Task) gcron.JobFunc {
	baseFunc := job.CreateJob(*taskModel)
	if baseFunc == nil {
		return nil
	}
	return func(ctx context.Context) {
		baseFunc(ctx)
		t.TriggerDependents(taskModel.ID)
	}
}

// wrapScheduledJob 包装定时调度的任务执行，额外为 once/times 策略累计已执行次数
func (t *tTaskManager) wrapScheduledJob(taskModel *models.Task) gcron.JobFunc {
	baseFunc := t.wrapJob(taskModel)
	if baseFunc == nil {
		return nil
	}
	return func(ctx context.Context) {
		baseFunc(ctx)
		t.recordScheduledRun(taskModel)
	}
}

// recordScheduledRun 为 once/times 策略任务累计已执行次数（持久化到数据库，重启后不重复执行）
// 依赖未满足被跳过的执行（取消日志）不计入次数；执行完毕的任务自动禁用
func (t *tTaskManager) recordScheduledRun(taskModel *models.Task) {
	if taskModel.Policy != global.TaskPolicyOnce && taskModel.Policy != global.TaskPolicyTimes {
		return
	}

	// 依赖未满足时任务不会真正执行（日志为取消状态），不累计次数
	lastLog := &models.TaskLog{}
	if err := global.GormDB.Where("task_id = ?", taskModel.ID).Order("id DESC").First(lastLog).Error; err == nil && lastLog.Status == global.TaskStatusCancel {
		return
	}

	if err := global.GormDB.Model(&models.Task{}).
		Where("id = ?", taskModel.ID).
		Update("executed_times", gorm.Expr("executed_times + 1")).Error; err != nil {
		zap.S().Errorf("更新任务执行次数失败: taskId=%d, error=%v", taskModel.ID, err)
		return
	}

	// 次数用尽后自动禁用，避免"启用但不再调度"的状态脱节
	if taskModel.Policy == global.TaskPolicyOnce {
		_ = global.GormDB.Model(&models.Task{}).Where("id = ?", taskModel.ID).Update("status", global.TaskStatusDisabled).Error
		return
	}
	var executed int
	if err := global.GormDB.Model(&models.Task{}).Where("id = ?", taskModel.ID).Select("executed_times").Scan(&executed).Error; err == nil && executed >= taskModel.Count {
		_ = global.GormDB.Model(&models.Task{}).Where("id = ?", taskModel.ID).Update("status", global.TaskStatusDisabled).Error
	}
}

// TriggerDependents 任务执行完成后检查依赖该任务的下游任务，满足条件时自动执行
func (t *tTaskManager) TriggerDependents(taskID uint) {
	// 查询本任务最近一次执行日志，判断执行结果
	taskLog := &models.TaskLog{}
	if err := global.GormDB.Where("task_id = ?", taskID).Order("id DESC").First(taskLog).Error; err != nil {
		zap.S().Debugf("任务无执行记录，跳过下游触发: taskId=%d", taskID)
		return
	}

	dependency := &models.TaskDependency{}
	dependents, err := dependency.GetDependentsByTaskID(taskID)
	if err != nil {
		zap.S().Errorf("查询依赖任务失败: taskId=%d, error=%v", taskID, err)
		return
	}
	if len(dependents) == 0 {
		return
	}

	// 根据执行结果更新依赖状态
	switch taskLog.Status {
	case global.TaskStatusFinish:
		for _, dep := range dependents {
			if _, err := dep.Update(dep.ID, g.Map{"status": global.TaskStatusFinish}); err != nil {
				zap.S().Errorf("更新依赖状态失败: depId=%d, error=%v", dep.ID, err)
			}
		}
	case global.TaskStatusFailure:
		for _, dep := range dependents {
			if _, err := dep.Update(dep.ID, g.Map{"status": global.TaskStatusFailure}); err != nil {
				zap.S().Errorf("更新依赖状态失败: depId=%d, error=%v", dep.ID, err)
			}
		}
		return
	default:
		// 任务被取消或未完成，不触发下游
		return
	}

	// 检查每个下游任务是否满足执行条件
	for _, dep := range dependents {
		dependentTask := &models.Task{}
		if err := dependentTask.FindOne(global.GormDB, map[string]interface{}{"id": dep.TaskID}); err != nil {
			continue
		}
		// 禁用的任务不自动触发
		if dependentTask.Status == global.TaskStatusDisabled {
			continue
		}
		// 正在执行的任务不重复触发
		runningLog := &models.TaskLog{}
		if err := global.GormDB.Where("task_id = ?", dep.TaskID).Order("id DESC").First(runningLog).Error; err == nil && runningLog.Status == global.TaskStatusRunning {
			continue
		}

		canExecute, err := dependentTask.CanExecute()
		if err != nil {
			zap.S().Errorf("检查依赖任务执行条件失败: taskId=%d, error=%v", dep.TaskID, err)
			continue
		}
		if !canExecute {
			continue
		}

		zap.S().Infof("依赖任务已完成，自动触发下游任务: taskId=%d", dep.TaskID)
		t.RunTask(dependentTask)
	}
}

func (t *tTaskManager) StartTask(taskModel *models.Task) error {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		zap.S().Infof("启动定时任务-%d", taskModel.ID)
		t.cron.Start(cronName)
		return nil
	}
	return t.AddTask(taskModel)
}

func (t *tTaskManager) StopTask(taskModel *models.Task) {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		zap.S().Infof("停止定时任务-%d", taskModel.ID)
		t.cron.Stop(cronName)
	}
}

func (t *tTaskManager) AddTask(taskModel *models.Task) error {
	// 预校验cron表达式：gcron不校验字段范围且 step=0 会挂死解析器，必须提前拦截
	if err := ValidateCronSpec(taskModel.Spec); err != nil {
		return gerror.Newf("任务【%d】cron表达式无效: %v", taskModel.ID, err)
	}

	taskFunc := t.wrapScheduledJob(taskModel)
	if taskFunc == nil {
		zap.S().Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return gerror.Newf("创建任务处理Job失败,不支持的任务协议#%v", taskModel.Protocol)
	}

	cronName := strconv.Itoa(int(taskModel.ID))
	zap.S().Infof("添加定时任务-%d, 名称:%s, 策略:%v", taskModel.ID, taskModel.Name, taskModel.Policy)
	var err error
	var cronJob *gcron.Entry
	switch taskModel.Policy {
	case global.TaskPolicyMulti:
		if taskModel.Delay > 0 {
			t.cron.DelayAdd(gctx.GetInitCtx(), time.Duration(taskModel.Delay)*time.Second, taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.Add(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicySingle:
		if taskModel.Delay > 0 {
			t.cron.DelayAddSingleton(gctx.GetInitCtx(), time.Duration(taskModel.Delay)*time.Second, taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddSingleton(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicyOnce:
		// 单次任务已执行过则不再注册，避免重启后重复执行
		if taskModel.ExecutedTimes >= 1 {
			errMsg := fmt.Sprintf("单次任务已执行过(%d次)，不再重复调度", taskModel.ExecutedTimes)
			zap.S().Warnf("跳过注册单次任务: taskId=%d, %s", taskModel.ID, errMsg)
			return gerror.New(errMsg)
		}
		if taskModel.Delay > 0 {
			t.cron.DelayAddOnce(gctx.GetInitCtx(), time.Duration(taskModel.Delay)*time.Second, taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddOnce(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicyTimes:
		// 剩余次数 = 配置次数 - 已执行次数，避免重启后计数归零重复执行
		remaining := taskModel.Count - taskModel.ExecutedTimes
		if remaining <= 0 {
			errMsg := fmt.Sprintf("多次任务已执行完毕(%d/%d)，不再重复调度", taskModel.ExecutedTimes, taskModel.Count)
			zap.S().Warnf("跳过注册多次任务: taskId=%d, %s", taskModel.ID, errMsg)
			return gerror.New(errMsg)
		}
		if taskModel.Delay > 0 {
			t.cron.DelayAddTimes(gctx.GetInitCtx(), time.Duration(taskModel.Delay)*time.Second, taskModel.Spec, remaining, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddTimes(gctx.GetInitCtx(), taskModel.Spec, remaining, taskFunc, cronName)
		}
	default:
		return gerror.Newf("使用无效的策略, cron.Policy=%v", taskModel.Policy)
	}

	if taskModel.Delay == 0 {
		if err != nil {
			return err
		}

		if cronJob == nil {
			return gerror.New("添加任务失败")
		}
	}

	return nil
}

func (t *tTaskManager) RemoveTask(taskModel *models.Task) {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		zap.S().Infof("删除定时任务-%d", taskModel.ID)
		t.cron.Remove(cronName)
	}
}

func (t *tTaskManager) UpdateTask(taskModel *models.Task) error {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		t.RemoveTask(taskModel)
		return t.AddTask(taskModel)
	}
	// 任务不在调度器中（如之前注册失败），直接尝试注册
	return t.AddTask(taskModel)
}

func (t *tTaskManager) WaitAndExit() {
	if t.cron != nil {
		t.cron.Stop()
		zap.S().Info("定时任务调度器已停止")
	}
}

// MigrateSSHSecrets 将历史明文SSH密码批量加密存储（幂等：已加密或无法解密的跳过）
func MigrateSSHSecrets() (int, error) {
	var tasks []models.Task
	if err := global.GormDB.Where("protocol = ?", global.TaskProtocolSSH).Find(&tasks).Error; err != nil {
		return 0, err
	}

	migrated := 0
	for _, task := range tasks {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(task.Params), &config); err != nil {
			continue
		}
		pw, ok := config["password"].(string)
		if !ok || pw == "" || strings.HasPrefix(pw, "enc:v1:") {
			continue
		}

		enc, err := utils.EncryptSecret(pw)
		if err != nil {
			zap.S().Warnf("迁移SSH密码失败，保持原样: taskId=%d, error=%v", task.ID, err)
			continue
		}
		config["password"] = enc
		out, err := json.Marshal(config)
		if err != nil {
			continue
		}
		if err := global.GormDB.Model(&models.Task{}).Where("id = ?", task.ID).Update("params", string(out)).Error; err != nil {
			zap.S().Errorf("更新任务参数失败: taskId=%d, error=%v", task.ID, err)
			continue
		}
		migrated++
	}
	return migrated, nil
}

// AddDependency 添加任务依赖
func (t *tTaskManager) AddDependency(taskID, dependentID uint, isMust bool) error {
	dependency := &models.TaskDependency{
		TaskID:      taskID,
		DependentID: dependentID,
		IsMust:      isMust,
	}

	_, err := dependency.Create()
	return err
}

// RemoveDependency 移除任务依赖
func (t *tTaskManager) RemoveDependency(taskID, dependentID uint) error {
	dependency := &models.TaskDependency{}

	// 查找依赖关系
	var dep models.TaskDependency
	result := global.GormDB.Where("task_id = ? AND dependent_id = ?", taskID, dependentID).First(&dep)
	if result.Error != nil {
		return result.Error
	}

	// 删除依赖关系
	return dependency.Delete(dep.ID)
}

// GetDependencies 获取任务的所有依赖
func (t *tTaskManager) GetDependencies(taskID uint) ([]models.TaskDependency, error) {
	dependency := &models.TaskDependency{}
	return dependency.GetDependenciesByTaskID(taskID)
}

// GetDependents 获取依赖于指定任务的所有任务
func (t *tTaskManager) GetDependents(dependentID uint) ([]models.TaskDependency, error) {
	dependency := &models.TaskDependency{}
	return dependency.GetDependentsByTaskID(dependentID)
}

// CheckAndExecute 检查并执行任务（用于手动触发依赖任务）
func (t *tTaskManager) CheckAndExecute(taskModel *models.Task) {
	canExecute, err := taskModel.CanExecute()
	if err != nil {
		zap.S().Errorf("检查任务执行条件失败: taskId=%d, error=%v", taskModel.ID, err)
		return
	}

	if canExecute {
		zap.S().Infof("依赖条件满足，执行任务: taskId=%d", taskModel.ID)
		t.RunTask(taskModel)
	} else {
		zap.S().Infof("任务依赖未满足，跳过执行: taskId=%d", taskModel.ID)
	}
}
