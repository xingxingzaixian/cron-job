# 邮件/Webhook通知功能实现说明

## 概述

本项目实现了完整的邮件/Webhook通知功能，支持在任务执行完成时发送通知。

## 功能特性

### 1. 通知类型支持
- **邮件通知**：通过SMTP服务器发送邮件通知
- **Webhook通知**：通过HTTP请求发送Webhook通知

### 2. 触发条件
- **成功**：任务执行成功时触发
- **失败**：任务执行失败时触发
- **所有**：任务执行完成时触发（无论成功或失败）
- **超时**：任务执行超时时触发
- **取消**：任务执行被取消时触发

### 3. 配置层级
- **全局配置**：在config.yaml中配置，适用于所有任务
- **任务级配置**：在数据库中配置，适用于特定任务

### 4. 重试机制
- 支持配置通知发送失败时的重试次数
- 支持配置重试间隔时间

## 数据库表设计

### 1. 通知配置表 (sched_notification_config)

```sql
CREATE TABLE sched_notification_config (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    
    task_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '任务ID，0表示全局配置',
    name VARCHAR(64) NOT NULL COMMENT '通知名称',
    type VARCHAR(20) NOT NULL COMMENT '通知类型：email/webhook',
    target VARCHAR(255) NOT NULL COMMENT '目标地址：邮箱或Webhook URL',
    trigger VARCHAR(20) NOT NULL DEFAULT 'all' COMMENT '触发条件：success/failure/all/timeout/cancel',
    enabled BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用',
    
    -- 邮件配置
    email_recipients TEXT COMMENT '收件人列表',
    email_subject VARCHAR(256) COMMENT '邮件主题模板',
    email_body TEXT COMMENT '邮件正文模板',
    
    -- Webhook配置
    webhook_url VARCHAR(512) COMMENT 'Webhook URL',
    webhook_method VARCHAR(10) NOT NULL DEFAULT 'POST' COMMENT 'HTTP方法',
    webhook_headers TEXT COMMENT '自定义请求头',
    webhook_body TEXT COMMENT '自定义请求体模板',
    
    -- 重试配置
    retry_times TINYINT NOT NULL DEFAULT 0 COMMENT '重试次数',
    retry_interval INT NOT NULL DEFAULT 5 COMMENT '重试间隔(秒)',
    
    INDEX idx_task_id (task_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2. 通知日志表 (sched_notification_log)

```sql
CREATE TABLE sched_notification_log (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    
    notification_id INT UNSIGNED NOT NULL COMMENT '通知配置ID',
    task_id INT UNSIGNED NOT NULL COMMENT '任务ID',
    task_log_id INT UNSIGNED NOT NULL COMMENT '任务日志ID',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '发送状态：1=成功, 0=失败',
    result TEXT COMMENT '发送结果/错误信息',
    retry_count TINYINT NOT NULL DEFAULT 0 COMMENT '重试次数',
    start_time BIGINT COMMENT '开始时间戳',
    end_time BIGINT COMMENT '结束时间戳',
    
    INDEX idx_notification_id (notification_id),
    INDEX idx_task_id (task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 配置文件结构

在 `config.yaml` 中添加以下配置：

```yaml
# 通知配置
notification:
    # 邮件配置
    email:
        enabled: false
        smtp_host: smtp.example.com
        smtp_port: 587
        username: your-email@example.com
        password: your-password
        from: no-reply@example.com
    
    # Webhook配置
    webhook:
        enabled: false
        url: ""
        headers: {}
    
    # 全局配置
    global:
        enabled: true
        default_trigger_condition: all
        max_retry_times: 3
        retry_interval: 5
```

## API接口

### 1. 通知配置管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/notification/list` | 获取通知配置列表 |
| GET | `/api/notification/view` | 获取单个通知配置 |
| POST | `/api/notification/create` | 创建通知配置 |
| POST | `/api/notification/update` | 更新通知配置 |
| POST | `/api/notification/delete` | 删除通知配置 |
| POST | `/api/notification/test` | 测试通知发送 |
| GET | `/api/notification/logs` | 查询通知发送日志 |

### 2. 请求参数示例

#### 创建通知配置
```json
{
    "task_id": 1,
    "name": "任务完成通知",
    "type": "email",
    "target": "user@example.com",
    "trigger": "all",
    "enabled": true,
    "email_recipients": "user1@example.com,user2@example.com",
    "email_subject": "任务执行完成: {{.TaskName}}",
    "email_body": "<h2>任务执行通知</h2><p>任务名称: {{.TaskName}}</p>",
    "retry_times": 3,
    "retry_interval": 5
}
```

#### 测试通知发送
```json
{
    "type": "email",
    "target": "test@example.com",
    "content": "这是一条测试消息"
}
```

## 代码架构

### 1. 新增文件

| 文件路径 | 说明 |
|---------|------|
| `internal/models/notification.go` | 通知配置和日志模型 |
| `internal/schemas/notification.go` | 通知相关的数据结构 |
| `internal/service/notification/service.go` | 通知核心服务 |
| `internal/service/notification/email.go` | 邮件发送服务 |
| `internal/service/notification/webhook.go` | Webhook调用服务 |
| `internal/service/notification/template.go` | 模板渲染服务 |
| `internal/service/notification/manager.go` | 通知管理器 |
| `internal/service/web/api/notification.go` | 通知配置API接口 |

### 2. 修改文件

| 文件路径 | 修改内容 |
|---------|---------|
| `internal/models/task.go` | 添加EnableNotification字段 |
| `internal/global/task.go` | 添加通知相关常量和类型 |
| `internal/service/cron/job/job.go` | 在afterExecJob中添加通知触发 |
| `internal/service/web/router/router.go` | 注册通知相关路由 |
| `lib/config/config.go` | 添加通知默认配置 |
| `lib/database/init.go` | 添加通知表自动迁移 |

## 使用示例

### 1. 配置邮件通知

1. 在config.yaml中配置SMTP服务器信息
2. 在数据库中创建通知配置，type设置为"email"
3. 任务执行完成后会自动发送邮件通知

### 2. 配置Webhook通知

1. 在数据库中创建通知配置，type设置为"webhook"
2. 配置Webhook URL和请求头
3. 任务执行完成后会自动调用Webhook

### 3. 自定义通知模板

支持使用Go模板语法自定义通知内容：

```yaml
# 邮件主题模板
email_subject: "任务执行{{if eq .Status \"success\"}}成功{{else}}失败{{end}}: {{.TaskName}}"

# 邮件正文模板
email_body: |
  <h2>任务执行通知</h2>
  <p><strong>任务名称:</strong> {{.TaskName}}</p>
  <p><strong>执行状态:</strong> {{.Status}}</p>
  <p><strong>开始时间:</strong> {{.StartTime}}</p>
  <p><strong>结束时间:</strong> {{.EndTime}}</p>
  <p><strong>执行时长:</strong> {{.Duration}}</p>
  {{if .Result}}<p><strong>执行结果:</strong> {{.Result}}</p>{{end}}
  {{if .Error}}<p><strong>错误信息:</strong> {{.Error}}</p>{{end}}
```

## 注意事项

1. **异步发送**：通知发送是异步的，不会阻塞任务执行
2. **错误处理**：通知发送失败不会影响任务状态
3. **重试机制**：支持配置重试次数和间隔
4. **安全性**：敏感配置（如SMTP密码）建议使用环境变量
5. **性能考虑**：通知配置有缓存机制，避免频繁查询数据库

## 扩展性

1. **新增通知类型**：可以轻松添加新的通知类型（如钉钉、企业微信等）
2. **自定义模板**：支持自定义邮件和Webhook模板
3. **通知日志**：完整的通知发送日志，便于排查问题
4. **统计功能**：可以扩展通知统计和监控功能
