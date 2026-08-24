package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	"cronJob/internal/service/notification"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NotificationRegister 注册通知相关路由
func NotificationRegister(router *gin.RouterGroup) {
	router.GET("/list", NotificationList)
	router.GET("/view", NotificationView)
	router.POST("/create", NotificationCreate)
	router.POST("/update", NotificationUpdate)
	router.POST("/delete", NotificationDelete)
	router.POST("/test", NotificationTest)
	router.GET("/logs", NotificationLogs)
}

// NotificationList 通知配置列表
func NotificationList(c *gin.Context) {
	var params schemas.SearchNotificationParams
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 设置默认分页参数
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 15
	}

	notificationModel := &models.NotificationConfig{}
	configs, count, err := notificationModel.PageList(global.GormDB, params.TaskID, params.PageNo, params.PageSize)
	if err != nil {
		zap.S().Error("查询通知配置列表失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 转换为输出格式
	var result []schemas.NotificationConfigOutput
	for _, config := range configs {
		result = append(result, schemas.NotificationConfigOutput{
			ID:        config.ID,
			TaskID:    config.TaskID,
			Name:      config.Name,
			Type:      config.Type,
			Target:    config.Target,
			Trigger:   config.Trigger,
			Enabled:   config.Enabled,
			CreatedAt: config.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: config.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"total": count,
			"list":  result,
		},
	})
}

// NotificationView 查看单个通知配置
func NotificationView(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "缺少通知配置ID",
		})
		return
	}

	notificationModel := &models.NotificationConfig{}
	idUint, _ := strconv.ParseUint(id, 10, 32)
	if err := notificationModel.FindByID(uint(idUint)); err != nil {
		zap.S().Error("查询通知配置失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "通知配置不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    notificationModel,
	})
}

// NotificationCreate 创建通知配置
func NotificationCreate(c *gin.Context) {
	var params schemas.NotificationConfigInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 验证参数
	if params.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "通知名称不能为空",
		})
		return
	}

	if params.Type == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "通知类型不能为空",
		})
		return
	}

	// 验证通知类型
	switch params.Type {
	case models.NotificationTypeEmail:
		if params.Target == "" && params.EmailRecipients == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "邮件收件人不能为空",
			})
			return
		}
	case models.NotificationTypeWebhook:
		if params.Target == "" && params.WebhookURL == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "Webhook URL不能为空",
			})
			return
		}
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "不支持的通知类型",
		})
		return
	}

	// 创建通知配置
	notificationModel := &models.NotificationConfig{
		TaskID:          params.TaskID,
		Name:            params.Name,
		Type:            params.Type,
		Target:          params.Target,
		Trigger:         params.Trigger,
		Enabled:         params.Enabled,
		EmailRecipients: params.EmailRecipients,
		EmailSubject:    params.EmailSubject,
		EmailBody:       params.EmailBody,
		WebhookURL:      params.WebhookURL,
		WebhookMethod:   params.WebhookMethod,
		WebhookHeaders:  params.WebhookHeaders,
		WebhookBody:     params.WebhookBody,
		RetryTimes:      params.RetryTimes,
		RetryInterval:   params.RetryInterval,
	}

	// 设置默认值
	if notificationModel.Trigger == "" {
		notificationModel.Trigger = models.NotificationTriggerAll
	}
	if notificationModel.WebhookMethod == "" {
		notificationModel.WebhookMethod = "POST"
	}

	id, err := notificationModel.Create()
	if err != nil {
		zap.S().Error("创建通知配置失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建失败",
		})
		return
	}

	// 清除通知管理器缓存
	manager := notification.GetNotificationService().GetManager()
	manager.ClearCacheByTaskID(params.TaskID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data": gin.H{
			"id": id,
		},
	})
}

