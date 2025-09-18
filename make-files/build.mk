# build.mk 仅负责构建与打包相关的目标

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

.PHONY: frontend-build
frontend-build: ## 构建 web-admin 前端产物
	@echo "构建 web-admin 前端产物..."
	$(FRONTEND_BUILD_CMD)

.PHONY: clean
clean: ## 清理构建产物
	@echo "清理构建产物..."
	rm -rf $(BUILD_DIR)
	rm -f $(BACKEND_DIR)/coverage.out $(BACKEND_DIR)/coverage.html

.PHONY: dist-clean
dist-clean: ## 清理 dist 目录
	@echo "清理 dist 目录..."
	rm -rf $(DIST_ROOT)

.PHONY: release-clean
release-clean: ## 清理 target 目录
	@echo "清理 target 目录..."
	rm -rf $(RELEASE_ROOT)

.PHONY: dist
dist: build frontend-build ## 生成供 install/local 使用的目录结构
	@echo "准备目录模式安装包..."
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_BACKEND_BIN) $(DIST_WEBADMIN_DIR)
	@cp plugin.yaml $(DIST_DIR)/
	@cp $(BUILD_DIR)/plugin $(DIST_BACKEND_BIN)/
	@if [ -d "$(FRONTEND_OUTPUT)" ] && [ -n "$$(ls -A $(FRONTEND_OUTPUT) 2>/dev/null)" ]; then \
			echo "复制前端构建产物 -> $(DIST_WEBADMIN_OUTPUT)"; \
			mkdir -p $(DIST_WEBADMIN_OUTPUT); \
			cp -R $(FRONTEND_OUTPUT)/. $(DIST_WEBADMIN_OUTPUT)/; \
	else \
			echo "提示: 未找到或前端构建目录为空 ($(FRONTEND_OUTPUT))"; \
			echo "     请先执行 make frontend-build"; \
		fi
	@if [ -f README.md ]; then \
			cp README.md $(DIST_DIR)/; \
		fi

.PHONY: release
release: build frontend-build ## 生成 target/<version> 发布目录（包含前后端产物）
	@echo "准备发布目录 $(RELEASE_DIR)..."
	@rm -rf $(RELEASE_DIR)
	@mkdir -p $(RELEASE_BACKEND_BIN)
	@cp plugin.yaml $(RELEASE_DIR)/
	@cp $(BUILD_DIR)/plugin $(RELEASE_BACKEND_BIN)/
	@mkdir -p $(RELEASE_WEBADMIN_DIR)
	@cp -R $(FRONTEND_OUTPUT) $(RELEASE_WEBADMIN_OUTPUT)
	@if [ -f README.md ]; then \
			cp README.md $(RELEASE_DIR)/; \
		fi

.PHONY: package
package: dist ## 打包 dist 目录为 zip（用于远程安装）
	@echo "打包插件压缩包..."
	@rm -f $(APP_NAME)-$(VERSION).zip
	@cd $(DIST_ROOT) && zip -r ../$(APP_NAME)-$(VERSION).zip $(VERSION)

.PHONY: package-release
package-release: release ## 将 target/<version> 打包为 zip
	@echo "打包发布压缩包..."
	@rm -f $(APP_NAME)-$(VERSION)-release.zip
	@cd $(RELEASE_ROOT) && zip -r ../$(APP_NAME)-$(VERSION)-release.zip $(VERSION)
