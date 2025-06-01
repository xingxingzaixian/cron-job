# 多阶段构建 Dockerfile
# 第一阶段：构建前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 安装 pnpm
RUN npm install -g pnpm

# 复制前端依赖文件
COPY web/package.json web/pnpm-lock.yaml ./

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制前端源码
COPY web/ ./

# 构建前端
RUN pnpm run build

# 第二阶段：构建后端
FROM golang:1.21-alpine AS backend-builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制 Go 模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 复制前端构建结果
COPY --from=frontend-builder /app/web/dist ./web/dist

# 构建参数
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w \
    -X 'cronJob/cmd.AppVersion=${VERSION}' \
    -X 'cronJob/cmd.BuildTime=${BUILD_TIME}' \
    -X 'cronJob/cmd.GitCommit=${GIT_COMMIT}'" \
    -o cronJob .

# 第三阶段：运行时镜像
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# 复制二进制文件
COPY --from=backend-builder /app/cronJob .

# 复制配置文件
COPY config.example.yaml config.yaml

# 创建数据目录
RUN mkdir -p /app/data && \
    chown -R cronjob:cronjob /app

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./cronJob", "server"]
