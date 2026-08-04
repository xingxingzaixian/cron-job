package task

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/job"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"go.uber.org/zap"
)

var TaskManager = NewTaskManager()

func NewTaskManager() *tTaskManager {
	return &tTaskManager{
		once: sync.Once{},
	}
}

type tTaskManager struct {
	once sync.Once
	cron *gcron.Cron
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
			t.AddTask(&item)
			taskNum++
		}

		zap.S().Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
	})
}

func (t *tTaskManager) RunTask(taskModel *models.Task) {
	jobFunc := t.wrapJob(taskModel)
	if jobFunc != nil {
		g.Go(gctx.GetInitCtx(), jobFunc, nil)
	}
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

func (t *tTaskManager) StartTask(taskModel *models.Task) {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		zap.S().Infof("启动定时任务-%d", taskModel.ID)
		t.cron.Start(cronName)
	} else {
		t.AddTask(taskModel)
	}
}

func (t *tTaskManager) StopTask(taskModel *models.Task) {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		zap.S().Infof("停止定时任务-%d", taskModel.ID)
		t.cron.Stop(cronName)
	}
}

func (t *tTaskManager) AddTask(taskModel *models.Task) error {
	taskFunc := t.wrapJob(taskModel)
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
		if taskModel.Delay > 0 {
			t.cron.DelayAddOnce(gctx.GetInitCtx(), time.Duration(taskModel.Delay)*time.Second, taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddOnce(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicyTimes:
		if taskModel.Delay > 0 {
			t.cron.DelayAddTimes(gctx.GetInitCtx(), time.Duration(taskModel.Delay)*time.Second, taskModel.Spec, taskModel.Count, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddTimes(gctx.GetInitCtx(), taskModel.Spec, taskModel.Count, taskFunc, cronName)
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

func (t *tTaskManager) UpdateTask(taskModel *models.Task) {
	cronName := strconv.Itoa(int(taskModel.ID))
	if t.cron.Search(cronName) != nil {
		t.RemoveTask(taskModel)
		t.AddTask(taskModel)
	}
}

func (t *tTaskManager) WaitAndExit() {
	if t.cron != nil {
		t.cron.Stop()
		zap.S().Info("定时任务调度器已停止")
	}
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
