# migrate.mk 管理数据库迁移与数据初始化任务

.PHONY: migrate
migrate: ## 运行数据库迁移（通过主进程）
	@echo "运行数据库迁移..."
	cd $(BACKEND_DIR) && \
		PX_RUN_MIGRATE=true \
		go run ./cmd/plugin

.PHONY: migrate-cmd
migrate-cmd: ## 使用独立迁移命令
	@echo "使用独立迁移命令..."
	cd $(BACKEND_DIR) && go run ./cmd/database/main.go migrate

.PHONY: seed
seed: ## 运行数据种子脚本
	@echo "运行数据种子..."
	cd $(BACKEND_DIR) && go run ./cmd/database/main.go seed

.PHONY: reset-db
reset-db: ## 重置数据库（危险操作）
	@echo "警告: 这将删除所有数据！"
	@read -p "确定要继续吗？[y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "重置数据库..."; \
		cd $(BACKEND_DIR) && go run ./cmd/database/main.go refresh; \
	else \
		echo "操作已取消"; \
	fi
