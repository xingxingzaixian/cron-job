package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	"cronJob/internal/service/cron/handler"
	"cronJob/internal/service/cron/lib/httpclient"
	taskManager "cronJob/internal/service/cron/task_manager"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// sshPasswordMask SSH密码在接口返回中的掩码占位符
const sshPasswordMask = "********"

func TaskRegister(router *gin.RouterGroup) {
	router.GET("/list", TaskList)
	router.GET("/view", TaskView)
	router.POST("/edit", TaskEdit)
	router.POST("/op", TaskOp)
	router.POST("/create", TaskCreate)
	router.POST("/update", TaskUpdate)
	router.POST("/delete", TaskDelete)
	router.POST("/start", TaskStart)
	router.POST("/stop", TaskStop)
	router.POST("/execute", TaskExecute)
	router.POST("/test", TaskTest)

	// 任务依赖相关接口
	router.POST("/dependency/add", TaskDependencyAdd)
	router.POST("/dependency/remove", TaskDependencyRemove)
	router.GET("/dependency/list", TaskDependencyList)
}

// maskSSHParams 将SSH任务参数中的密码替换为掩码，避免明文密码通过接口泄露
func maskSSHParams(protocol global.TaskProtocol, params string) string {
	if protocol != global.TaskProtocolSSH || params == "" {
		return params
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(params), &config); err != nil {
		return params
	}

	if pw, ok := config["password"].(string); ok && pw != "" {
		config["password"] = sshPasswordMask
	}

	out, err := json.Marshal(config)
	if err != nil {
		return params
	}
	return string(out)
}

// sshPasswordNeedsPreserve 判断SSH配置中密码是否为空或为掩码，需要保留已存储的密码
func sshPasswordNeedsPreserve(params string) bool {
	if params == "" {
		return true
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(params), &config); err != nil {
		return false
	}

	pw, ok := config["password"].(string)
	return !ok || pw == "" || pw == sshPasswordMask
}

// taskDependencyWouldCycle 判断添加 taskID→dependentID 依赖是否会形成循环依赖（纯函数，便于单测）
// 从 dependentID 出发沿着"任务依赖链"（task_id → dependent_id）搜索，若能回到 taskID 则存在环
func taskDependencyWouldCycle(taskID, dependentID uint, allDeps []models.TaskDependency) bool {
	// 邻接表：task_id -> 它依赖的任务列表
	adjacency := make(map[uint][]uint)
	for _, dep := range allDeps {
		adjacency[dep.TaskID] = append(adjacency[dep.TaskID], dep.DependentID)
	}

	visited := make(map[uint]bool)
	stack := []uint{dependentID}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur == taskID {
			return true
		}
		if visited[cur] {
			continue
		}
		visited[cur] = true
		for _, next := range adjacency[cur] {
			if !visited[next] {
				stack = append(stack, next)
			}
		}
	}
	return false
}

// loadTaskDependencies 从数据库加载全部任务依赖关系
func loadTaskDependencies() ([]models.TaskDependency, error) {
	var allDeps []models.TaskDependency
	if err := global.GormDB.Find(&allDeps).Error; err != nil {
		return nil, err
	}
	return allDeps, nil
}

// GetErrorMsg 获取错误信息
func GetErrorMsg(s interface{}, err error) string {
	// 获取验证器错误信息
	getOne := func(err error) string {
		if err == nil {
			return ""
		}

		if fieldErr, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range fieldErr {
				// 返回第一个错误
				return strings.ToLower(fe.Field()) + "参数错误"
			}
		}

		return "参数错误"
	}

	return getOne(err)
}

