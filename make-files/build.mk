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
