# LXZ 版本管理机制

本文档描述了 LXZ 项目的版本管理机制，包括版本信息注入、GitHub 版本检查和自动更新功能。

## 功能特性

### 1. 版本信息注入
- 在编译时自动注入 Git 标签、提交哈希、构建时间等信息
- 支持跨平台构建，自动识别目标平台和架构
- 版本信息通过 `ldflags` 注入到二进制文件中

### 2. GitHub 版本检查
- 自动从 GitHub API 获取最新的 release 版本
- 比较当前版本与最新版本，提示用户更新
- 支持自定义仓库地址和检查 URL

### 3. 自动更新提示
- 在应用启动时自动检查更新（可配置）
- 提供详细的更新信息和下载链接
- 支持手动检查更新

## 使用方法

### 构建项目

#### 使用 Makefile（推荐）
```bash
# 构建当前平台
make build

# 构建特定平台
make build-linux
make build-windows
make build-darwin

# 交叉编译所有平台
make cross-build

# 显示帮助信息
make help
```

#### 使用构建脚本
```bash
# 构建当前平台
./scripts/build.sh build

# 构建特定平台
./scripts/build.sh build linux amd64 lxz-linux-amd64

# 交叉编译
./scripts/build.sh cross

# 清理构建产物
./scripts/build.sh clean
```

### 版本管理

#### 查看版本信息
```bash
# 显示完整版本信息
lxz version

# 显示简短版本信息
lxz version --short

# 检查更新
lxz version --check-update
```

#### 发布新版本
```bash
# 使用发布脚本
./release.sh

# 或使用 Makefile
make release
```

## 配置说明

### 环境变量

可以通过环境变量自定义版本检查行为：

```bash
# 设置仓库地址
export LXZ_REPOSITORY="your-username/your-repo"

# 禁用自动检查
export LXZ_AUTO_CHECK="false"

# 自定义检查 URL（用于私有部署）
export LXZ_CHECK_URL="https://your-api.example.com"
```

### 配置文件

版本配置支持从配置文件读取，配置文件格式支持 JSON 和 YAML：

```yaml
# config.yaml
version:
  repository: "liangzhaoliang95/lxz"
  auto_check: true
  check_url: "https://api.github.com"
```

## 版本号规范

项目使用语义化版本号（Semantic Versioning）：

- **主版本号（Major）**：不兼容的 API 修改
- **次版本号（Minor）**：向下兼容的功能性新增
- **修订号（Patch）**：向下兼容的问题修正

版本号格式：`vX.Y.Z`

示例：
- `v1.0.0` - 第一个正式版本
- `v1.2.3` - 主版本1，次版本2，修订号3
- `v2.0.0` - 主版本2，可能包含破坏性变更

## 发布流程

### 1. 准备发布
```bash
# 确保代码已提交
git status

# 拉取最新代码
git pull origin main
```

### 2. 运行发布脚本
```bash
./release.sh
```

发布脚本会：
- 检查 Git 仓库状态
- 提示选择版本递增类型
- 创建新的 Git 标签
- 推送到远程仓库
- 触发 GitHub Actions 自动构建

### 3. 自动构建
GitHub Actions 会自动：
- 构建多平台二进制文件
- 创建 GitHub Release
- 上传构建产物
- 生成安装说明

## 开发模式

在开发模式下，版本信息会显示为：
- Version: `dev`
- Commit: `dev`
- Date: `unknown`

开发版本无法进行版本检查和更新检查。

## 故障排除

### 常见问题

1. **版本检查失败**
   - 检查网络连接
   - 确认仓库地址正确
   - 检查 GitHub API 限制

2. **构建失败**
   - 确认 Go 版本兼容性
   - 检查依赖是否正确安装
   - 查看构建日志

3. **版本信息不显示**
   - 确认使用正确的构建脚本
   - 检查 ldflags 参数
   - 验证版本包导入

### 调试模式

启用调试信息：
```bash
export LXZ_DEBUG=true
go run main.go version --check-update
```

## 贡献指南

### 添加新功能
1. 在 `internal/version/` 包中添加新功能
2. 更新相关测试
3. 更新文档

### 修改构建流程
1. 更新 `scripts/build.sh`
2. 更新 `Makefile`
3. 更新 GitHub Actions 配置

### 报告问题
1. 检查现有 issue
2. 创建新的 issue
3. 提供详细的错误信息和复现步骤

## 相关文件

- `internal/version/` - 版本管理核心包
- `scripts/build.sh` - 构建脚本
- `Makefile` - 构建工具
- `release.sh` - 发布脚本
- `.github/workflows/release.yml` - GitHub Actions 配置

## 更新日志

### v1.0.0
- 初始版本管理功能
- 支持版本信息注入
- 支持 GitHub 版本检查
- 支持自动更新提示

### v1.1.0
- 添加配置管理
- 支持环境变量配置
- 改进错误处理
- 添加调试模式