// TaskList 任务列表
func TaskList(c *gin.Context) {
	var params schemas.SearchTaskParmas
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.SearchTaskParmas{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Errorw("参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.SearchTaskParmas{}, err),
		})
		return
	}

	task := &models.Task{}
	tasks, count, err := task.PageList(global.GormDB, &params)
	if err != nil {
		zap.S().Error("查询任务列表失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	var result []schemas.TaskItemOutput
	for _, v := range tasks {
		result = append(result, schemas.TaskItemOutput{
			TaskEditHTTPInput: schemas.TaskEditHTTPInput{
				ID:            v.ID,
				Name:          v.Name,
				Spec:          v.Spec,
				Protocol:      v.Protocol,
				Command:       v.Command,
				Params:        maskSSHParams(v.Protocol, v.Params),
				Timeout:       v.Timeout,
				Policy:        v.Policy,
				Count:         v.Count,
				Delay:         v.Delay,
				RetryTimes:    v.RetryTimes,
				RetryInterval: v.RetryInterval,
				Tag:           v.Tag,
				Remark:        v.Remark,
				Status:        v.Status,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": schemas.SearchTaskResponse{
			Total: count,
			List:  result,
		},
	})
}

// TaskView 查看任务详情
func TaskView(c *gin.Context) {
	var params schemas.TaskViewInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskViewInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Errorw("参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskViewInput{}, err),
		})
		return
	}

	task := &models.Task{}
	if err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": schemas.TaskItemOutput{
			TaskEditHTTPInput: schemas.TaskEditHTTPInput{
				ID:            task.ID,
				Name:          task.Name,
				Spec:          task.Spec,
				Protocol:      task.Protocol,
				Command:       task.Command,
				Params:        maskSSHParams(task.Protocol, task.Params),
				Timeout:       task.Timeout,
				Policy:        task.Policy,
				Count:         task.Count,
				Delay:         task.Delay,
				RetryTimes:    task.RetryTimes,
				RetryInterval: task.RetryInterval,
				Tag:           task.Tag,
				Remark:        task.Remark,
				Status:        task.Status,
			},
		},
	})
}

