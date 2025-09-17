# test.mk 汇总测试与代码质量相关目标

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

.PHONY: check
check: lint test ## 运行 lint + test
