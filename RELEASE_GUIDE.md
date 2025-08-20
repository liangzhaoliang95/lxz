# 自动发布指南

本项目使用GitHub Actions自动构建和发布release产物。当您创建并推送一个tag时，系统会自动构建6个不同平台的二进制文件。

## 支持的平台

- **Linux**: AMD64, ARM64
- **macOS**: AMD64, ARM64  
- **Windows**: AMD64, ARM64

## 如何发布新版本

### 1. 创建并推送tag

```bash
# 确保代码已提交并推送到远程仓库
git add .
git commit -m "准备发布 v1.0.0"
git push origin main

# 创建tag
git tag v1.0.0

# 推送tag到远程仓库（这会触发GitHub Action）
git push origin v1.0.0
```

### 2. 监控构建过程

1. 前往GitHub仓库页面
2. 点击"Actions"标签页
3. 查看"Release Build"工作流的执行状态
4. 等待所有6个平台的构建完成

### 3. 自动创建Release

构建完成后，系统会自动：
- 在GitHub Releases页面创建新的release
- 上传所有平台的二进制文件
- 生成安装说明和下载链接

## 构建产物

每个平台都会生成以下文件：

### Linux/macOS
- 可执行文件（如：`lxz-linux-amd64`）
- 压缩包（如：`lxz-linux-amd64.tar.gz`）

### Windows
- 可执行文件（如：`lxz-windows-amd64.exe`）

## 安装说明

### Linux/macOS用户

```bash
# 下载对应平台的压缩包
wget https://github.com/liangzhaoliang95/lxz/releases/download/v1.0.0/lxz-linux-amd64.tar.gz

# 解压
tar -xzf lxz-linux-amd64.tar.gz

# 移动到PATH目录
sudo mv lxz-linux-amd64 /usr/local/bin/lxz

# 添加执行权限
sudo chmod +x /usr/local/bin/lxz

# 验证安装
lxz --version
```

### Windows用户

1. 下载对应平台的exe文件
2. 将文件放在任意目录
3. 将该目录添加到PATH环境变量
4. 或在命令行中直接运行exe文件

## 故障排除

### 构建失败
- 检查Go版本是否兼容（需要Go 1.24.3+）
- 确保所有依赖都已正确安装
- 查看GitHub Actions日志获取详细错误信息

### Release未自动创建
- 确保tag名称以`v`开头（如：`v1.0.0`）
- 检查GitHub Actions是否成功完成
- 验证仓库有足够的权限创建releases

## 注意事项

- 构建过程大约需要5-10分钟
- 确保仓库有足够的GitHub Actions配额
- 建议在发布前先在本地测试构建过程
- 可以手动取消失败的构建并重新触发
