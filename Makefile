# Makefile for PowerX Scrum Plugin

# 变量定义
APP_NAME := powerx-plugin-scrum
VERSION := 0.1.0
PLUGIN_ID := com.powerx.plugins.scrum

# Go 相关变量
GO_MODULE := github.com/powerx-plugins/scrum
BACKEND_DIR := backend
BUILD_DIR := $(BACKEND_DIR)/bin
MAIN_FILE := $(BACKEND_DIR)/cmd/plugin/main.go

# Docker 相关
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_REGISTRY ?=

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "PowerX Scrum Plugin Makefile"
	@echo ""
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

.PHONY: build
build: ## 构建后端二进制文件
	@echo "构建后端二进制文件..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o bin/plugin ./cmd/plugin

.PHONY: build-linux
build-linux: ## 构建 Linux 版本的二进制文件
	@echo "构建 Linux 版本的二进制文件..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o bin/plugin ./cmd/plugin

.PHONY: run
run: ## 启动后端服务（开发模式）
	@echo "启动后端服务..."
	cd $(BACKEND_DIR) && \
	PX_BIND_ADDR=":8091" \
	PX_DB_SCHEMA="scrum" \
	PX_LOG_LEVEL="debug" \
	PX_DEV_MODE=1 \
	go run ./cmd/plugin

.PHONY: migrate
migrate: ## 运行数据库迁移
	@echo "运行数据库迁移..."
	cd $(BACKEND_DIR) && \
	PX_RUN_MIGRATE=true \
	go run ./cmd/plugin

.PHONY: migrate-cmd
migrate-cmd: ## 使用独立迁移命令
	@echo "使用独立迁移命令..."
	cd $(BACKEND_DIR) && go run ./cmd/database/migrate

.PHONY: seed
seed: ## 运行数据种子
	@echo "运行数据种子..."
	cd $(BACKEND_DIR) && go run ./cmd/database/seed

.PHONY: test
test: ## 运行测试
	@echo "运行测试..."
	cd $(BACKEND_DIR) && go test ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## 运行代码检查
	@echo "运行代码检查..."
	cd $(BACKEND_DIR) && golangci-lint run

.PHONY: fmt
fmt: ## 格式化代码
	@echo "格式化代码..."
	cd $(BACKEND_DIR) && go fmt ./...

.PHONY: mod-tidy
mod-tidy: ## 整理 Go 模块依赖
	@echo "整理 Go 模块依赖..."
	cd $(BACKEND_DIR) && go mod tidy

.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -f $(BACKEND_DIR)/coverage.out $(BACKEND_DIR)/coverage.html

.PHONY: package
package: build ## 打包插件
	@echo "打包插件..."
	@rm -f $(APP_NAME)-$(VERSION).zip
	zip -r $(APP_NAME)-$(VERSION).zip \
		plugin.yaml \
		$(BUILD_DIR)/plugin \
		web-admin/ \
		README.md \
		-x "web-admin/node_modules/*" "web-admin/.git/*"

.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run --rm -p 8091:8091 \
		-e PX_BIND_ADDR=":8091" \
		-e PX_DB_SCHEMA="scrum" \
		-e PX_LOG_LEVEL="debug" \
		$(DOCKER_IMAGE)

.PHONY: docker-push
docker-push: ## 推送 Docker 镜像
	@echo "推送 Docker 镜像..."
	@if [ -z "$(DOCKER_REGISTRY)" ]; then \
		echo "错误: DOCKER_REGISTRY 未设置"; \
		exit 1; \
	fi
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

.PHONY: dev-setup
dev-setup: ## 开发环境设置
	@echo "设置开发环境..."
	cd $(BACKEND_DIR) && go mod download
	@echo "安装 golangci-lint..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

.PHONY: check
check: lint test ## 运行所有检查

.PHONY: all
all: clean mod-tidy build test package ## 完整构建流程

# 开发便捷命令
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