// NotificationUpdate 更新通知配置
func NotificationUpdate(c *gin.Context) {
	var params schemas.NotificationConfigInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if params.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "缺少通知配置ID",
		})
		return
	}

	// 检查通知配置是否存在
	notificationModel := &models.NotificationConfig{}
	if err := notificationModel.FindByID(params.ID); err != nil {
		zap.S().Error("查询通知配置失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "通知配置不存在",
		})
		return
	}

	// 更新通知配置
	updateData := map[string]interface{}{
		"name":              params.Name,
		"type":              params.Type,
		"target":            params.Target,
		"trigger":           params.Trigger,
		"enabled":           params.Enabled,
		"email_recipients":  params.EmailRecipients,
		"email_subject":     params.EmailSubject,
		"email_body":        params.EmailBody,
		"webhook_url":       params.WebhookURL,
		"webhook_method":    params.WebhookMethod,
		"webhook_headers":   params.WebhookHeaders,
		"webhook_body":      params.WebhookBody,
		"retry_times":       params.RetryTimes,
		"retry_interval":    params.RetryInterval,
	}

	_, err := notificationModel.Update(params.ID, updateData)
	if err != nil {
		zap.S().Error("更新通知配置失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	// 清除通知管理器缓存
	manager := notification.GetNotificationService().GetManager()
	manager.ClearCacheByTaskID(notificationModel.TaskID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// NotificationDelete 删除通知配置
func NotificationDelete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "缺少通知配置ID",
		})
		return
	}

	// 检查通知配置是否存在
	notificationModel := &models.NotificationConfig{}
	idUint, _ := strconv.ParseUint(id, 10, 32)
	if err := notificationModel.FindByID(uint(idUint)); err != nil {
		zap.S().Error("查询通知配置失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "通知配置不存在",
		})
		return
	}

	// 删除通知配置
	if err := notificationModel.Delete(global.GormDB, notificationModel.ID); err != nil {
		zap.S().Error("删除通知配置失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "删除失败",
		})
		return
	}

	// 清除通知管理器缓存
	manager := notification.GetNotificationService().GetManager()
	manager.ClearCacheByTaskID(notificationModel.TaskID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// NotificationTest 测试通知发送
func NotificationTest(c *gin.Context) {
	var params schemas.NotificationTestInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 验证参数
	if params.Type == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "通知类型不能为空",
		})
		return
	}

	if params.Target == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "目标地址不能为空",
		})
		return
	}

	// 创建测试通知配置
	testConfig := &models.NotificationConfig{
		Type:    params.Type,
		Target:  params.Target,
		Trigger: models.NotificationTriggerAll, // 设置默认触发条件
	}

	// 测试通知
	notificationService := notification.GetNotificationService()
	err := notificationService.TestNotification(testConfig)
	if err != nil {
		zap.S().Error("测试通知发送失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "测试失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试成功",
	})
}

// NotificationLogs 查询通知发送日志
func NotificationLogs(c *gin.Context) {
	var params schemas.NotificationLogInput
	if err := c.ShouldBind(&params); err != nil {
		zap.S().Error("参数绑定失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if params.TaskID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "缺少任务ID",
		})
		return
	}

	// 设置默认分页参数
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 15
	}

	// 查询通知日志
	logModel := &models.NotificationLog{}
	logs, count, err := logModel.GetLogsByTaskID(params.TaskID, params.PageNo, params.PageSize)
	if err != nil {
		zap.S().Error("查询通知日志失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 转换为输出格式
	var result []schemas.NotificationLogOutput
	for _, log := range logs {
		result = append(result, schemas.NotificationLogOutput{
			ID:             log.ID,
			NotificationID: log.NotificationID,
			TaskID:         log.TaskID,
			TaskLogID:      log.TaskLogID,
			Status:         log.Status,
			Result:         log.Result,
			RetryCount:     log.RetryCount,
			StartTime:      formatTimestamp(log.StartTime),
			EndTime:        formatTimestamp(log.EndTime),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"total": count,
			"list":  result,
		},
	})
}

// formatTimestamp 格式化时间戳
func formatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
