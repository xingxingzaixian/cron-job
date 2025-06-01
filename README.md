# 定时任务管理平台

基于 Golang 实现的高性能定时任务管理平台，支持秒级任务调度和分布式执行。

## ✨ 特点

- 🚀 **高性能**: 资源占用少，支持高并发任务执行
- ⏰ **秒级调度**: 支持秒级定时任务配置，精确到秒
- 📦 **部署简单**: 单二进制文件部署，仅需 YAML 配置文件
- 🎨 **管理界面**: 集成现代化前端管理页面，支持响应式设计
- 🔧 **多协议支持**: 支持 HTTP、Shell、SSH 任务类型
- 🗄️ **多数据库支持**: 支持 MySQL、PostgreSQL、SQLite 三种数据库
- 🔒 **安全可靠**: JWT 认证、CORS 配置、输入验证，支持可选认证模式
- 📊 **监控完善**: 任务执行日志、状态监控、性能指标，支持多行日志优化显示
- 🔄 **容错机制**: 任务重试、超时控制、异常恢复
- 🌐 **跨平台支持**: 支持 Windows、Linux、macOS 多平台交叉编译
- 🔐 **SSH 远程执行**: 支持密码和密钥认证，多命令执行模式
- ⚙️ **灵活配置**: 前后端统一配置管理，支持环境变量覆盖

# 界面
![](images/1.png)

![](images/2.png)

![](images/5.png)

![](images/3.png)

![](images/4.png)

![](images/6.png)

## 🏗️ 系统架构

系统采用前后端分离架构，支持多种数据库和执行器：

- **前端层**: Vue3 + Naive UI，提供现代化的用户界面
- **后端层**: Gin Web框架，提供RESTful API服务
- **业务层**: 任务管理、用户管理、日志管理等核心业务
- **执行层**: 支持HTTP、Shell、SSH三种任务执行器
- **数据层**: 支持MySQL、PostgreSQL、SQLite三种数据库

## 🚀 快速开始

### 环境要求
- Go 1.21+
- 数据库支持（任选其一）：
  - MySQL 5.7+ / 8.0+
  - PostgreSQL 12+
  - SQLite 3.x（无需额外安装，适合轻量级部署）
- Node.js 16+ 和 pnpm（前端开发）

### 安装部署

1. **克隆项目**
   ```bash
   git clone git@github.com:xingxingzaixian/cron-job.git
   cd cron-job
   ```

2. **配置数据库**

   **MySQL:**
   ```sql
   CREATE DATABASE cronJob CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

   **PostgreSQL:**
   ```sql
   CREATE DATABASE cronJob WITH ENCODING 'UTF8';
   ```

   **SQLite:**
   ```bash
   # SQLite 无需手动创建数据库，系统会自动创建
   # 默认数据库文件位置：./data/cronJob.db
   ```

3. **配置文件**
   ```bash
   cp config.example.yaml config.yaml
   # 编辑 config.yaml 配置数据库连接等信息
   ```

   **MySQL配置示例:**
   ```yaml
   db:
     engine: "mysql"
     host: "127.0.0.1"
     port: 3306
     name: "cronJob"
     user: "root"
     password: "your_password"
   ```

   **PostgreSQL配置示例:**
   ```yaml
   db:
     engine: "postgresql"
     host: "127.0.0.1"
     port: 5432
     name: "cronJob"
     user: "postgres"
     password: "your_password"
   ```

   **SQLite配置示例:**
   ```yaml
   db:
     engine: "sqlite"
     # SQLite 特有配置
     sqlite:
       file_path: "./data/cronJob.db"  # 数据库文件路径
       cache_size: -64000              # 缓存大小（KB）
       synchronous: "NORMAL"           # 同步模式：OFF/NORMAL/FULL
   ```

4. **测试数据库连接**
   ```bash
   # 运行所有数据库测试
   go test ./lib/database -v

   # 或运行特定测试
   go test ./lib/database -v -run TestDatabaseConnection
   go test ./lib/database -v -run TestDatabaseFactory
   go test ./lib/database -v -run TestDatabaseConfigValidation
   ```

5. **编译运行**
   ```bash
   go mod tidy
   go build -o cronJob main.go
   ./cronJob server
   ```

   **或使用 Makefile 构建:**
   ```bash
   # 构建当前平台
   make build

   # 交叉编译所有平台
   make build-all

   # 开发模式（包含前端构建）
   make dev
   ```

6. **访问系统**
   - 管理界面: http://localhost:8210/admin
   - API文档: http://localhost:8210/swagger/index.html
   - 默认管理员账号: admin/admin（启用认证时）

### Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o cronJob main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/cronJob .
COPY --from=builder /app/config.yaml .
CMD ["./cronJob"]
```