// TaskEdit 新增/修改任务（兼容旧版 /api/task/edit 路由）
func TaskEdit(c *gin.Context) {
	var params schemas.TaskEditHTTPInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskEditHTTPInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Errorw("参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskEditHTTPInput{}, err),
		})
		return
	}

	// ID > 0 时走更新逻辑，否则走创建逻辑
	if params.ID > 0 {
		task := &models.Task{}
		if err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID}); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "任务不存在",
			})
			return
		}

		// SSH任务密码为空或为掩码时，保留数据库中的原密码
		updateParams := params.Params
		if params.Protocol == global.TaskProtocolSSH && sshPasswordNeedsPreserve(params.Params) {
			updateParams = task.Params
		}

		_, err := task.Update(params.ID, map[string]interface{}{
			"name":           params.Name,
			"command":        params.Command,
			"params":         updateParams,
			"spec":           params.Spec,
			"protocol":       params.Protocol,
			"timeout":        params.Timeout,
			"policy":         params.Policy,
			"count":          params.Count,
			"delay":          params.Delay,
			"retry_times":    params.RetryTimes,
			"retry_interval": params.RetryInterval,
			"tag":            params.Tag,
			"remark":         params.Remark,
		})
		if err != nil {
			zap.S().Error("更新任务失败", err)
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
			return
		}

		// 如果任务的当前状态是禁用，就删除正在调度的任务；否则更新调度
		if task.Status == global.TaskStatusDisabled {
			taskManager.TaskManager.RemoveTask(task)
		} else {
			taskManager.TaskManager.UpdateTask(task)
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新成功",
		})
		return
	}

	// 创建逻辑
	task := &models.Task{}
	if ok := task.IsNameExist(params.Name); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务名已存在",
		})
		return
	}

	// SSH任务创建时必须提供真实密码（不允许掩码占位符）
	if params.Protocol == global.TaskProtocolSSH && sshPasswordNeedsPreserve(params.Params) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "SSH任务必须提供密码",
		})
		return
	}

	task.Name = params.Name
	task.Command = params.Command
	task.Params = params.Params
	task.Spec = params.Spec
	task.Protocol = params.Protocol
	task.Timeout = params.Timeout
	task.Policy = params.Policy
	task.Count = params.Count
	task.Delay = params.Delay
	task.RetryTimes = params.RetryTimes
	task.RetryInterval = params.RetryInterval
	task.Tag = params.Tag
	task.Remark = params.Remark
	task.Status = params.Status
	_, err := task.Create()
	if err != nil {
		zap.S().Error("创建任务失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 创建时状态为启用则直接注册到调度器
	if task.Status == global.TaskStatusEnabled {
		if err := taskManager.TaskManager.AddTask(task); err != nil {
			zap.S().Errorf("任务创建成功但注册调度器失败: taskId=%d, error=%v", task.ID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

// TaskOp 任务操作（兼容旧版 /api/task/op 路由）
func TaskOp(c *gin.Context) {
	var params schemas.TaskOptionInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Errorw("参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	task := &models.Task{}
	if err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	switch params.Op {
	case "start":
		if _, err := task.Update(params.ID, map[string]interface{}{"status": global.TaskStatusEnabled}); err != nil {
			zap.S().Error("更新任务状态失败", err)
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
			return
		}
		taskManager.TaskManager.StartTask(task)
	case "stop":
		if _, err := task.Update(params.ID, map[string]interface{}{"status": global.TaskStatusDisabled}); err != nil {
			zap.S().Error("更新任务状态失败", err)
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
			return
		}
		taskManager.TaskManager.StopTask(task)
	case "run":
		taskManager.TaskManager.RunTask(task)
	case "delete":
		if err := task.Delete(global.GormDB, params.ID); err != nil {
			zap.S().Error("删除任务失败", err)
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
			return
		}
		// 清理依赖关系，避免残留脏数据
		dependency := &models.TaskDependency{}
		if err := dependency.DeleteByTaskID(params.ID); err != nil {
			zap.S().Errorf("清理任务依赖失败: taskId=%d, error=%v", params.ID, err)
		}
		if err := dependency.DeleteByDependentID(params.ID); err != nil {
			zap.S().Errorf("清理被依赖关系失败: taskId=%d, error=%v", params.ID, err)
		}
		taskManager.TaskManager.RemoveTask(task)
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "不支持的操作: " + params.Op,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
	})
}

// TaskCreate 创建任务
func TaskCreate(c *gin.Context) {
	var params schemas.TaskEditHTTPInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskEditHTTPInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Errorw("参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskEditHTTPInput{}, err),
		})
		return
	}

	task := &models.Task{}
	if ok := task.IsNameExist(params.Name); ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务名已存在",
		})
		return
	}

	// SSH任务创建时必须提供真实密码（不允许掩码占位符）
	if params.Protocol == global.TaskProtocolSSH && sshPasswordNeedsPreserve(params.Params) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "SSH任务必须提供密码",
		})
		return
	}

	task.Name = params.Name
	task.Command = params.Command
	task.Params = params.Params
	task.Spec = params.Spec
	task.Protocol = params.Protocol
	task.Timeout = params.Timeout
	task.Policy = params.Policy
	task.Count = params.Count
	task.Delay = params.Delay
	task.RetryTimes = params.RetryTimes
	task.RetryInterval = params.RetryInterval
	task.Tag = params.Tag
	task.Remark = params.Remark
	task.Status = params.Status
	_, err := task.Create()
	if err != nil {
		zap.S().Error("创建任务失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 创建时状态为启用则直接注册到调度器
	if task.Status == global.TaskStatusEnabled {
		if err := taskManager.TaskManager.AddTask(task); err != nil {
			zap.S().Errorf("任务创建成功但注册调度器失败: taskId=%d, error=%v", task.ID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

// TaskUpdate 更新任务
func TaskUpdate(c *gin.Context) {
	var params schemas.TaskEditHTTPInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskEditHTTPInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskEditHTTPInput{}, err),
		})
		return
	}

	task := &models.Task{}
	err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	// SSH任务密码为空或为掩码时，保留数据库中的原密码
	updateParams := params.Params
	if params.Protocol == global.TaskProtocolSSH && sshPasswordNeedsPreserve(params.Params) {
		updateParams = task.Params
	}

	_, err = task.Update(params.ID, map[string]interface{}{
		"name":           params.Name,
		"command":        params.Command,
		"params":         updateParams,
		"spec":           params.Spec,
		"protocol":       params.Protocol,
		"timeout":        params.Timeout,
		"policy":         params.Policy,
		"count":          params.Count,
		"delay":          params.Delay,
		"retry_times":    params.RetryTimes,
		"retry_interval": params.RetryInterval,
		"tag":            params.Tag,
		"remark":         params.Remark,
	})
	if err != nil {
		zap.S().Error("更新任务失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 如果更新的任务状态是禁用，就删除当前正在调度的任务
	if task.Status == global.TaskStatusDisabled {
		taskManager.TaskManager.RemoveTask(task)
	} else {
		// 添加任务到调度进程中
		taskManager.TaskManager.UpdateTask(task)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// TaskDelete 删除任务
func TaskDelete(c *gin.Context) {
	var params schemas.TaskOptionInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	task := &models.Task{}
	err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	err = task.Delete(global.GormDB, params.ID)
	if err != nil {
		zap.S().Error("删除任务失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 清理依赖关系，避免残留脏数据
	dependency := &models.TaskDependency{}
	if err := dependency.DeleteByTaskID(params.ID); err != nil {
		zap.S().Errorf("清理任务依赖失败: taskId=%d, error=%v", params.ID, err)
	}
	if err := dependency.DeleteByDependentID(params.ID); err != nil {
		zap.S().Errorf("清理被依赖关系失败: taskId=%d, error=%v", params.ID, err)
	}

	// 从调度进程中移除任务
	taskManager.TaskManager.RemoveTask(task)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// TaskStart 启动任务
func TaskStart(c *gin.Context) {
	var params schemas.TaskOptionInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	task := &models.Task{}
	err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	_, err = task.Update(params.ID, map[string]interface{}{"status": global.TaskStatusEnabled})
	if err != nil {
		zap.S().Error("更新任务状态失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	taskManager.TaskManager.StartTask(task)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "启动成功",
	})
}

// TaskStop 停止任务
func TaskStop(c *gin.Context) {
	var params schemas.TaskOptionInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	task := &models.Task{}
	err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	_, err = task.Update(params.ID, map[string]interface{}{"status": global.TaskStatusDisabled})
	if err != nil {
		zap.S().Error("更新任务状态失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	taskManager.TaskManager.StopTask(task)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "停止成功",
	})
}

// TaskExecute 执行任务
func TaskExecute(c *gin.Context) {
	var params schemas.TaskOptionInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskOptionInput{}, err),
		})
		return
	}

	task := &models.Task{}
	err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	taskManager.TaskManager.RunTask(task)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "执行成功",
	})
}

// TaskTest 测试任务配置
// @Summary 测试任务配置
// @Description 按当前表单配置同步执行一次，不保存任务
// @Tags 任务
// @Accept json
// @Produce json
// @Param data body schemas.TaskTestInput true "测试任务配置"
// @Success 200 {object} schemas.Response{data=schemas.TaskTestOutput} "success"
// @Router /api/task/test [post]
func TaskTest(c *gin.Context) {
	params := &schemas.TaskTestInput{}
	if err := c.ShouldBind(params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskTestInput{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Errorw("参数验证失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskTestInput{}, err),
		})
		return
	}

	h, err := newTaskHandler(params.Protocol)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	start := time.Now()
	result := schemas.TaskTestOutput{}

	if params.Protocol == global.TaskProtocolHttp {
		// HTTP 任务：直接执行以获取状态码、耗时、响应大小
		var cmd struct {
			URL    string `json:"url"`
			Method string `json:"method"`
		}
		if err := json.Unmarshal([]byte(params.Command), &cmd); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "command 解析失败",
			})
			return
		}
		timeout := params.Timeout
		if timeout <= 0 {
			timeout = global.HttpExecTimeout
		}
		httpResult, runErr := httpclient.Do(cmd.Method, cmd.URL, params.Params, time.Duration(timeout)*time.Second)
		if httpResult != nil {
			result.Output = httpResult.Output
			result.StatusCode = httpResult.StatusCode
			result.Size = httpResult.Size
		}
		result.DurationMs = time.Since(start).Milliseconds()
		result.Success = runErr == nil
		if runErr != nil {
			result.Error = runErr.Error()
		}
	} else {
		// Shell / SSH：复用任务处理器
		taskModel := &models.Task{
			Protocol: params.Protocol,
			Command:  params.Command,
			Params:   params.Params,
			Timeout:  params.Timeout,
		}
		output, runErr := h.Run(taskModel, 0)
		result.Output = output
		result.DurationMs = time.Since(start).Milliseconds()
		result.Success = runErr == nil
		if runErr != nil {
			result.Error = runErr.Error()
		}
	}
	schemas.ResponseSuccess(c, result)
}

