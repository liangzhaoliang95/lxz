#!/bin/bash

# LXZ Release Script
# 自动版本管理和发布脚本

set -e  # 遇到错误立即退出

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

# 检查是否在git仓库中
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "当前目录不是git仓库！"
        exit 1
    fi
}

# 检查是否有未提交的更改
check_uncommitted_changes() {
    if ! git diff-index --quiet HEAD --; then
        print_warning "检测到未提交的更改！"
        echo "请先提交或暂存您的更改，然后重新运行脚本。"
        echo "或者使用 'git stash' 暂存更改。"
        exit 1
    fi
}

# 检查远程仓库配置
check_remote() {
    if ! git remote get-url origin > /dev/null 2>&1; then
        print_error "未找到远程仓库 'origin'！"
        exit 1
    fi
}

# 获取最后一个tag
get_last_tag() {
    local last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -z "$last_tag" ]; then
        print_warning "未找到任何tag，将使用初始版本 v1.0.0"
        echo "v1.0.0"
    else
        echo "$last_tag"
    fi
}

# 解析版本号并递增
increment_version() {
    local version=$1
    local increment_type=${2:-patch}  # 默认递增patch版本
    
    # 移除v前缀
    version=${version#v}
    
    # 分割版本号
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    case $increment_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch|*)
            patch=$((patch + 1))
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# 显示版本信息
show_version_info() {
    local current_version=$1
    local new_version=$2
    
    echo
    print_info "版本信息："
    echo "  当前版本: $current_version"
    echo "  新版本:   $new_version"
    echo
}

# 确认发布
confirm_release() {
    local new_version=$1
    
    echo -e "${YELLOW}即将发布版本: $new_version${NC}"
    echo "此操作将："
    echo "  1. 创建新的tag: $new_version"
    echo "  2. 推送到远程仓库"
    echo "  3. 触发GitHub Actions自动构建"
    echo
    
    read -p "确认发布？(y/N): " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "发布已取消"
        exit 0
    fi
}

# 创建并推送tag
create_and_push_tag() {
    local new_version=$1
    
    print_info "创建tag: $new_version"
    if git tag "$new_version"; then
        print_success "Tag创建成功: $new_version"
    else
        print_error "Tag创建失败！"
        exit 1
    fi
    
    print_info "推送tag到远程仓库..."
    if git push origin "$new_version"; then
        print_success "Tag推送成功！"
    else
        print_error "Tag推送失败！"
        exit 1
    fi
}

# 显示后续步骤
show_next_steps() {
    local new_version=$1
    
    echo
    print_success "版本 $new_version 发布成功！"
    echo
    print_info "后续步骤："
    echo "  1. 查看GitHub Actions构建状态："
    echo "     https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/actions"
    echo "  2. 等待构建完成后查看Release："
    echo "     https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/releases"
    echo "  3. 下载对应平台的二进制文件"
    echo
}

# 主函数
main() {
    echo "🚀 LXZ Release Script"
    echo "======================"
    
    # 检查环境
    check_git_repo
    check_uncommitted_changes
    check_remote
    
    # 获取当前版本
    local current_version=$(get_last_tag)
    print_info "当前版本: $current_version"
    
    # 确定版本递增类型
    echo
    echo "请选择版本递增类型："
    echo "  1) patch (补丁版本，默认)"
    echo "  2) minor (次要版本)"
    echo "  3) major (主要版本)"
    echo "  4) 自定义版本号"
    echo
    
    read -p "请选择 (1-4，默认1): " -n 1 -r
    echo
    
    local increment_type="patch"
    local new_version=""
    
    case $REPLY in
        1|"")
            increment_type="patch"
            new_version=$(increment_version "$current_version" "patch")
            ;;
        2)
            increment_type="minor"
            new_version=$(increment_version "$current_version" "minor")
            ;;
        3)
            increment_type="major"
            new_version=$(increment_version "$current_version" "major")
            ;;
        4)
            echo
            read -p "请输入自定义版本号 (格式: v1.2.3): " custom_version
            if [[ $custom_version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                new_version=$custom_version
            else
                print_error "版本号格式错误！请使用格式: v1.2.3"
                exit 1
            fi
            ;;
        *)
            print_error "无效选择！"
            exit 1
            ;;
    esac
    
    # 显示版本信息
    show_version_info "$current_version" "$new_version"
    
    # 确认发布
    confirm_release "$new_version"
    
    # 拉取最新代码
    print_info "拉取最新代码..."
    if ! git pull origin main; then
        print_warning "拉取代码失败，尝试拉取当前分支..."
        current_branch=$(git branch --show-current)
        if ! git pull origin "$current_branch"; then
            print_error "拉取代码失败！"
            exit 1
        fi
    fi
    
    # 创建并推送tag
    create_and_push_tag "$new_version"
    
    # 显示后续步骤
    show_next_steps "$new_version"
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
