# LXZ Makefile
# 简化构建和版本管理流程

.PHONY: help build clean test release version check-update install-golangci-lint install-goimports install-golines lint-staged format-staged

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
fmt: install-goimports install-golines ## 格式化代码（使用 golines 和 goimports）
	@echo "✨ 格式化代码..."
	@if [ -n "$(FILES)" ]; then \
		echo "格式化指定文件..."; \
		echo "1. 使用 goimports 优化 import 顺序..."; \
		for file in $(FILES); do \
			goimports -w "$$file" || true; \
		done; \
		echo "2. 使用 golines 格式化代码..."; \
		for file in $(FILES); do \
			golines -w --max-len=100 "$$file" || true; \
		done; \
	else \
		echo "格式化整个项目..."; \
		echo "1. 使用 goimports 优化 import 顺序..."; \
		find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -not -path "./bin/*" -not -path "./dist/*" -type f -exec goimports -w {} \; 2>/dev/null || true; \
		echo "2. 使用 golines 格式化代码..."; \
		golines -w --max-len=100 --ignored-dirs=vendor,.git,bin,dist . || true; \
	fi
	@echo "✅ 代码格式化完成!"

vet: ## 代码静态分析
	@echo "代码静态分析..."
	@go vet ./...

# 工具安装
install-golangci-lint: ## 安装 golangci-lint
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "📥 安装 golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
		echo "✅ golangci-lint 安装完成"; \
	else \
		echo "✅ golangci-lint 已安装"; \
	fi

install-goimports: ## 安装 goimports
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "📥 安装 goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest || exit 1; \
		echo "✅ goimports 安装完成"; \
	else \
		echo "✅ goimports 已安装"; \
	fi
	@command -v goimports >/dev/null 2>&1 || (echo "❌ goimports 安装失败或未找到" && exit 1)

install-golines: ## 安装 golines
	@if ! command -v golines >/dev/null 2>&1; then \
		echo "📥 安装 golines..."; \
		go install github.com/segmentio/golines@latest || exit 1; \
		echo "✅ golines 安装完成"; \
	else \
		echo "✅ golines 已安装"; \
	fi
	@command -v golines >/dev/null 2>&1 || (echo "❌ golines 安装失败或未找到" && exit 1)

lint: install-golangci-lint ## 代码检查
	@echo "🔍 运行代码检查..."
	@if [ -n "$(PACKAGES)" ]; then \
		for pkg in $(PACKAGES); do \
			echo "检查包: $$pkg"; \
			golangci-lint run "$$pkg" || exit 1; \
		done; \
	else \
		golangci-lint run ./...; \
	fi
	@echo "✅ 代码检查完成!"

lint-staged: install-golangci-lint ## 检查暂存的 Go 文件
	@echo "🔍 检查暂存的 Go 文件..."
	@staged_files=$$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$$' || true); \
	if [ -z "$$staged_files" ]; then \
		echo "ℹ️  没有暂存的 Go 文件"; \
		exit 0; \
	fi; \
	packages=$$(echo "$$staged_files" | xargs -n1 dirname | sort -u | grep -v '^\.$$' | sed 's|^|./|' || true); \
	if [ -z "$$packages" ]; then \
		echo "ℹ️  没有找到对应的包"; \
		exit 0; \
	fi; \
	echo "发现包: $$(echo "$$packages" | wc -l | tr -d ' ') 个"; \
	echo "$$packages" | while IFS= read -r pkg; do \
		if [ -n "$$pkg" ] && [ -d "$$pkg" ]; then \
			echo "  检查包: $$pkg"; \
			golangci-lint run "$$pkg" || exit 1; \
		fi; \
	done
	@echo "✅ 检查完成!"

format-staged: install-goimports install-golines ## 格式化暂存的 Go 文件
	@echo "✨ 格式化暂存的 Go 文件..."
	@staged_files=$$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$$' || true); \
	if [ -z "$$staged_files" ]; then \
		echo "ℹ️  没有暂存的 Go 文件"; \
		exit 0; \
	fi; \
	echo "发现暂存文件: $$(echo "$$staged_files" | wc -l | tr -d ' ') 个"; \
	echo "$$staged_files" | while IFS= read -r file; do \
		if [ -n "$$file" ] && [ -f "$$file" ]; then \
			echo "  格式化: $$file"; \
			goimports -w "$$file" || true; \
			golines -w --max-len=100 "$$file" || true; \
		fi; \
	done
	@echo "✅ 格式化完成!"

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