## 📖 使用指南

### 创建任务

#### HTTP任务
```json
{
  "name": "API健康检查",
  "spec": "0 */5 * * * *",
  "protocol": 1,
  "command": "{\"method\":\"GET\",\"url\":\"https://api.example.com/health\"}",
  "timeout": 30
}
```

#### Shell任务
```json
{
  "name": "数据备份",
  "spec": "0 0 2 * * *",
  "protocol": 2,
  "command": "mysqldump -u root -p database > backup.sql",
  "timeout": 3600
}
```

#### SSH任务
```json
{
  "name": "远程服务器维护",
  "spec": "0 0 3 * * *",
  "protocol": 3,
  "command": "systemctl status nginx\ndf -h\nfree -m",
  "params": "{\"host\":\"192.168.1.100\",\"port\":22,\"username\":\"admin\",\"password\":\"password\",\"mode\":\"sequential\"}",
  "timeout": 300
}
```

**SSH任务参数说明:**
- `host`: 目标服务器地址
- `port`: SSH端口（默认22）
- `username`: 登录用户名
- `password`: 登录密码
- `mode`: 执行模式
  - `sequential`: 顺序执行（默认），命令间共享会话上下文
  - `script`: 脚本模式，将所有命令作为一个脚本执行

### Cron表达式说明
```
秒 分 时 日 月 周
*  *  *  *  *  *
```

示例：
- `0 0 12 * * *` - 每天中午12点
- `0 */5 * * * *` - 每5分钟执行
- `0 0 0 1 * *` - 每月1号零点

## 🔧 配置说明

### 核心配置
- `http.addr`: 服务监听地址
- `db.*`: 数据库连接配置
- `auth.enable`: 是否启用用户认证（可选认证模式）
- `jwt.*`: JWT认证配置
- `cors.*`: 跨域配置

### 认证配置
```yaml
# 认证配置
auth:
  # 是否启用用户认证功能，false时系统将跳过所有用户认证和用户管理功能
  # 适用于只需要任务调度功能，不需要用户管理的场景
  enable: true

# JWT配置（仅在启用认证时生效）
jwt:
  secret: "your_jwt_secret_here"  # JWT密钥，建议使用环境变量 ${JWT_SECRET}
  expires: 7200                   # 过期时间(秒)
```

### 数据库配置
支持通过环境变量覆盖配置：
```bash
export DB_PASSWORD="your_password"
export JWT_SECRET="your_jwt_secret"
```

### 安全配置
- 使用环境变量存储敏感信息
- 配置CORS允许的域名
- 设置合适的JWT过期时间
- 支持禁用认证的轻量级部署模式

## 🛡️ 安全特性

- **可选认证模式**: 支持启用/禁用用户认证，适应不同部署场景
- **JWT身份认证**: 安全的无状态认证机制
- **输入参数验证**: 严格的参数校验和类型检查
- **SQL注入防护**: 使用参数化查询防止SQL注入
- **XSS攻击防护**: 前端输入过滤和输出编码
- **CORS跨域控制**: 可配置的跨域访问策略
- **密码强度验证**: 用户密码复杂度要求
- **SSH连接安全**: 支持密码和密钥认证方式

## 📊 监控指标

- **任务执行统计**: 成功率、失败率、执行次数
- **性能监控**: 任务执行耗时、系统资源使用
- **日志管理**: 详细的执行日志，支持多行日志优化显示
- **错误追踪**: 异常日志统计和错误分析
- **实时状态**: 任务运行状态实时监控

