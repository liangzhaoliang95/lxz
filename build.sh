#!/bin/zsh

binPath=""

# 判断当前操作系统 windows macos linux
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
    if [[ -d "/usr/local/bin" ]]; then
        binPath="/usr/local/bin/"
    elif [[ -d "/usr/bin" ]]; then
        binPath="/usr/bin/"
    else
        echo "No suitable bin directory found."
        exit 1
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
    binPath="/usr/local/bin/"
elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" ]]; then
    OS="windows"
    binPath="/c/shell/"
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

# 判断架构
if [[ "$(uname -m)" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$(uname -m)" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $(uname -m)"
    exit 1
fi

#macos
GOOS=$OS GOARCH=$ARCH go build -o lxz.exe ./main.go
sudo cp lxz.exe "$binPath"
sudo chmod +x "$binPath/lxz"
