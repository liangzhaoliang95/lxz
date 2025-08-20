<div align="center">
  <img src="https://img.shields.io/badge/Go-1.24.3+-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg" alt="Platform">
</div>

<div align="center">
  <h1>LXZ - DevOps 图形化 CLI 工具</h1>
  <p><strong>一个强大的 DevOps 图形化命令行界面工具，支持数据库、Docker、Redis、文件系统、Kubernetes 等多种资源管理</strong></p>
  
  <div>
    <a href="README.md" style="background: #f0f0f0; color: #333; padding: 8px 16px; margin: 5px; border-radius: 4px; text-decoration: none; display: inline-block;">English</a>
    <a href="README_CN.md" style="background: #007acc; color: white; padding: 8px 16px; margin: 5px; border-radius: 4px; text-decoration: none; display: inline-block;">中文</a>
  </div>
</div>

---

## 🌟 功能特性

### 🚀 核心功能
- **数据库管理**: 连接和管理 MySQL 数据库，提供直观的用户界面
- **Docker 操作**: 容器管理、日志查看、Shell 访问等
- **Redis 管理**: Redis 连接管理和数据操作
- **文件浏览器**: 高级文件系统导航，支持文件预览
- **Kubernetes 集成**: K9s 配置管理和集群访问
- **SSH 连接管理**: 集中式 SSH 主机管理
- **跨平台支持**: 支持 Linux、macOS 和 Windows (AMD64 & ARM64)

### 🎯 主要优势
- **终端界面**: 基于 tview 构建的丰富 TUI，提供出色的终端体验
- **热键支持**: 全面的键盘快捷键，提高导航效率
- **插件系统**: 可扩展架构，支持自定义功能
- **配置管理**: 基于 YAML 的配置，支持验证
- **日志监控**: 内置日志系统，支持可配置级别

## 📦 安装

### 前置要求
- Go 1.24.3 或更高版本
- Git

### 快速安装
```bash
# 克隆仓库
git clone https://github.com/liangzhaoliang95/lxz.git
cd lxz

# 构建并安装
go build -o lxz ./main.go
sudo mv lxz /usr/local/bin/
sudo chmod +x /usr/local/bin/lxz
```

### 下载预构建二进制文件
访问 [Releases](https://github.com/liangzhaoliang95/lxz/releases) 页面下载适合您平台的预构建二进制文件。

## 🚀 使用方法

### 基本命令
```bash
# 启动 LXZ
lxz

# 自定义刷新率启动
lxz --refresh 5

# 调试日志级别启动
lxz --logLevel debug

# 无头模式启动
lxz --headless
```

### 配置
LXZ 使用 YAML 配置文件，位置在：
- Linux/macOS: `~/.config/lxz/`
- Windows: `%APPDATA%\lxz\`

### 快捷键
- `F` - 切换全屏模式
- `Ctrl+R` - 刷新数据
- `Ctrl+N` - 创建新项目
- `Ctrl+D` - 删除项目
- `Enter` - 选择/执行
- `Tab` - 切换焦点
- `Escape` - 退出全屏/返回

## 🏗️ 架构

### 项目结构
```
lxz/
├── cmd/           # 命令行接口
├── internal/      # 内部包
│   ├── config/    # 配置管理
│   ├── drivers/   # 外部服务驱动
│   ├── model/     # 数据模型
│   ├── ui/        # 用户界面组件
│   ├── view/      # 视图层组件
│   └── helper/    # 工具函数
├── main.go        # 应用程序入口点
└── go.mod         # Go 模块定义
```

### 核心组件
- **视图层**: 处理不同资源视图（数据库、Docker、Redis 等）
- **UI 层**: 管理终端 UI 组件和交互
- **驱动层**: 抽象外部服务连接
- **配置层**: 管理应用程序配置和验证

## 🔧 开发

### 从源码构建
```bash
# 克隆和设置
git clone https://github.com/liangzhaoliang95/lxz.git
cd lxz

# 安装依赖
go mod download

# 构建
go build -o lxz ./main.go

# 运行测试
go test ./...
```

### 贡献
1. Fork 仓库
2. 创建功能分支
3. 进行更改
4. 添加测试（如果适用）
5. 提交 Pull Request

### 开发依赖
- Go 1.24.3+
- tview (终端 UI 框架)
- tcell (终端单元格库)
- cobra (CLI 框架)

## 📚 文档

- [发布指南](RELEASE_GUIDE.md) - 如何创建发布版本
- [配置指南](docs/configuration.md) - 配置选项
- [API 参考](docs/api.md) - API 文档
- [贡献指南](CONTRIBUTING.md) - 如何贡献

## 🤝 贡献

我们欢迎贡献！请查看我们的 [贡献指南](CONTRIBUTING.md) 了解详情。

### 贡献方式
- 🐛 报告错误
- 💡 建议新功能
- 📝 改进文档
- 🔧 提交 Pull Request
- ⭐ 给仓库加星标

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [tview](https://github.com/rivo/tview) - 终端 UI 框架
- [tcell](https://github.com/gdamore/tcell) - 终端单元格库
- [cobra](https://github.com/spf13/cobra) - CLI 框架
- [Go 社区](https://golang.org/) - 感谢这个优秀的语言

---

<div align="center">
  <p>用 ❤️ 制作，来自 LXZ 团队</p>
</div>
