# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LXZ is a Go-based Terminal User Interface (TUI) application for DevOps resource management. It provides graphical interfaces for managing databases (MySQL), Redis, Docker containers, file systems, Kubernetes (K9s), and SSH connections—all within the terminal.

Built with the tview library (forked version), it follows a clean layered architecture with separation between UI components, business logic, drivers, and configuration.

## Build and Development Commands

### Building

```bash
# Build for current platform
make build
# OR
./scripts/build.sh build

# Build for specific platforms
make build-linux      # Linux AMD64
make build-windows    # Windows AMD64
make build-darwin     # macOS AMD64

# Cross-compile all platforms (6 targets: Linux/macOS/Windows on AMD64/ARM64)
make cross-build
# OR
./scripts/build.sh cross
```

The build script (`scripts/build.sh`) automatically injects version information via ldflags:
- Git tag version
- Commit hash
- Build date/time
- Go version
- Platform and architecture

Output binaries go to `dist/` directory for cross-builds.

### Testing

```bash
# Run all tests
make test
# OR
go test ./...

# Run tests with coverage
make test-coverage   # Generates coverage.html

# Run benchmarks
make bench
```

Note: The codebase has minimal test coverage currently (only 2 test files found).

### Development

```bash
# Run in development mode
make dev
# OR
go run main.go

# Run with specific flags
go run main.go --logLevel debug
go run main.go --refresh 5      # 5 second refresh rate
go run main.go --headless       # No header UI
go run main.go --splashless     # Skip splash screen

# Format code
make fmt

# Static analysis
make vet

# Lint (requires golangci-lint)
make lint

# All quality checks (fmt + vet + lint + security)
make quality
```

### Version Management

```bash
# Show version info
make version
# OR
go run main.go version

# Check for updates
make check-update
# OR
go run main.go version --check-update
```

### Releasing

```bash
# Create new release (interactive script)
make release
# OR
./release.sh
```

The release process:
1. Checks git status
2. Prompts for version bump type (major/minor/patch)
3. Creates and pushes git tag (e.g., `v1.0.0`)
4. GitHub Actions automatically builds 6 platform binaries and creates GitHub Release

See RELEASE_GUIDE.md and VERSION_MANAGEMENT.md for details.

### Other Commands

```bash
# Clean build artifacts
make clean

# Download dependencies
make deps

# Update dependencies
make deps-update

# Install to system
make install
```

## Architecture Overview

### Layer Structure

```
cmd/                    Entry point (Cobra CLI)
    ↓
internal/view/          Resource browsers (business logic)
    ↓
internal/ui/            tview component wrappers
    ↓
internal/drivers/       External service adapters
internal/config/        Configuration management
```

### Key Architectural Patterns

**1. Component Lifecycle**
All views implement the `Component` interface (internal/model/types.go):
- `Init(ctx)` - Initialization with dependency injection
- `Start()` - Called when component becomes active
- `Stop()` - Called when component is removed
- `Hints()` - Returns keyboard shortcuts for menu

Lifecycle managed by `PageStack` which pushes/pops components onto a navigation stack.

**2. View Layer (internal/view/)**

The `view.App` struct is the main controller:
- Routes F1-F6 keys to different resource browsers
- Manages the `PageStack` for navigation
- Injects dependencies into components via `inject()`

All browsers inherit from `BaseFlex`:
- `DatabaseBrowser` - MySQL table browsing, query execution
- `RedisBrowser` - Redis key-value operations
- `DockerBrowser` - Container management, logs, shell access
- `FileBrowser` - Directory tree navigation with preview
- `K9SBrowser` - Kubernetes monitoring
- `SshConnect` - SSH connection manager

Each browser can have sub-views that get pushed onto the stack (e.g., DatabaseDbListView → DatabaseTableView → DatabaseQueryView).

**3. UI Layer (internal/ui/)**

`ui.App` wraps `tview.Application` and manages:
- Global components: Logo, Menu (breadcrumbs), SubMenu (keyboard hints), Status, Flash (notifications)
- Keyboard routing via `KeyActions` map
- Thread-safe UI updates via `QueueUpdateDraw()`

Modal dialogs in `ui/dialog/`:
- Connection forms (database, Redis)
- File operations (create, delete, rename)
- Confirmation dialogs
- Loading screens

**4. Driver Layer (internal/drivers/)**

Abstract interfaces for external services:
- `IDatabaseConn` - Database operations (currently MySQL via GORM)
- `RedisClient` - Redis operations (go-redis wrapper)
- Docker driver (moby/moby SDK)