// newTaskHandler 根据协议返回对应的任务执行处理器
func newTaskHandler(protocol global.TaskProtocol) (handler.Handler, error) {
	switch protocol {
	case global.TaskProtocolHttp:
		return &handler.HTTPHandler{}, nil
	case global.TaskProtocolShell:
		return &handler.SHELLHandler{}, nil
	case global.TaskProtocolSSH:
		return &handler.SSHHandler{}, nil
	default:
		return nil, fmt.Errorf("不支持的协议类型: %d", protocol)
	}
}

// TaskDependencyAdd 添加任务依赖
func TaskDependencyAdd(c *gin.Context) {
	var params schemas.TaskDependency
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskDependency{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskDependency{}, err),
		})
		return
	}

	// 检查任务是否存在
	task := &models.Task{}
	if err := task.FindOne(global.GormDB, map[string]interface{}{"id": params.TaskID}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "任务不存在",
		})
		return
	}

	// 检查依赖任务是否存在
	dependentTask := &models.Task{}
	if err := dependentTask.FindOne(global.GormDB, map[string]interface{}{"id": params.DependentID}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "依赖任务不存在",
		})
		return
	}

	// 检查是否试图创建循环依赖
	if params.TaskID == params.DependentID {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "不能将任务设置为依赖自己",
		})
		return
	}

	// 检查重复依赖
	var dupCount int64
	global.GormDB.Model(&models.TaskDependency{}).
		Where("task_id = ? AND dependent_id = ?", params.TaskID, params.DependentID).
		Count(&dupCount)
	if dupCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "依赖关系已存在",
		})
		return
	}

	// 深度检查循环依赖：若 dependentID 的依赖链中已经包含 taskID，则新增 taskID→dependentID 会形成环
	allDeps, err := loadTaskDependencies()
	if err != nil {
		zap.S().Error("查询依赖关系失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}
	if taskDependencyWouldCycle(params.TaskID, params.DependentID, allDeps) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "检测到循环依赖，请检查任务依赖关系",
		})
		return
	}

	// 添加依赖关系
	err = taskManager.TaskManager.AddDependency(params.TaskID, params.DependentID, params.IsMust)
	if err != nil {
		zap.S().Error("添加任务依赖失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
	})
}

