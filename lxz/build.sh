#!/bin/zsh

#macos
GOOS=darwin GOARCH=arm64 go build -o lxz ./main.go
sudo cp lxz /usr/local/bin/lxz
sudo chmod +x /usr/local/bin/lxz