Drivers use:
- Factory pattern for initialization
- Connection pooling via `sync.Map`
- Lazy initialization with `GetConnectOrInit()`

**5. Configuration (internal/config/)**

YAML-based config stored in `~/.config/lxz/`:
- `Config` → `LXZ` → resource-specific configs (DatabaseConfig, RedisConfig)
- JSON schema validation (config/json/schemas/)
- Load/Save/Merge methods
- `Styles` system for theming with listener pattern

### Data Flow Example

User presses F3 (Redis Browser):
1. `view.App.menuPageChange()` captures F3
2. Creates `NewRedisBrowser(app)`
3. Calls `app.inject(browser, true)` → `browser.Init(ctx)` → loads RedisConfig
4. `PageStack.Push()` triggers `StackPushed()` → calls `browser.Start()`
5. UI renders connection table
6. User interaction → `browser.Keyboard()` handles keys
7. `RedisDriver` methods fetch/update data
8. Results displayed in table

### Keyboard Event Routing

1. `tview.Application.SetInputCapture()` captures all key events
2. Routed to `view.App.keyboard()`
3. Checks `ui.App.HasAction(key)` for global actions
4. Otherwise passed to focused component's `Keyboard()` method
5. Component checks its `KeyActions` map

### Observer Pattern Usage

- `StackListener` - Components notified on push/pop from PageStack
- `StyleListener` - Components receive theme change notifications
- Menu listeners - Update breadcrumb navigation

## Extension Points

### Adding a New Resource Browser

1. Create struct embedding `*BaseFlex` in `internal/view/`
2. Implement `Component` interface methods
3. Register in `view.App.menuPageChange()` with function key binding
4. Add keyboard actions via `bindKeys()`

### Adding a New Database Driver

1. Implement `IDatabaseConn` interface in `internal/drivers/database_drivers/`
2. Update `_initDriver()` factory to recognize new provider
3. Add configuration struct in `internal/config/database_config.go`

### Adding New Dialogs

1. Create form in `internal/ui/dialog/`
2. Optionally extend `BaseModelForm` for validation
3. Push onto `ui.Pages` as modal
4. Use callbacks for form submission

## Important Notes

- **tview Fork**: Uses custom fork `github.com/liangzhaoliang95/tview` (not upstream rivo/tview)
- **Thread Safety**: Always use `app.QueueUpdateDraw()` for UI updates from goroutines
- **CGO**: Builds with `CGO_ENABLED=0` for static binaries
- **Go Version**: Requires Go 1.24.3+
- **Config Location**: `~/.lxz/` (all platforms) - stores all configuration files, logs, and data

### Configuration Directory Structure

All configuration files are stored in `~/.lxz/`:

```
~/.lxz/
├── config.yaml                  # Main configuration
├── app_database_config.yaml     # Database connections
├── app_redis_config.yaml        # Redis connections
├── hotkeys.yaml                 # Keyboard shortcuts
├── aliases.yaml                 # Command aliases
├── plugins.yaml                 # Plugin configuration
├── views.yaml                   # Custom views
├── lxz.log                      # Application logs
├── skins/                       # Color themes
└── screen-dumps/                # Screenshots
```

### Migrating from Old Config Location

If you have existing configurations in the old XDG locations:
- **macOS**: `~/Library/Application Support/lxz/`
- **Linux**: `~/.config/lxz/`

Simply copy all files to `~/.lxz/`:
```bash
# macOS
cp -r ~/Library/Application\ Support/lxz/* ~/.lxz/

# Linux
cp -r ~/.config/lxz/* ~/.lxz/
```

### Custom Config Directory

You can override the default location using the `LXZ_CONFIG_DIR` environment variable:
```bash
export LXZ_CONFIG_DIR=/path/to/custom/config
lxz
```

## Key Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/liangzhaoliang95/tview` - TUI framework (forked)
- `github.com/gdamore/tcell/v2` - Terminal cell library
- `gorm.io/gorm` - ORM for database operations
- `github.com/go-redis/redis/v8` - Redis client
- `github.com/moby/moby` - Docker SDK
- `gopkg.in/yaml.v3` - YAML configuration
- `github.com/xeipuuv/gojsonschema` - JSON schema validation

## Testing Strategy

Currently minimal test coverage. When adding tests:
- Place `*_test.go` files next to source files
- Use `github.com/stretchr/testify` for assertions
- Run via `go test ./...` or `make test`

## Code Conventions

- Package-level interfaces in `internal/model/types.go`
- Structured logging via `internal/slogs` with slog attributes
- Helper utilities in `internal/helper/`
- Color utilities in `internal/color/`
- Version info injected at build time in `internal/version/`