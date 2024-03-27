package cron

import (
	"github.com/gogf/gf/v2"
	"cronJob/internal/models"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
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
	tasks []*models.Task
	cron *gcron.Cron
}

func (t *tTaskManager) loadOnce() {
	t.once.Do(func() {
		task := models.Task{}
		// 获取任务列表
		t.tasks = task.GetActiveTasks()

		// 初始化定时任务执行对象
		t.cron = gcron.New()

		t.cron.Add(gctx.GetInitCtx(), "*/1 * * * * *", func() {
			for _, task := range t.tasks {
				task.Run()
			}
		}, "task1")
	})
}

func (t *tTaskManager) Run() {

}

func (t *tTaskManager) Stop() {
}

func (t *tTaskManager) AddTask(task *Task) {

}

func (t *tTaskManager) DelTask(task *Task) {