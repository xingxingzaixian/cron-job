# CronJob Makefile
# 提供便捷的构建、测试、部署命令

# 变量定义
APP_NAME := cronJob
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -s -w \
	-X '$(APP_NAME)/cmd.AppVersion=$(VERSION)' \
	-X '$(APP_NAME)/cmd.BuildTime=$(BUILD_TIME)' \
	-X '$(APP_NAME)/cmd.GitCommit=$(GIT_COMMIT)' \
	-X '$(APP_NAME)/cmd.GoVersion=$(GO_VERSION)'

# 目录
BUILD_DIR := build
RELEASE_DIR := release
WEB_DIR := web

# 平台定义
PLATFORMS := \
	windows-amd64 \
	linux-amd64 \
	linux-arm64 \
	darwin-amd64 \
	darwin-arm64

# 默认目标
.PHONY: all
all: clean frontend build

# 显示帮助信息
.PHONY: help
help:
	@echo "CronJob 构建系统"
	@echo ""
	@echo "可用命令:"
	@echo "  make build          - 构建当前平台的二进制文件"
	@echo "  make build-all      - 交叉编译所有平台"
	@echo "  make frontend       - 构建前端"
	@echo "  make clean          - 清理构建文件"
	@echo "  make test           - 运行测试"
	@echo "  make lint           - 代码检查"
	@echo "  make dev            - 开发模式运行"
	@echo "  make docker         - 构建 Docker 镜像"
	@echo "  make release        - 创建发布包"
	@echo ""
	@echo "变量:"
	@echo "  VERSION=$(VERSION)"
	@echo "  BUILD_TIME=$(BUILD_TIME)"
	@echo "  GIT_COMMIT=$(GIT_COMMIT)"

# 构建前端
.PHONY: frontend
frontend:
	@echo "构建前端..."
	@cd $(WEB_DIR) && pnpm install --frozen-lockfile
	@cd $(WEB_DIR) && pnpm run build
	@echo "前端构建完成"

# 构建当前平台
.PHONY: build
build:
	@echo "构建 $(APP_NAME) v$(VERSION)..."
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .
	@echo "构建完成: $(APP_NAME)"

# 交叉编译所有平台
.PHONY: build-all
build-all: frontend
	@echo "开始交叉编译..."
	@mkdir -p $(RELEASE_DIR)
	@$(foreach platform,$(PLATFORMS), \
		$(call build_platform,$(platform));)
	@echo "所有平台构建完成"

# 构建单个平台的函数
define build_platform
	$(eval OS := $(word 1,$(subst -, ,$(1))))
	$(eval ARCH := $(word 2,$(subst -, ,$(1))))
	$(eval BINARY := $(APP_NAME)$(if $(filter windows,$(OS)),.exe))
	$(eval PLATFORM_DIR := $(RELEASE_DIR)/$(1))
	
	@echo "构建 $(1)..."
	@mkdir -p $(PLATFORM_DIR)
	@GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build \
		-ldflags "$(LDFLAGS)" \
		-o $(PLATFORM_DIR)/$(BINARY) .
	
	@cp config.example.yaml $(PLATFORM_DIR)/config.yaml
	@if [ -f README.md ]; then cp README.md $(PLATFORM_DIR)/; fi
	@if [ -f LICENSE ]; then cp LICENSE $(PLATFORM_DIR)/; fi
	
	@if [ "$(OS)" = "windows" ]; then \
		echo '@echo off' > $(PLATFORM_DIR)/start.bat; \
		echo 'echo 启动 CronJob 服务...' >> $(PLATFORM_DIR)/start.bat; \
		echo '$(BINARY) server' >> $(PLATFORM_DIR)/start.bat; \
		echo 'pause' >> $(PLATFORM_DIR)/start.bat; \
	else \
		echo '#!/bin/bash' > $(PLATFORM_DIR)/start.sh; \
		echo 'echo "启动 CronJob 服务..."' >> $(PLATFORM_DIR)/start.sh; \
		echo './$(BINARY) server' >> $(PLATFORM_DIR)/start.sh; \
		chmod +x $(PLATFORM_DIR)/start.sh; \
	fi
	
	@cd $(PLATFORM_DIR) && \
	if [ "$(OS)" = "windows" ]; then \
		zip -r ../$(APP_NAME)-$(1)-v$(VERSION).zip * >/dev/null 2>&1 || true; \
	else \
		tar -czf ../$(APP_NAME)-$(1)-v$(VERSION).tar.gz *; \
	fi
	
	@echo "$(1) 构建完成"
endef

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf $(RELEASE_DIR)
	@rm -f $(APP_NAME) $(APP_NAME).exe
	@rm -rf $(WEB_DIR)/dist
	@rm -rf $(WEB_DIR)/node_modules/.cache
	@go clean -cache
	@echo "清理完成"

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	@go test -v ./...

# 代码检查
.PHONY: lint
lint:
	@echo "代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过代码检查"; \
	fi

# 开发模式
.PHONY: dev
dev: frontend
	@echo "启动开发模式..."
	@go run . server

# 构建 Docker 镜像
.PHONY: docker
docker:
	@echo "构建 Docker 镜像..."
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "Docker 镜像构建完成"

# 创建发布包
.PHONY: release
release: clean build-all
	@echo "创建发布包..."
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git 提交: $(GIT_COMMIT)"
	@echo ""
	@echo "生成的文件:"
	@ls -la $(RELEASE_DIR)/*.{zip,tar.gz} 2>/dev/null || true
	@echo ""
	@echo "发布包创建完成"

# 安装依赖
.PHONY: deps
deps:
	@echo "安装 Go 依赖..."
	@go mod download
	@go mod tidy
	@echo "安装前端依赖..."
	@cd $(WEB_DIR) && pnpm install

# 更新依赖
.PHONY: update
update:
	@echo "更新 Go 依赖..."
	@go get -u ./...
	@go mod tidy
	@echo "更新前端依赖..."
	@cd $(WEB_DIR) && pnpm update

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化 Go 代码..."
	@go fmt ./...
	@echo "格式化前端代码..."
	@cd $(WEB_DIR) && pnpm run format 2>/dev/null || true

# 生成文档
.PHONY: docs
docs:
	@echo "生成 API 文档..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init; \
	else \
		echo "swag 未安装，跳过文档生成"; \
	fi

# 显示版本信息
.PHONY: version
version:
	@echo "应用名称: $(APP_NAME)"
	@echo "版本号: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git 提交: $(GIT_COMMIT)"
	@echo "Go 版本: $(GO_VERSION)"

# 快速构建（跳过前端）
.PHONY: quick
quick:
	@echo "快速构建（跳过前端）..."
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .
	@echo "快速构建完成"

# 检查构建环境
.PHONY: check
check:
	@echo "检查构建环境..."
	@echo "Go 版本: $(GO_VERSION)"
	@command -v pnpm >/dev/null 2>&1 && echo "pnpm: 已安装" || echo "pnpm: 未安装"
	@command -v docker >/dev/null 2>&1 && echo "Docker: 已安装" || echo "Docker: 未安装"
	@command -v git >/dev/null 2>&1 && echo "Git: 已安装" || echo "Git: 未安装"
	@echo "构建环境检查完成"
