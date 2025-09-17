# build.mk 负责具体的构建与打包逻辑

# 基础变量
APP_NAME := powerx-plugin-note
VERSION := 0.1.0
PLUGIN_ID := com.powerx.plugins.note

GO_MODULE := github.com/powerx-plugins/note
BACKEND_DIR := backend
BUILD_DIR := $(BACKEND_DIR)/bin
MAIN_FILE := $(BACKEND_DIR)/cmd/plugin/main.go

DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_REGISTRY ?=

DIST_ROOT := dist
DIST_DIR := $(DIST_ROOT)/$(VERSION)
DIST_BACKEND_BIN := $(DIST_DIR)/backend/bin
DIST_WEBADMIN_DIR := $(DIST_DIR)/web-admin
DIST_WEBADMIN_OUTPUT := $(DIST_WEBADMIN_DIR)/.output

# 默认目标
.PHONY: help
help: ## 显示可用命令列表
	@echo "PowerX Note Plugin Makefile"
	@echo ""
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's|^[^:]*:||' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

.PHONY: build
build: ## 构建后端二进制文件
	@echo "构建后端二进制文件..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o bin/plugin ./cmd/plugin

.PHONY: build-linux
build-linux: ## 构建 Linux 版本的后端二进制
	@echo "构建 Linux 版本的后端二进制..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o bin/plugin ./cmd/plugin

.PHONY: run
run: ## 启动后端服务（开发模式）
	@echo "启动后端服务..."
	cd $(BACKEND_DIR) && \
	PX_BIND_ADDR=":8091" \
	PX_DB_SCHEMA="note" \
	PX_LOG_LEVEL="debug" \
	PX_DEV_MODE=1 \
	go run ./cmd/plugin

.PHONY: migrate
migrate: ## 运行数据库迁移（通过主进程）
	@echo "运行数据库迁移..."
	cd $(BACKEND_DIR) && \
	PX_RUN_MIGRATE=true \
	go run ./cmd/plugin

.PHONY: migrate-cmd
migrate-cmd: ## 使用独立迁移命令
	@echo "使用独立迁移命令..."
	cd $(BACKEND_DIR) && go run ./cmd/database/migrate

.PHONY: seed
seed: ## 运行数据种子脚本
	@echo "运行数据种子..."
	cd $(BACKEND_DIR) && go run ./cmd/database/seed

.PHONY: test
test: ## 执行 Go 单元测试
	@echo "运行测试..."
	cd $(BACKEND_DIR) && go test ./...

.PHONY: test-coverage
test-coverage: ## 执行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## 运行 golangci-lint
	@echo "运行代码检查..."
	cd $(BACKEND_DIR) && golangci-lint run

.PHONY: fmt
fmt: ## 使用 go fmt 格式化代码
	@echo "格式化代码..."
	cd $(BACKEND_DIR) && go fmt ./...

.PHONY: mod-tidy
mod-tidy: ## 整理 Go 模块依赖
	@echo "整理 Go 模块依赖..."
	cd $(BACKEND_DIR) && go mod tidy

.PHONY: clean
clean: ## 清理构建产物
	@echo "清理构建产物..."
	rm -rf $(BUILD_DIR)
	rm -f $(BACKEND_DIR)/coverage.out $(BACKEND_DIR)/coverage.html

.PHONY: dist-clean
dist-clean: ## 清理 dist 目录
	@echo "清理 dist 目录..."
	rm -rf $(DIST_ROOT)

.PHONY: dist
dist: build ## 生成供 install/local 使用的目录结构
	@echo "准备目录模式安装包..."
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_BACKEND_BIN)
	@cp plugin.yaml $(DIST_DIR)/
	@cp $(BUILD_DIR)/plugin $(DIST_BACKEND_BIN)/
	@if [ -d "web-admin/.output" ]; then \
		echo "复制前端构建产物..."; \
		mkdir -p $(DIST_WEBADMIN_DIR); \
		cp -R web-admin/.output $(DIST_WEBADMIN_OUTPUT); \
	fi
	@if [ -f README.md ]; then \
		cp README.md $(DIST_DIR)/; \
	fi

.PHONY: package
package: dist ## 打包插件为 zip（用于远程安装）
	@echo "打包插件压缩包..."
	@rm -f $(APP_NAME)-$(VERSION).zip
	@cd $(DIST_ROOT) && zip -r ../$(APP_NAME)-$(VERSION).zip $(VERSION)

.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## 使用 Docker 运行插件
	@echo "运行 Docker 容器..."
	docker run --rm -p 8091:8091 \
		-e PX_BIND_ADDR=":8091" \
		-e PX_DB_SCHEMA="note" \
		-e PX_LOG_LEVEL="debug" \
		$(DOCKER_IMAGE)

.PHONY: docker-push
docker-push: ## 推送 Docker 镜像到仓库
	@echo "推送 Docker 镜像..."
	@if [ -z "$(DOCKER_REGISTRY)" ]; then \
		echo "错误: DOCKER_REGISTRY 未设置"; \
		exit 1; \
	fi
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

.PHONY: dev-setup
dev-setup: ## 初始化本地 Go 依赖与 lint 工具
	@echo "下载 Go 依赖..."
	cd $(BACKEND_DIR) && go mod download
	@echo "检测 golangci-lint..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

.PHONY: check
check: lint test ## 运行 lint + test

.PHONY: all
all: clean mod-tidy build test package ## 完整构建与打包流程

.PHONY: dev
dev: migrate run ## 迁移并启动开发服务

.PHONY: reset-db
reset-db: ## 重置数据库（危险操作）
	@echo "警告: 这将删除所有数据！"
	@read -p "确定要继续吗？[y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "重置数据库..."; \
		cd $(BACKEND_DIR) && PX_RESET_DB=true go run ./cmd/database/migrate; \
	else \
		echo "操作已取消"; \
	fi