## 🚀 构建和部署

### 跨平台构建
项目支持多平台交叉编译：

```bash
# 使用 Makefile 构建所有平台
make build-all

# 或使用构建脚本
./build/build.sh --platforms "windows-amd64,linux-amd64,darwin-amd64"

# Windows 批处理脚本
./build/cross_build.bat
```

支持的平台：
- Windows (amd64)
- Linux (amd64, arm64)
- macOS (amd64, arm64)

### 生产部署建议

1. **数据库选择**:
   - 轻量级部署: SQLite
   - 中等规模: MySQL
   - 大规模/高并发: PostgreSQL

2. **认证模式**:
   - 内网环境: 可禁用认证 (`auth.enable: false`)
   - 公网环境: 启用认证 (`auth.enable: true`)

3. **性能优化**:
   - 配置合适的数据库连接池
   - 设置适当的任务超时时间
   - 定期清理历史日志

## 🔄 功能路线图

### 已完成 ✅
- ✅ 多数据库支持 (MySQL/PostgreSQL/SQLite)
- ✅ SSH 远程任务执行
- ✅ 可选认证模式
- ✅ 跨平台构建支持
- ✅ 用户管理和权限配置
- ✅ 任务执行日志优化显示
- ✅ 前后端配置统一管理

### 开发中 🔄
- 🔄 任务执行通知模块（邮件和webhook）
- 🔄 分布式任务执行功能
- 🔄 任务依赖关系支持

### 计划中 📋
- 📋 可视化任务流程图
- 📋 集群部署支持
- 📋 任务执行统计报表
- 📋 API 限流和熔断
- 📋 任务模板管理

## 💡 使用示例

### 场景1：轻量级内网部署
```yaml
# config.yaml
auth:
  enable: false  # 禁用认证，简化部署

db:
  engine: "sqlite"  # 使用SQLite，无需额外数据库
  sqlite:
    file_path: "./data/cronJob.db"
```

### 场景2：生产环境部署
```yaml
# config.yaml
auth:
  enable: true  # 启用认证

db:
  engine: "mysql"
  host: "${DB_HOST}"
  password: "${DB_PASSWORD}"  # 使用环境变量

jwt:
  secret: "${JWT_SECRET}"  # 使用环境变量
```

### 场景3：SSH批量服务器维护
创建SSH任务，定期检查多台服务器状态：
```bash
# 命令字段（多行命令）
systemctl status nginx
df -h /
free -m
uptime

# 参数字段（JSON格式）
{
  "host": "192.168.1.100",
  "port": 22,
  "username": "admin",
  "password": "your_password",
  "mode": "sequential"
}
```

## ❓ 常见问题

### Q: 如何选择数据库？
**A:**
- **SQLite**: 适合单机部署、数据量小于100万条记录的场景
- **MySQL**: 适合中等规模部署、需要主从复制的场景
- **PostgreSQL**: 适合大规模部署、需要复杂查询和高并发的场景

### Q: 认证模式如何选择？
**A:**
- **启用认证** (`auth.enable: true`): 适合多用户环境、公网部署
- **禁用认证** (`auth.enable: false`): 适合内网环境、单用户使用

### Q: SSH任务执行失败怎么办？
**A:**
1. 检查网络连接和SSH服务状态
2. 验证用户名密码是否正确
3. 确认目标服务器SSH配置允许密码登录
4. 检查命令语法和权限

### Q: 如何备份数据？
**A:**
- **SQLite**: 直接复制数据库文件
- **MySQL**: 使用 `mysqldump` 命令
- **PostgreSQL**: 使用 `pg_dump` 命令

### Q: 如何升级系统？
**A:**
1. 停止服务
2. 备份数据库
3. 替换二进制文件
4. 检查配置文件兼容性
5. 启动服务

## 📞 技术支持

- 🐛 **问题反馈**: [GitHub Issues](https://github.com/xingxingzaixian/cron-job/issues)
- 📖 **文档**: 
- 💬 **讨论**: 

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议。