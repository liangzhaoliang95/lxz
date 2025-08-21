# LXZ Makefile
# 简化构建和版本管理流程

.PHONY: help build clean test release version check-update

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
help: ## 显示帮助信息
	@echo "LXZ 项目构建工具"
	@echo ""
	@echo "可用目标:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "示例:"
	@echo "  make build          # 构建当前平台"
	@echo "  make build-linux    # 构建Linux版本"
	@echo "  make build-windows  # 构建Windows版本"
	@echo "  make release        # 发布新版本"
	@echo "  make version        # 显示版本信息"

# 构建相关
build: ## 构建当前平台
	@echo "构建当前平台..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh build

build-linux: ## 构建Linux版本
	@echo "构建Linux版本..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh build linux amd64 lxz-linux-amd64

build-windows: ## 构建Windows版本
	@echo "构建Windows版本..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh build windows amd64 lxz-windows-amd64.exe

build-darwin: ## 构建macOS版本
	@echo "构建macOS版本..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh build darwin amd64 lxz-darwin-amd64

cross-build: ## 交叉编译所有平台
	@echo "交叉编译所有平台..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh cross

# 清理
clean: ## 清理构建产物
	@echo "清理构建产物..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh clean

# 测试
test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 版本管理
version: ## 显示版本信息
	@echo "版本信息:"
	@go run main.go version

check-update: ## 检查更新
	@echo "检查更新..."
	@go run main.go version --check-update

# 发布
release: ## 发布新版本
	@echo "发布新版本..."
	@chmod +x release.sh
	@./release.sh

# 开发工具
fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...

vet: ## 代码静态分析
	@echo "代码静态分析..."
	@go vet ./...

lint: ## 代码检查
	@echo "代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过代码检查"; \
	fi

# 依赖管理
deps: ## 下载依赖
	@echo "下载依赖..."
	@go mod download

deps-update: ## 更新依赖
	@echo "更新依赖..."
	@go get -u ./...
	@go mod tidy

# 安装
install: ## 安装到系统
	@echo "安装到系统..."
	@go install ./...

# 开发环境
dev: ## 开发模式运行
	@echo "开发模式运行..."
	@go run main.go

# 文档
docs: ## 生成文档
	@echo "生成文档..."
	@if command -v godoc >/dev/null 2>&1; then \
		godoc -http=:6060 & \
		echo "文档服务器已启动: http://localhost:6060"; \
		echo "按 Ctrl+C 停止服务器"; \
		wait; \
	else \
		echo "godoc 未安装，无法生成文档"; \
	fi

# 性能分析
bench: ## 运行性能测试
	@echo "运行性能测试..."
	@go test -bench=. ./...

profile: ## 生成性能分析文件
	@echo "生成性能分析文件..."
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo "性能分析文件已生成: cpu.prof, mem.prof"

# 容器化
docker-build: ## 构建Docker镜像
	@echo "构建Docker镜像..."
	@docker build -t lxz:latest .

docker-run: ## 运行Docker容器
	@echo "运行Docker容器..."
	@docker run -it --rm lxz:latest

# 安全检查
security: ## 安全漏洞检查
	@echo "检查安全漏洞..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec 未安装，跳过安全检查"; \
	fi

# 代码质量
quality: fmt vet lint security ## 代码质量检查（格式化、静态分析、代码检查、安全检查）

# 完整构建流程
all: clean deps test quality build ## 完整构建流程（清理、依赖、测试、质量检查、构建）

# 发布前检查
pre-release: clean deps test quality cross-build ## 发布前检查（清理、依赖、测试、质量检查、交叉编译）
