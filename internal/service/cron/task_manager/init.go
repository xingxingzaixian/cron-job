package task

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/job"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"go.uber.org/zap"
	"strconv"
	"sync"
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

func (t *tTaskManager) loadOnce() {
	t.once.Do(func() {
		task := models.Task{}
		// 获取任务列表
		taskList := task.GetActiveTasks()

		// 初始化定时任务执行对象
		t.cron = gcron.New()
		t.cron.Start()

		taskNum := 0
		for _, item := range taskList {
			t.AddTask(item)
			taskNum++
		}

		zap.S().Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
	})
}

func (t *tTaskManager) Run() {

}

func (t *tTaskManager) Stop() {
}

func (t *tTaskManager) AddTask(taskModel *models.Task) error {
	taskFunc := job.CreateJob(taskModel)
	if taskFunc == nil {
		zap.S().Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return gerror.Newf("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
	}

	cronName := strconv.Itoa(int(taskModel.ID))
	var err error
	var cronJob *gcron.Entry
	switch taskModel.Policy {
	case global.TaskPolicyMulti:
		cronJob, err = t.cron.Add(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
	case global.TaskPolicySingle:
		cronJob, err = t.cron.AddSingleton(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
	case global.TaskPolicyOnce:
		cronJob, err = t.cron.AddOnce(gctx.GetInitCtx(), taskModel.Spec, taskFunc, cronName)
	case global.TaskPolicyTimes:
		cronJob, err = t.cron.AddTimes(gctx.GetInitCtx(), taskModel.Spec, taskModel.Count, taskFunc, cronName)
	default:
		return gerror.Newf("使用无效的策略, cron.Policy=%v", taskModel.Policy)
	}

	if err != nil {
		return err
	}

	if cronJob == nil {
		return gerror.New("启动任务失败")
	}

	return nil
}

func (t *tTaskManager) DelTask(task *models.Task) {

}

func (t *tTaskManager) UpdateTask(task *models.Task) {

}
