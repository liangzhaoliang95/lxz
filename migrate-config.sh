#!/bin/bash

# LXZ Configuration Migration Script
# Migrates configuration from old XDG locations to ~/.lxz

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Determine OS and old config location
detect_old_config_dir() {
    local os=$(uname -s)
    case "$os" in
        Darwin)
            echo "$HOME/Library/Application Support/lxz"
            ;;
        Linux)
            echo "$HOME/.config/lxz"
            ;;
        *)
            echo ""
            ;;
    esac
}

# New config location
NEW_CONFIG_DIR="$HOME/.lxz"
OLD_CONFIG_DIR=$(detect_old_config_dir)

print_info "LXZ Configuration Migration"
echo ""
print_info "Old location: $OLD_CONFIG_DIR"
print_info "New location: $NEW_CONFIG_DIR"
echo ""

# Check if old config exists
if [ -z "$OLD_CONFIG_DIR" ]; then
    print_error "Unsupported operating system"
    exit 1
fi

if [ ! -d "$OLD_CONFIG_DIR" ]; then
    print_warning "Old configuration directory not found: $OLD_CONFIG_DIR"
    print_info "Nothing to migrate. The new configuration will be created at: $NEW_CONFIG_DIR"
    exit 0
fi

# Check if new config already exists
if [ -d "$NEW_CONFIG_DIR" ]; then
    print_warning "New configuration directory already exists: $NEW_CONFIG_DIR"
    read -p "Do you want to merge/overwrite? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Migration cancelled"
        exit 0
    fi
else
    # Create new config directory
    print_info "Creating new configuration directory..."
    mkdir -p "$NEW_CONFIG_DIR"
fi

# Copy configuration files
print_info "Copying configuration files..."

# List of files to migrate
files=(
    "config.yaml"
    "app_database_config.yaml"
    "app_redis_config.yaml"
    "hotkeys.yaml"
    "aliases.yaml"
    "plugins.yaml"
    "views.yaml"
)

for file in "${files[@]}"; do
    if [ -f "$OLD_CONFIG_DIR/$file" ]; then
        cp "$OLD_CONFIG_DIR/$file" "$NEW_CONFIG_DIR/"
        print_success "Copied: $file"
    fi
done

# Copy directories
dirs=(
    "skins"
    "screen-dumps"
)

for dir in "${dirs[@]}"; do
    if [ -d "$OLD_CONFIG_DIR/$dir" ]; then
        cp -r "$OLD_CONFIG_DIR/$dir" "$NEW_CONFIG_DIR/"
        print_success "Copied: $dir/"
    fi
done

# Copy log file if exists
if [ -f "$OLD_CONFIG_DIR/lxz.log" ]; then
    cp "$OLD_CONFIG_DIR/lxz.log" "$NEW_CONFIG_DIR/"
    print_success "Copied: lxz.log"
fi

echo ""
print_success "Migration completed successfully!"
echo ""
print_info "You can now safely remove the old configuration directory:"
print_info "  rm -rf '$OLD_CONFIG_DIR'"
echo ""
print_info "Or keep it as a backup."
