package task

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/job"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
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
	g.Go(gctx.GetInitCtx(), job.CreateJob(*taskModel), nil)
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
	taskFunc := job.CreateJob(*taskModel)
	if taskFunc == nil {
		zap.S().Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return gerror.Newf("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
	}

	cronName := strconv.Itoa(int(taskModel.ID))
	zap.S().Infof("添加定时任务-%d", taskModel.ID, taskModel.Name, taskModel.Policy)
	var err error
	var cronJob *gcron.Entry
	switch taskModel.Policy {
	case global.TaskPolicyMulti:
		if taskModel.Delay > 0 {
			t.cron.DelayAdd(gctx.GetInitCtx(), time.Duration(taskModel.Delay), taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.Add(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicySingle:
		if taskModel.Delay > 0 {
			t.cron.DelayAddSingleton(gctx.GetInitCtx(), time.Duration(taskModel.Delay), taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddSingleton(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicyOnce:
		if taskModel.Delay > 0 {
			t.cron.DelayAddOnce(gctx.GetInitCtx(), time.Duration(taskModel.Delay), taskModel.Spec, taskFunc, cronName)
		} else {
			cronJob, err = t.cron.AddOnce(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
		}
	case global.TaskPolicyTimes:
		if taskModel.Delay > 0 {
			t.cron.DelayAddTimes(gctx.GetInitCtx(), time.Duration(taskModel.Delay), taskModel.Spec, taskModel.Count, taskFunc, cronName)
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
	t.cron.Stop()
}
