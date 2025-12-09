#!/bin/bash

# LXZ 构建脚本
# 自动注入版本信息并构建二进制文件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
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

# 获取Git信息
get_git_info() {
    # 获取最新的tag
    local version=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    
    # 获取commit hash
    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    # 获取commit日期
    local date=$(git log -1 --format=%cd --date=short 2>/dev/null || echo "unknown")
    
    echo "$version|$commit|$date"
}

# 获取Go版本
get_go_version() {
    go version | awk '{print $3}'
}

# 获取平台信息
get_platform_info() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    # 标准化架构名称
    case $arch in
        x86_64)
            arch="amd64"
            ;;
        aarch64)
            arch="arm64"
            ;;
    esac
    
    # Windows环境特殊处理
    if [[ "$os" == *"mingw"* ]] || [[ "$os" == *"msys"* ]]; then
        os="windows"
    fi
    
    echo "$os|$arch"
}

# 构建参数
build_args() {
    local git_info=$(get_git_info)
    local platform_info=$(get_platform_info)
    local go_version=$(get_go_version)
    
    IFS='|' read -r version commit date <<< "$git_info"
    IFS='|' read -r os arch <<< "$platform_info"
    
    # 构建ldflags
    local ldflags="-X 'github.com/liangzhaoliang95/lxz/internal/version.Version=$version'"
    ldflags="$ldflags -X 'github.com/liangzhaoliang95/lxz/internal/version.Commit=$commit'"
    ldflags="$ldflags -X 'github.com/liangzhaoliang95/lxz/internal/version.Date=$date'"
    ldflags="$ldflags -X 'github.com/liangzhaoliang95/lxz/internal/version.GoVersion=$go_version'"
    ldflags="$ldflags -X 'github.com/liangzhaoliang95/lxz/internal/version.Platform=$os'"
    ldflags="$ldflags -X 'github.com/liangzhaoliang95/lxz/internal/version.Arch=$arch'"
    ldflags="$ldflags -X 'github.com/liangzhaoliang95/lxz/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
    
    echo "$ldflags"
}

# 构建函数
build() {
    local target_os=${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}
    local target_arch=${2:-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')}
    local output_name=${3:-"lxz"}
    
    # Windows环境特殊处理
    if [[ "$target_os" == *"mingw"* ]] || [[ "$target_os" == *"msys"* ]]; then
        target_os="windows"
    fi
    
    # 标准化架构名称
    if [[ "$target_arch" == "x86_64" ]]; then
        target_arch="amd64"
    elif [[ "$target_arch" == "aarch64" ]]; then
        target_arch="arm64"
    fi
    
    # 自动为Windows平台添加.exe后缀
    if [[ "$target_os" == "windows" ]] && [[ "$output_name" != *".exe" ]]; then
        output_name="${output_name}.exe"
    fi
    
    print_info "开始构建 $target_os/$target_arch..."
    
    # 设置环境变量
    export GOOS=$target_os
    export GOARCH=$target_arch
    export CGO_ENABLED=0
    
    # 获取构建参数
    local ldflags=$(build_args)
    
    print_info "构建参数: $ldflags"
    
    # 执行构建
    if go build -ldflags="$ldflags" -o "$output_name" ./main.go; then
        print_success "构建成功: $output_name"
        
        # 显示文件信息
        if command -v file >/dev/null 2>&1; then
            print_info "文件信息:"
            file "$output_name"
        fi
        
        if command -v ls >/dev/null 2>&1; then
            print_info "文件大小:"
            ls -lh "$output_name"
        fi
    else
        print_error "构建失败！"
        exit 1
    fi
}

# 交叉编译
cross_build() {
    print_info "开始交叉编译..."
    
    # 创建输出目录
    local output_dir="dist"
    mkdir -p "$output_dir"
    
    # 支持的平台
    local platforms=(
        "linux:amd64:lxz-linux-amd64"
        "linux:arm64:lxz-linux-arm64"
        "darwin:amd64:lxz-darwin-amd64"
        "darwin:arm64:lxz-darwin-arm64"
        "windows:amd64:lxz-windows-amd64.exe"
        "windows:arm64:lxz-windows-arm64.exe"
    )
    
    for platform in "${platforms[@]}"; do
        IFS=':' read -r os arch name <<< "$platform"
        print_info "构建 $os/$arch -> $name"
        build "$os" "$arch" "$output_dir/$name"
    done
    
    print_success "交叉编译完成！输出目录: $output_dir"
}

# 显示帮助信息
show_help() {
    echo "LXZ 构建脚本"
    echo
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  build [os] [arch] [output]  构建指定平台的二进制文件"
    echo "  cross                       交叉编译所有支持的平台"
    echo "  clean                       清理构建产物"
    echo "  help                        显示此帮助信息"
    echo
    echo "示例:"
    echo "  $0 build                    # 构建当前平台"
    echo "  $0 build linux amd64        # 构建Linux AMD64"
    echo "  $0 build windows amd64 lxz.exe  # 构建Windows AMD64"
    echo "  $0 cross                    # 交叉编译所有平台"
    echo
}

# 清理函数
clean() {
    print_info "清理构建产物..."
    
    # 删除构建产物
    rm -f lxz lxz.exe
    rm -rf dist/
    
    print_success "清理完成！"
}

# 主函数
main() {
    case "${1:-build}" in
        "build")
            if [ "$2" = "cross" ]; then
                cross_build
            else
                build "$2" "$3" "$4"
            fi
            ;;
        "cross")
            cross_build
            ;;
        "clean")
            clean
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
