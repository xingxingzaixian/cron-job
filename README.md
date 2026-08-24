# CronJob 定时任务管理系统

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-brightgreen.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-brightgreen.svg)](https://vuejs.org)
[![Vite Version](https://img.shields.io/badge/Vite-5.1+-brightgreen.svg)](https://vitejs.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.4+-blue.svg)](https://www.typescriptlang.org)

CronJob 是一个基于 Go 语言开发的现代化定时任务管理系统，提供了直观的 Web 管理界面，支持多种任务执行方式和灵活的调度策略。

## 📑 目录

- [🚀 简介](#-简介)
- [🌟 核心特性](#-核心特性)
- [🛠️ 技术栈](#️-技术栈)
- [📦 安装部署](#-安装部署)
  - [方式一：二进制部署（推荐）](#方式一二进制部署推荐)
  - [方式二：Docker 部署](#方式二docker-部署)
  - [方式三：源码编译](#方式三源码编译)
- [⚙️ 配置说明](#️-配置说明)
- [🎯 快速开始](#-快速开始)
- [📖 API 文档](#-api-文档)
- [🔧 开发指南](#-开发指南)
- [📊 监控与运维](#-监控与运维)
- [🤝 贡献指南](#-贡献指南)
- [📄 许可证](#-许可证)
- [🎯 开发计划](#-开发计划)

## 🚀 简介

CronJob 是一个基于 Go 语言开发的现代化定时任务管理系统，提供了直观的 Web 管理界面，支持多种任务执行方式和灵活的调度策略。

系统内置安装向导，首次启动即可通过浏览器完成配置，SQLite 开箱即用，也可切换 MySQL/PostgreSQL 用于生产环境。

## 🌟 核心特性

### 任务管理
- **多协议支持**：HTTP 请求、Shell 命令、SSH 远程执行
- **灵活调度**：支持标准 Cron 表达式，精确到秒级
- **执行策略**：并行（multi）、单例（single）、单次（once）、多次（times）四种执行策略（单次/多次的执行次数会持久化，服务重启后不会重复执行）
- **任务依赖**：支持任务间依赖关系，可配置是否强依赖
- **任务分组**：通过标签对任务进行分类管理
- **超时控制**：自定义任务执行超时时间
- **重试机制**：可配置重试次数和间隔时间
- **延迟执行**：支持任务启动延迟

### 监控与日志
- **实时监控**：任务执行状态实时展示
- **详细日志**：完整的任务执行日志记录（状态、结果、耗时）
- **执行统计**：任务执行时长、成功率统计

### 通知系统
- **邮件通知**：支持 SMTP 邮件发送，任务执行完成时自动通知
- **Webhook 通知**：支持 HTTP 回调，可对接钉钉、企业微信等
- **灵活配置**：支持任务级和全局两级配置
- **触发条件**：支持成功、失败、超时、取消等多种触发条件
- **重试机制**：通知发送失败时自动重试

### 系统特性
- **用户认证**：可选的 JWT 用户身份验证（可通过配置关闭）
- **RESTful API**：完整的 API 接口，支持 Swagger 在线文档
- **响应式界面**：现代化的 Web 管理界面，支持中英文国际化
- **多数据库支持**：SQLite（默认）、MySQL、PostgreSQL
- **CLI 管理**：基于 Cobra 的命令行工具

## 🛠️ 技术栈

### 后端技术
- **Go 1.22+** - 核心开发语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **Cobra** - CLI 工具
- **Viper** - 配置管理
- **Zap** - 日志框架
- **JWT** - 身份验证
- **Swagger** - API 文档

### 前端技术
- **Vue 3** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Naive UI** - 组件库
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Vue I18n** - 中英文国际化
- **VueUse** - 组合式 API 工具库

### 数据库支持
- **SQLite** - 默认数据库，纯 Go 驱动（modernc.org/sqlite），无需 CGO，适合单机部署
- **MySQL 8.0+** - 生产环境推荐
- **PostgreSQL 12+** - 高级功能支持

## 📦 安装部署

### 方式一：二进制部署（推荐）

1. **下载最新版本**

   从 GitHub Releases 下载对应平台的压缩包（如 `cronJob-linux-amd64-v1.0.0.tar.gz`）。

2. **解压并配置**

   ```bash
   tar -xzf cronJob-linux-amd64-v1.0.0.tar.gz
   cd cronJob-linux-amd64-v1.0.0
   ```

   项目中 `config.yaml` 为当前运行配置，部署时可将 `config.example.yaml` 复制为 `config.yaml` 后按需修改。

   > 提示：从非项目目录启动时，请设置环境变量 `CRONJOB_BASE_DIR` 指向安装目录（包含 `config.yaml` 和 `install.lock` 的目录），否则系统可能误判为未安装而进入安装模式。

3. **启动服务**

   ```bash
   ./cronJob server
   ```

   首次启动（无 `install.lock` 文件时）将自动进入安装模式，仅开放安装相关 API。完成安装后会生成 `install.lock`，服务自动切换到正常模式。

### 方式二：Docker 部署

1. **构建镜像**

   ```bash
   docker build -t cronjob:latest .
   ```

2. **运行容器**

   ```bash
   docker run -d -p 8210:8210 \
     -v $(pwd)/config.yaml:/app/config.yaml \
     -v $(pwd)/data:/app/data \
     -v $(pwd)/logs:/app/logs \
     --restart unless-stopped \
     cronjob:latest
   ```

3. **Docker Compose（含 MySQL）**

   ```yaml
   version: '3.8'
   services:
     cronjob:
       image: cronjob:latest
       ports:
         - "8210:8210"
       volumes:
         - ./config.yaml:/app/config.yaml
         - ./data:/app/data
         - ./logs:/app/logs
       environment:
         - GIN_MODE=release
       restart: unless-stopped

   ```

### 方式三：源码编译

1. **克隆项目**

   ```bash
   git clone https://github.com/your-org/cron-job.git
   cd cron-job
   ```

2. **安装依赖**

   ```bash
   # 安装 Go 依赖
   make deps

   # 或手动安装
   go mod download
   cd web && pnpm install
   ```

3. **构建项目**

   ```bash
   # 一键构建（前端 + 后端）
   make all

   # 仅构建后端
   make build

   # 仅构建前端
   make frontend

   # 交叉编译所有平台
   make build-all

   # 快速构建（跳过前端，使用已有前端资源）
   make quick
   ```

## ⚙️ 配置说明

配置文件为 `config.yaml`，以下为各字段说明：

### 完整配置示例

```yaml
# 运行模式
debug: release               # release / debug

# HTTP 服务配置
http:
    addr: :8210              # 监听地址和端口
    read_timeout: 10         # HTTP 读取超时（秒）
    write_timeout: 10        # HTTP 写入超时（秒）
    max_header_bytes: 20     # HTTP 请求头最大字节数（MB）

# Swagger 文档
swagger:
    title: "定时任务服务swagger API"
    desc: 这是一个简单的定时任务执行系统
    host: 127.0.0.1:8210
    base_path: ""
    enable: false            # 是否启用 Swagger 文档（生产环境建议关闭）

# 跨域配置
cors:
    allowed_origins: http://localhost:8210,http://127.0.0.1:8210

# 身份验证
auth:
    enable: true             # 是否启用 JWT 认证（false 时所有 API 无需登录）

# JWT 配置
jwt:
    expires: 7200            # Token 过期时间（秒），默认 2 小时
    # secret: ""             # 生产环境请固定为随机密钥（如 `openssl rand -hex 64` 生成）
                            # 不配置时每次启动生成随机密钥，重启后所有登录态失效

# 数据库配置
db:
    engine: sqlite           # 数据库引擎：sqlite / mysql / postgres
    name: cronJob            # 数据库名称
    data_dir: data           # SQLite 数据文件目录
    prefix: sched_           # 表名前缀

    # SQLite 专属配置
    sqlite:
        synchronous: NORMAL  # 同步模式：OFF / NORMAL / FULL
        cache_size: -64000   # 缓存大小（-64MB，负数表示 KB）
        temp_store: MEMORY   # 临时存储：FILE / MEMORY
        busy_timeout: 30000  # 数据库锁定等待时间（毫秒）
        vacuum_interval_hours: 168  # SQLite VACUUM 间隔（小时），回收删除后的文件空间

    # PostgreSQL 专用（engine: postgres 时生效）
    # sslmode: disable       # PostgreSQL SSL 模式：disable / require / verify-ca / verify-full
    # timezone: Asia/Shanghai # PostgreSQL 连接时区

# 任务日志管理
log:
    retention_days: 30       # 任务日志保留天数，0 表示不自动清理
    cleanup_interval_hours: 6  # 日志清理检查间隔（小时）

# 通知配置
notification:
    email:
        enabled: false               # 是否启用邮件通知
        smtp_host: ""                # SMTP 服务器地址
        smtp_port: 587               # SMTP 端口（587=STARTTLS, 465=SSL/TLS）
        username: ""                 # SMTP 用户名
        password: ""                 # SMTP 密码
        from: ""                     # 发件人地址
    webhook:
        enabled: false               # 是否启用 Webhook 通知
        url: ""                      # Webhook URL（任务级配置优先）
        headers: {}                  # 自定义请求头

# 敏感数据加密
security:
    # secret: ""             # SSH密码等敏感数据加密密钥，请固定为随机字符串（如 `openssl rand -hex 64`）
                            # 留空时回退到 jwt.secret；若两者都未固定，每次启动密钥会变化，历史密文将无法解密
```

### 数据库配置切换

**MySQL**
```yaml
db:
    engine: mysql
    host: localhost
    port: 3306
    name: cronjob
    username: root
    password: password
    charset: utf8mb4
    prefix: sched_
```

> 任务日志量较大时，MySQL 可按月分区以加快过期日志归档/删除，参考 `docs/sql/mysql_partition_log.sql`。

**PostgreSQL**
```yaml
db:
    engine: postgres
    host: localhost
    port: 5432
    name: cronjob
    username: postgres
    password: password
    sslmode: disable
    timezone: Asia/Shanghai
    prefix: sched_
```

## 🎯 快速开始

### 1. 首次安装

启动服务后访问 `http://localhost:8210`，系统会自动进入安装模式（仅开放安装相关 API）：

1. 配置数据库连接（支持测试连接）
2. 创建管理员账户
3. 完成初始化

安装完成后会生成 `install.lock` 文件，服务自动切换到正常模式，所有 API 可用。

### 2. 登录系统

使用创建的管理员账户登录。若配置文件中 `auth.enable` 为 `false`，则无需登录即可使用全部功能。

### 3. 创建任务

#### HTTP 任务

向指定 URL 发送 HTTP 请求：

```json
{
  "name": "检查服务状态",
  "spec": "0 */5 * * * *",
  "protocol": 1,
  "command": "https://api.example.com/health",
  "params": {
    "method": "GET",
    "headers": {
      "User-Agent": "CronJob/1.0"
    }
  },
  "timeout": 30,
  "retry_times": 3,
  "retry_interval": 10
}
```

#### Shell 任务

在服务器本地执行 Shell 命令：

```json
{
  "name": "数据备份",
  "spec": "0 0 2 * * *",
  "protocol": 2,
  "command": "/usr/local/bin/backup.sh",
  "params": "--full --compress",
  "timeout": 3600,
  "retry_times": 2,
  "retry_interval": 30
}
```

#### SSH 任务

远程执行 SSH 命令（支持密码认证）：

```json
{
  "name": "远程服务重启",
  "spec": "0 0 4 * * 0",
  "protocol": 3,
  "command": "systemctl restart nginx",
  "params": {
    "host": "192.168.1.100",
    "port": 22,
    "username": "admin",
    "password": "password"
  },
  "timeout": 120
}
```

### 4. 任务依赖

可以为任务设置前置依赖，使某个任务在另一个任务完成后才执行：

```json
{
  "task_id": 2,
  "dependent_id": 1,
  "is_must": true
}
```

- `is_must: true` — 强依赖，前置任务失败则本任务不执行
- `is_must: false` — 弱依赖，前置任务失败仍可执行本任务

## 📖 API 文档

### 接口概览

系统提供完整的 RESTful API，所有 Web 界面功能均可以通过 API 操作：

| 模块 | 端点 | 说明 |
|------|------|------|
| 安装 | `GET /api/install/check` | 检查安装状态 |
| | `POST /api/install/install` | 执行安装 |
| | `POST /api/install/test-db` | 测试数据库连接 |
| 认证 | `POST /api/login` | 用户登录 |
| | `GET /api/check-auth` | 检查认证是否启用 |
| 任务 | `GET /api/task/list` | 任务列表（分页） |
| | `POST /api/task/create` | 创建任务 |
| | `POST /api/task/update` | 更新任务 |
| | `POST /api/task/delete` | 删除任务 |
| | `POST /api/task/start` | 启用任务 |
| | `POST /api/task/stop` | 禁用任务 |
| | `POST /api/task/execute` | 手动执行任务 |
| 依赖 | `POST /api/task/dependency/add` | 添加依赖 |
| | `POST /api/task/dependency/remove` | 移除依赖 |
| | `GET /api/task/dependency/list` | 依赖列表 |
| 日志 | `GET /api/taskLog/list` | 执行日志（分页） |
| | `GET /api/taskLog/query` | 单条日志详情 |
| | `DELETE /api/taskLog/delete` | 删除日志 |
| 用户 | `GET /api/user/list` | 用户列表 |
| | `POST /api/user/edit` | 创建/编辑用户 |
| | `GET /api/user/view` | 用户详情 |
| | `POST /api/user/del` | 删除用户 |
| | `POST /api/user/update-password` | 修改密码 |
| | `GET /api/user/info` | 当前用户信息 |
| 通知 | `GET /api/notification/list` | 通知配置列表 |
| | `GET /api/notification/view` | 通知配置详情 |
| | `POST /api/notification/create` | 创建通知配置 |
| | `POST /api/notification/update` | 更新通知配置 |
| | `POST /api/notification/delete` | 删除通知配置 |
| | `POST /api/notification/test` | 测试通知发送 |
| | `GET /api/notification/logs` | 通知发送日志 |

### Swagger 在线文档

启动服务后，在配置中开启 `swagger.enable: true`，访问 `http://localhost:8210/swagger/index.html` 查看完整的 API 文档（默认关闭）。

### 认证方式

API 使用 JWT Bearer Token 进行身份验证：

```bash
# 获取 Token
curl -X POST http://localhost:8210/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# 使用 Token 调用 API
curl -X GET http://localhost:8210/api/task/list \
  -H "Authorization: Bearer YOUR_TOKEN"
```

> 当配置文件中 `auth.enable` 为 `false` 时，所有 API 无需认证。

## 🔧 开发指南

### 开发环境要求

- Go 1.22+
- Node.js 18+
- pnpm
- Git

### 项目结构

```
cron-job/
├── cmd/                    # CLI 入口（Cobra）
├── main.go                 # 程序入口
├── config.yaml             # 运行时配置
├── config.example.yaml     # 配置模板
├── Dockerfile              # Docker 多阶段构建
├── Makefile                # 构建脚本
├── internal/               # 内部包
│   ├── global/             # 全局变量和常量
│   ├── models/             # GORM 数据模型
│   ├── schemas/            # 请求/响应 DTO
│   ├── service/            # 业务逻辑层
│   │   ├── cron/           # 定时任务调度子系统
│   │   │   ├── handler/    # 任务执行处理器（HTTP/Shell/SSH）
│   │   │   ├── job/        # Job 创建与执行逻辑
│   │   │   ├── lib/        # 子客户端（httpclient, sshclient）
│   │   │   └── task_manager/ # 调度器编排
│   │   └── web/            # HTTP API 层
│   │       ├── api/        # 控制器
│   │       ├── middleware/  # 中间件（认证、跨域、翻译）
│   │       └── router/     # 路由注册
│   └── utils/              # 工具函数
├── lib/                    # 公共库
│   ├── config/             # Viper 配置加载
│   ├── database/           # 数据库工厂（SQLite/MySQL/PostgreSQL）
│   ├── jwt/                # JWT 工具
│   └── logger/             # Zap 日志
├── web/                    # 前端（Vue 3 + TypeScript）
│   ├── src/
│   │   ├── api/            # API 请求封装
│   │   ├── assets/         # 静态资源（SVG 图标等）
│   │   ├── components/     # 公共组件
│   │   ├── enum/           # 枚举定义
│   │   ├── hooks/          # 组合式函数（Composables）
│   │   ├── layouts/        # 布局组件（Header、Logo 等）
│   │   ├── locales/        # 国际化配置（中文/英文）
│   │   ├── router/         # 路由配置
│   │   ├── store/          # Pinia 状态管理
│   │   ├── styles/         # 全局样式（重置 CSS、过渡动画）
│   │   ├── types/          # TypeScript 类型定义
│   │   ├── utils/          # 工具函数（消息提示、本地存储）
│   │   └── views/          # 页面视图
│   ├── static.go           # Go embed 嵌入前端资源
│   └── vite.config.ts      # Vite 构建配置
├── docs/                   # Swagger 文档（自动生成）
├── data/                   # SQLite 数据文件
└── logs/                   # 应用日志
```

### 常用命令

```bash
# 安装依赖
make deps

# 一键构建（前端 + 后端）
make all

# 仅构建后端
make build

# 开发模式（构建前端后启动 Go 服务）
make dev

# 交叉编译所有平台
make build-all

# 构建 Docker 镜像
make docker

# 创建发布包
make release

# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 生成 Swagger 文档
make docs

# 清理构建产物
make clean

# 快速构建（跳过前端，使用已有前端资源）
make quick

# 检查构建环境
make check

# 显示版本信息
make version
```

### 前端开发

前后端分离开发时，可分别启动：

```bash
# 终端 1：启动后端（开发模式）
make dev

# 终端 2：启动前端 Vite 开发服务器
cd web && pnpm dev
```

Vite 开发配置中已设置代理，前端请求 `/api` 会自动转发到后端 `http://127.0.0.1:8210`。

## 📊 监控与运维

### 日志

系统使用 Uber Zap 结构化日志，日志文件位于 `logs/` 目录。

### 数据库备份

SQLite 模式下，直接备份 `data/` 目录下的数据库文件即可：

```bash
# 备份
cp data/cronJob.db data/cronJob.db.backup

# 恢复
cp data/cronJob.db.backup data/cronJob.db
```

MySQL/PostgreSQL 模式下请使用各自数据库的备份工具。

### 配置调整

- 关闭认证：将 `auth.enable` 设为 `false`
- 修改端口：调整 `http.addr`，如 `:8210`
- 跨域设置：修改 `cors.allowed_origins`

## 🤝 贡献指南

我们欢迎所有形式的贡献，包括但不限于：

- 提交 Bug 报告
- 提出新功能建议
- 提交代码补丁
- 完善文档

### 贡献流程

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发规范

- 遵循 Go 语言编码规范
- 遵循 Vue.js 开发最佳实践
- 提交前运行 `make fmt` 和 `make lint` 检查代码
- 为新功能添加相应测试
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

## 🎯 开发计划

### 已完成
- [x] HTTP / Shell / SSH 三种任务协议
- [x] Cron 表达式调度（精确到秒）
- [x] 四种执行策略（并行/单例/单次/多次）
- [x] 任务依赖关系（强依赖/弱依赖）
- [x] 重试机制（可配置次数和间隔）
- [x] 安装向导（浏览器可视化配置）
- [x] JWT 认证（可选开关）
- [x] 多数据库支持（SQLite/MySQL/PostgreSQL）
- [x] Swagger API 文档
- [x] 中英文国际化
- [x] 响应式 Web 管理界面
- [x] 快速构建命令（make quick）
- [x] 构建环境检查（make check）
- [x] 邮件/Webhook 通知（支持任务级和全局配置）

### 近期计划
- [ ] 任务执行趋势图表
- [ ] 批量任务操作
- [ ] 任务模板功能
- [ ] 更多 SSH 认证方式（私钥）

### 远期规划
- [ ] 集群部署支持
- [ ] 可视化任务编排
- [ ] 多租户支持
- [ ] 插件系统
- [ ] 移动端应用
