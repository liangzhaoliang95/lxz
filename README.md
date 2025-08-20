<div align="center">
  <img src="https://img.shields.io/badge/Go-1.24.3+-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg" alt="Platform">
</div>

<div align="center">
  <h1>LXZ - DevOps 图形化 CLI 工具</h1>
  <p><strong>一个强大的 DevOps 图形化命令行界面工具，支持数据库、Docker、Redis、文件系统、Kubernetes 等多种资源管理</strong></p>
  
  <div>
    <button onclick="switchLanguage('en')" id="en-btn" style="background: #007acc; color: white; border: none; padding: 8px 16px; margin: 5px; border-radius: 4px; cursor: pointer;">English</button>
    <button onclick="switchLanguage('zh')" id="zh-btn" style="background: #f0f0f0; color: #333; border: none; padding: 8px 16px; margin: 5px; border-radius: 4px; cursor: pointer;">中文</button>
  </div>
</div>

---

## 🌟 Features / 功能特性

<div id="en-content">
### 🚀 Core Features
- **Database Management**: Connect and manage MySQL databases with intuitive UI
- **Docker Operations**: Container management, logs viewing, shell access, and more
- **Redis Management**: Redis connection management and data operations
- **File Browser**: Advanced file system navigation with preview capabilities
- **Kubernetes Integration**: K9s configuration management and cluster access
- **SSH Connection Manager**: Centralized SSH host management
- **Cross-Platform**: Support for Linux, macOS, and Windows (AMD64 & ARM64)

### 🎯 Key Benefits
- **Terminal-Based UI**: Rich TUI built with tview for excellent terminal experience
- **Hotkey Support**: Comprehensive keyboard shortcuts for efficient navigation
- **Plugin System**: Extensible architecture for custom functionality
- **Configuration Management**: YAML-based configuration with validation
- **Logging & Monitoring**: Built-in logging system with configurable levels
</div>

<div id="zh-content" style="display: none;">
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
</div>

## 📦 Installation / 安装

<div id="en-content">
### Prerequisites
- Go 1.24.3 or higher
- Git

### Quick Install
```bash
# Clone the repository
git clone https://github.com/liangzhaoliang95/lxz.git
cd lxz

# Build and install
go build -o lxz ./main.go
sudo mv lxz /usr/local/bin/
sudo chmod +x /usr/local/bin/lxz
```

### Download Pre-built Binaries
Visit [Releases](https://github.com/liangzhaoliang95/lxz/releases) page to download pre-built binaries for your platform.
</div>

<div id="zh-content" style="display: none;">
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
</div>

## 🚀 Usage / 使用方法

<div id="en-content">
### Basic Commands
```bash
# Start LXZ
lxz

# Start with custom refresh rate
lxz --refresh 5

# Start with debug logging
lxz --logLevel debug

# Start in headless mode
lxz --headless
```

### Configuration
LXZ uses YAML configuration files located in:
- Linux/macOS: `~/.config/lxz/`
- Windows: `%APPDATA%\lxz\`

### Key Bindings
- `F` - Toggle fullscreen mode
- `Ctrl+R` - Refresh data
- `Ctrl+N` - Create new item
- `Ctrl+D` - Delete item
- `Enter` - Select/Execute
- `Tab` - Switch focus
- `Escape` - Exit fullscreen/Go back
</div>

<div id="zh-content" style="display: none;">
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
</div>

## 🏗️ Architecture / 架构

<div id="en-content">
### Project Structure
```
lxz/
├── cmd/           # Command line interface
├── internal/      # Internal packages
│   ├── config/    # Configuration management
│   ├── drivers/   # External service drivers
│   ├── model/     # Data models
│   ├── ui/        # User interface components
│   ├── view/      # View layer components
│   └── helper/    # Utility functions
├── main.go        # Application entry point
└── go.mod         # Go module definition
```

### Core Components
- **View Layer**: Handles different resource views (Database, Docker, Redis, etc.)
- **UI Layer**: Manages terminal UI components and interactions
- **Driver Layer**: Abstracts external service connections
- **Config Layer**: Manages application configuration and validation
</div>

<div id="zh-content" style="display: none;">
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
</div>

## 🔧 Development / 开发

<div id="en-content">
### Building from Source
```bash
# Clone and setup
git clone https://github.com/liangzhaoliang95/lxz.git
cd lxz

# Install dependencies
go mod download

# Build
go build -o lxz ./main.go

# Run tests
go test ./...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### Development Dependencies
- Go 1.24.3+
- tview (terminal UI framework)
- tcell (terminal cell library)
- cobra (CLI framework)
</div>

<div id="zh-content" style="display: none;">
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
</div>

## 📚 Documentation / 文档

<div id="en-content">
- [Release Guide](RELEASE_GUIDE.md) - How to create releases
- [Configuration Guide](docs/configuration.md) - Configuration options
- [API Reference](docs/api.md) - API documentation
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
</div>

<div id="zh-content" style="display: none;">
- [发布指南](RELEASE_GUIDE.md) - 如何创建发布版本
- [配置指南](docs/configuration.md) - 配置选项
- [API 参考](docs/api.md) - API 文档
- [贡献指南](CONTRIBUTING.md) - 如何贡献
</div>

## 🤝 Contributing / 贡献

<div id="en-content">
We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Ways to Contribute
- 🐛 Report bugs
- 💡 Suggest new features
- 📝 Improve documentation
- 🔧 Submit pull requests
- ⭐ Star the repository
</div>

<div id="zh-content" style="display: none;">
我们欢迎贡献！请查看我们的 [贡献指南](CONTRIBUTING.md) 了解详情。

### 贡献方式
- 🐛 报告错误
- 💡 建议新功能
- 📝 改进文档
- 🔧 提交 Pull Request
- ⭐ 给仓库加星标
</div>

## 📄 License / 许可证

<div id="en-content">
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
</div>

<div id="zh-content" style="display: none;">
本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。
</div>

## 🙏 Acknowledgments / 致谢

<div id="en-content">
- [tview](https://github.com/rivo/tview) - Terminal UI framework
- [tcell](https://github.com/gdamore/tcell) - Terminal cell library
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [Go community](https://golang.org/) - For the amazing language
</div>

<div id="zh-content" style="display: none;">
- [tview](https://github.com/rivo/tview) - 终端 UI 框架
- [tcell](https://github.com/gdamore/tcell) - 终端单元格库
- [cobra](https://github.com/spf13/cobra) - CLI 框架
- [Go 社区](https://golang.org/) - 感谢这个优秀的语言
</div>

---

<div align="center">
  <p>Made with ❤️ by the LXZ team</p>
  <p>用 ❤️ 制作，来自 LXZ 团队</p>
</div>

<script>
function switchLanguage(lang) {
  const enContent = document.getElementById('en-content');
  const zhContent = document.getElementById('zh-content');
  const enBtn = document.getElementById('en-btn');
  const zhBtn = document.getElementById('zh-btn');
  
  if (lang === 'en') {
    enContent.style.display = 'block';
    zhContent.style.display = 'none';
    enBtn.style.background = '#007acc';
    enBtn.style.color = 'white';
    zhBtn.style.background = '#f0f0f0';
    zhBtn.style.color = '#333';
  } else {
    enContent.style.display = 'none';
    zhContent.style.display = 'block';
    zhBtn.style.background = '#007acc';
    zhBtn.style.color = 'white';
    enBtn.style.background = '#f0f0f0';
    enBtn.style.color = '#333';
  }
}

// Initialize with English
document.addEventListener('DOMContentLoaded', function() {
  switchLanguage('en');
});
</script>