// TaskDependencyRemove 移除任务依赖
func TaskDependencyRemove(c *gin.Context) {
	var params schemas.TaskDependency
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskDependency{}, err),
		})
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		zap.S().Error("参数验证失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": GetErrorMsg(schemas.TaskDependency{}, err),
		})
		return
	}

	// 移除依赖关系
	err := taskManager.TaskManager.RemoveDependency(params.TaskID, params.DependentID)
	if err != nil {
		zap.S().Error("移除任务依赖失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "移除成功",
	})
}

// TaskDependencyList 获取任务依赖列表
func TaskDependencyList(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	if taskIDStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "缺少task_id参数",
		})
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "task_id参数格式错误",
		})
		return
	}

	// 获取任务的所有依赖
	dependencies, err := taskManager.TaskManager.GetDependencies(uint(taskID))
	if err != nil {
		zap.S().Error("获取任务依赖失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 构造返回数据
	var result []schemas.TaskDependencyResponse
	for _, dep := range dependencies {
		// 获取依赖任务的名称
		dependentTask := &models.Task{}
		dependentTaskName := "任务已删除"
		if err := dependentTask.FindOne(global.GormDB, map[string]interface{}{"id": dep.DependentID}); err == nil {
			dependentTaskName = dependentTask.Name
		}

		response := schemas.TaskDependencyResponse{
			ID:                dep.ID,
			TaskID:            dep.TaskID,
			DependentID:       dep.DependentID,
			IsMust:            dep.IsMust,
			Status:            dep.Status,
			CreatedAt:         dep.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         dep.UpdatedAt.Format("2006-01-02 15:04:05"),
			DependentTaskName: dependentTaskName,
		}
		result = append(result, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    result,
	})
}
