# CronJob 定时任务管理系统

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-brightgreen.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-brightgreen.svg)](https://vuejs.org)

CronJob 是一个基于 Go 语言开发的现代化定时任务管理系统，提供了直观的 Web 管理界面，支持多种任务执行方式和灵活的调度策略。

## 🚀 核心特性

### 任务管理
- **多协议支持**：HTTP 请求、Shell 命令、SSH 远程执行
- **灵活调度**：支持标准 Cron 表达式，精确到秒级
- **执行策略**：并行、单次、单例、多次等多种执行策略
- **任务分组**：通过标签对任务进行分类管理
- **超时控制**：自定义任务执行超时时间

### 监控与日志
- **实时监控**：任务执行状态实时展示
- **详细日志**：完整的任务执行日志记录
- **重试机制**：自定义重试次数和间隔时间
- **执行统计**：任务执行时长、成功率统计

### 系统特性
- **用户认证**：JWT 基础的用户身份验证
- **权限管理**：基于角色的访问控制
- **RESTful API**：完整的 API 接口文档
- **响应式界面**：现代化的 Web 管理界面
- **多数据库支持**：SQLite、MySQL、PostgreSQL

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

### 数据库支持
- **SQLite** - 默认数据库（适合单机部署）
- **MySQL 8.0+** - 生产环境推荐
- **PostgreSQL 12+** - 高级功能支持

## 📦 安装部署

### 方式一：二进制部署（推荐）

1. **下载最新版本**
```bash
# 从 GitHub Releases 下载对应平台的二进制文件
wget https://github.com/your-org/cron-job/releases/download/v1.0.0/cronJob-linux-amd64-v1.0.0.tar.gz
tar -xzf cronJob-linux-amd64-v1.0.0.tar.gz
cd cronJob-linux-amd64-v1.0.0
```

2. **配置系统**
```bash
# 复制配置文件
cp config.yaml config.yaml
# 根据需要修改配置文件
vim config.yaml
```

3. **启动服务**
```bash
# 直接启动
./cronJob

# 或使用启动脚本
./start.sh
```

### 方式二：Docker 部署

1. **使用 Docker Compose**
```yaml
version: '3.8'
services:
  cronjob:
    image: cronjob:latest
    ports:
      - "8200:8200"
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data
      - ./logs:/app/logs
    environment:
      - GIN_MODE=release
    restart: unless-stopped
```

2. **启动服务**
```bash
docker-compose up -d
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
go mod download

# 安装前端依赖
cd web && pnpm install
```

3. **构建项目**
```bash
# 构建前端
make frontend

# 构建后端
make build

# 或者一键构建
make all
```

## ⚙️ 配置说明

### 配置文件结构
```yaml
# HTTP 服务配置
http:
  addr: :8200              # 监听地址
  read_timeout: 10         # 读取超时（秒）
  write_timeout: 10        # 写入超时（秒）
  max_header_bytes: 20     # 最大头部字节数（MB）

# 数据库配置
db:
  engine: sqlite           # 数据库引擎：sqlite/mysql/postgres
  name: cronJob           # 数据库名称
  data_dir: data          # 数据目录
  prefix: sched_          # 表前缀
  
# 身份验证配置
auth:
  enable: true            # 是否启用认证
  
# JWT 配置
jwt:
  secret: your-secret-key # JWT 密钥
  expires: 7200          # 过期时间（秒）

# 跨域配置
cors:
  allowed_origins: "http://localhost:8200,http://127.0.0.1:8200"

# Swagger 配置
swagger:
  title: "定时任务服务swagger API"
  desc: "这是一个简单的定时任务执行系统"
  host: "127.0.0.1:8200"
  base_path: ""
```

### 数据库配置示例

**MySQL 配置**
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

**PostgreSQL 配置**
```yaml
db:
  engine: postgres
  host: localhost
  port: 5432
  name: cronjob
  username: postgres
  password: password
  sslmode: disable
  prefix: sched_
```

## 🎯 快速开始

### 1. 首次安装

访问 `http://localhost:8200`，系统会自动进入安装模式：

1. 配置数据库连接
2. 创建管理员账户
3. 完成初始化设置

### 2. 登录系统

使用创建的管理员账户登录系统。

### 3. 创建任务

#### HTTP 任务示例
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

#### Shell 任务示例
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

#### SSH 任务示例
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

## 📖 API 文档

### 接口概览

系统提供完整的 RESTful API，支持所有 Web 界面功能：

- **用户管理**：登录、注册、密码修改
- **任务管理**：创建、编辑、删除、启用/禁用任务
- **任务执行**：手动触发、停止任务
- **日志查询**：任务执行日志查询和统计
- **系统监控**：系统状态和性能指标

### Swagger 文档

启动服务后，访问 `http://localhost:8200/swagger/index.html` 查看完整的 API 文档。

### 认证方式

API 使用 JWT 令牌进行身份验证：

```bash
# 获取令牌
curl -X POST http://localhost:8200/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# 使用令牌
curl -X GET http://localhost:8200/api/tasks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🔧 开发指南

### 开发环境要求

- Go 1.22+
- Node.js 18+
- pnpm
- Git

### 开发环境搭建

1. **克隆项目**
```bash
git clone https://github.com/your-org/cron-job.git
cd cron-job
```

2. **安装依赖**
```bash
make deps
```

3. **启动开发服务**
```bash
# 启动后端开发服务
make dev

# 启动前端开发服务（新终端）
cd web && pnpm dev
```

### 项目结构

```
cron-job/
├── cmd/                    # 命令行工具
├── internal/              # 内部包
│   ├── global/           # 全局常量和变量
│   ├── models/           # 数据模型
│   ├── schemas/          # 请求/响应结构
│   ├── service/          # 业务逻辑
│   └── utils/            # 工具函数
├── lib/                   # 公共库
├── web/                   # 前端代码
│   ├── src/
│   │   ├── api/          # API 接口
│   │   ├── components/   # 组件
│   │   ├── views/        # 页面
│   │   └── utils/        # 工具函数
│   └── dist/             # 构建产物
├── docs/                  # 文档
├── config.yaml           # 配置文件
├── Dockerfile            # Docker 配置
├── Makefile              # 构建脚本
└── README.md             # 项目说明
```

### 常用开发命令

```bash
# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 构建项目
make build

# 清理构建文件
make clean

# 生成 API 文档
make docs
```

## 🐳 Docker 部署

### 构建镜像

```bash
# 构建镜像
make docker

# 或手动构建
docker build -t cronjob:latest .
```

### 运行容器

```bash
# 简单运行
docker run -p 8200:8200 cronjob:latest

# 挂载配置和数据
docker run -p 8200:8200 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  cronjob:latest
```

### Docker Compose 部署

```yaml
version: '3.8'
services:
  cronjob:
    image: cronjob:latest
    ports:
      - "8200:8200"
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data
      - ./logs:/app/logs
    environment:
      - GIN_MODE=release
    restart: unless-stopped
    depends_on:
      - mysql

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: cronjob
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

## 📊 监控与运维

### 日志管理

系统使用结构化日志，日志文件位于 `logs/` 目录：

- `app.log` - 应用日志
- `error.log` - 错误日志
- `access.log` - 访问日志

### 性能监控

- 任务执行状态监控
- 系统资源使用情况
- API 响应时间统计
- 错误率统计

### 备份与恢复

```bash
# 备份数据库（SQLite）
cp data/cronJob.db data/cronJob.db.backup

# 备份配置文件
cp config.yaml config.yaml.backup

# 恢复数据库
cp data/cronJob.db.backup data/cronJob.db
```

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
- 提交前运行 `make lint` 检查代码
- 为新功能添加测试用例
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

## 🎯 开发计划

### 短期计划（v1.1.0）
- [ ] 任务依赖关系支持
- [ ] 邮件通知功能
- [ ] 任务执行历史图表
- [ ] 批量任务操作
- [ ] 任务模板功能

### 中期计划（v1.2.0）
- [ ] 集群部署支持
- [ ] 任务分组和权限管理
- [ ] Webhook 通知
- [ ] 任务执行队列优化
- [ ] 更多协议支持（gRPC、消息队列等）

### 长期计划（v2.0.0）
- [ ] 可视化任务编排
- [ ] 任务执行机器学习预测
- [ ] 多租户支持
- [ ] 插件系统
- [ ] 移动端应用

## 💬 社区支持

- 📧 邮件：support@cronjob.com
- 💬 讨论：[GitHub Discussions](https://github.com/your-org/cron-job/discussions)
- 🐛 问题：[GitHub Issues](https://github.com/your-org/cron-job/issues)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者和用户。

---

**CronJob** - 让定时任务管理变得简单高效 🚀