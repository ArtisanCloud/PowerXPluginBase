# build.mk 仅负责构建与打包相关的目标
# 固定插件 ID（你说不需要搞动态）
PLUGIN_ID           = com.powerx.plugins.base
POWERX_ADMIN_BASE   = /_p/$(PLUGIN_ID)/admin/
FRONTEND_DIR        = web-admin
FRONTEND_OUTPUT     = $(FRONTEND_DIR)/.output
CHECK_PORT          = 4100
CHECK_URL           = http://127.0.0.1:$(CHECK_PORT)$(POWERX_ADMIN_BASE)

.PHONY: build
build: ## 构建后端二进制文件
	@echo "构建后端二进制文件..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o bin/plugin ./cmd/plugin
	cd $(BACKEND_DIR) && go build -o bin/migrate ./cmd/database

.PHONY: build-linux
build-linux: ## 构建 Linux 版本的后端二进制
	@echo "构建 Linux 版本的后端二进制..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o bin/plugin ./cmd/plugin
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o bin/migrate ./cmd/database


.PHONY: frontend-build
frontend-build:
	@echo "构建 web-admin 前端产物（写入 baseURL=$(POWERX_ADMIN_BASE)）..."
	cd $(FRONTEND_DIR) && \
		POWERX_PROXY=1 \
		POWERX_ADMIN_BASE="$(POWERX_ADMIN_BASE)" \
		NODE_ENV=production \
		npx nuxi build

.PHONY: run-frontend
run-frontend: ## 仅运行已编译的 Nitro server（前台）
	@echo "启动 Nitro: $(CHECK_URL)"
	cd $(FRONTEND_OUTPUT) && \
		PORT=$(CHECK_PORT) NODE_ENV=production node server/index.mjs

.PHONY: check-base
check-base: frontend-build ## 构建后拉起 Nitro，抓首页 HTML 验证 baseURL
	@echo "启动 Nitro 并检查 baseURL..."
	cd $(FRONTEND_OUTPUT) && \
		PORT=$(CHECK_PORT) NODE_ENV=production node server/index.mjs & echo $$! > .nuxt_pid; \
		for i in $$(seq 1 40); do \
			sleep 0.25; \
			curl -fsS "$(CHECK_URL)" >/dev/null 2>&1 && break; \
			if [ $$i -eq 40 ]; then echo "Nitro 未就绪"; kill $$(cat .nuxt_pid) 2>/dev/null || true; exit 1; fi; \
		done; \
		echo -n "HTML 中的 "; \
		curl -s "$(CHECK_URL)" | grep -o 'app:{baseURL:"[^"]*"}' | head -1 || true; \
		kill $$(cat .nuxt_pid) 2>/dev/null || true; rm -f .nuxt_pid


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
	@if [ -f $(BUILD_DIR)/migrate ]; then \
		cp $(BUILD_DIR)/migrate $(DIST_BACKEND_BIN)/; \
	fi
	@if [ -d config ]; then \
		mkdir -p $(DIST_DIR)/config; \
		for f in schema.yaml values.example.yaml; do \
			if [ -f config/$$f ]; then cp config/$$f $(DIST_DIR)/config/; fi; \
		done; \
	fi
	@rm -f $(DIST_DIR)/config/host-values.yaml
	@if [ -d "$(FRONTEND_OUTPUT)" ] && [ -n "$$(ls -A $(FRONTEND_OUTPUT) 2>/dev/null)" ]; then \
			echo "复制前端构建产物 -> $(DIST_WEBADMIN_OUTPUT)"; \
			mkdir -p $(DIST_WEBADMIN_OUTPUT); \
			cp -R $(FRONTEND_OUTPUT)/. $(DIST_WEBADMIN_OUTPUT)/; \
	else \
			echo "提示: 未找到或前端构建目录为空 ($(FRONTEND_OUTPUT))"; \
			echo "     请先执行 make frontend-build"; \
		fi
	@if [ -d $(FRONTEND_DIR)/i18n ]; then \
			echo "复制前端语言资源 -> $(DIST_WEBADMIN_DIR)/i18n"; \
			mkdir -p $(DIST_WEBADMIN_DIR)/i18n; \
			cp -R $(FRONTEND_DIR)/i18n/. $(DIST_WEBADMIN_DIR)/i18n/; \
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
	@if [ -f $(BUILD_DIR)/migrate ]; then \
		cp $(BUILD_DIR)/migrate $(RELEASE_BACKEND_BIN)/; \
	fi
	@if [ -d config ]; then \
		mkdir -p $(RELEASE_DIR)/config; \
		for f in schema.yaml values.example.yaml; do \
			if [ -f config/$$f ]; then cp config/$$f $(RELEASE_DIR)/config/; fi; \
		done; \
	fi
	@rm -f $(RELEASE_DIR)/config/host-values.yaml
	@mkdir -p $(RELEASE_WEBADMIN_DIR)
	@cp -R $(FRONTEND_OUTPUT) $(RELEASE_WEBADMIN_OUTPUT)
	@if [ -d $(FRONTEND_DIR)/i18n ]; then \
			echo "复制前端语言资源 -> $(RELEASE_WEBADMIN_DIR)/i18n"; \
			mkdir -p $(RELEASE_WEBADMIN_DIR)/i18n; \
			cp -R $(FRONTEND_DIR)/i18n/. $(RELEASE_WEBADMIN_DIR)/i18n/; \
	fi
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
