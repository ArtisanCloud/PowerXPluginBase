# common.mk 用于定义 Makefile 公共变量与帮助信息

# 基础变量
APP_NAME := powerx-plugin-base
VERSION := 0.1.0
PLUGIN_ID := com.powerx.plugins.base

GO_MODULE := github.com/powerx-plugins/base
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

FRONTEND_DIR := web-admin
FRONTEND_OUTPUT := $(FRONTEND_DIR)/.output
FRONTEND_BUILD_CMD ?= npm --prefix $(FRONTEND_DIR) run build

RELEASE_ROOT := target
RELEASE_DIR := $(RELEASE_ROOT)/$(VERSION)
RELEASE_BACKEND_BIN := $(RELEASE_DIR)/backend/bin
RELEASE_WEBADMIN_DIR := $(RELEASE_DIR)/web-admin
RELEASE_WEBADMIN_OUTPUT := $(RELEASE_WEBADMIN_DIR)/.output

.DEFAULT_GOAL := help

.PHONY: help
help: ## 显示可用命令列表
        @echo "PowerX Base Template Plugin Makefile"
	@echo ""
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's|^[^:]*:||' | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf " %-18s %s\n", $$1, $$2}'
