<div align="center">
  <img src="https://img.shields.io/badge/Go-1.24.3+-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg" alt="Platform">
</div>

<div align="center">
  <h1>LXZ - DevOps Graphical CLI Tool</h1>
  <p><strong>A powerful DevOps graphical command-line interface tool supporting database, Docker, Redis, file system, Kubernetes and other resource management</strong></p>
  
  <div>
    <a href="README.md" style="background: #007acc; color: white; padding: 8px 16px; margin: 5px; border-radius: 4px; text-decoration: none; display: inline-block;">English</a>
    <a href="README_CN.md" style="background: #f0f0f0; color: #333; padding: 8px 16px; margin: 5px; border-radius: 4px; text-decoration: none; display: inline-block;">中文</a>
  </div>
</div>

---

## 🌟 Features

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

## 📦 Installation

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

## 🚀 Usage

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

## 🏗️ Architecture

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

## 🔧 Development

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

## 📚 Documentation

- [Release Guide](RELEASE_GUIDE.md) - How to create releases
- [Configuration Guide](docs/configuration.md) - Configuration options
- [API Reference](docs/api.md) - API documentation
- [Contributing Guide](CONTRIBUTING.md) - How to contribute

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Ways to Contribute
- 🐛 Report bugs
- 💡 Suggest new features
- 📝 Improve documentation
- 🔧 Submit pull requests
- ⭐ Star the repository

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [tview](https://github.com/rivo/tview) - Terminal UI framework
- [tcell](https://github.com/gdamore/tcell) - Terminal cell library
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [Go community](https://golang.org/) - For the amazing language

---

<div align="center">
  <p>Made with ❤️ by the LXZ team</p>
</div>
