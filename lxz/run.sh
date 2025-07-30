#!/bin/bash

# 监控目录
WATCH_DIR="./"

# Go 主程序入口
GO_FILE="main.go"

# 轮询时间（秒）
INTERVAL=10

# 存储上一次目录状态的 Hash
PREV_HASH=""

# 后台Go进程的 PID
GO_PID=0

# 启动 Go 程序函数
start_go_process() {
  echo "🚀 启动 go run $GO_FILE ..."
  go run "$GO_FILE" &
  GO_PID=$!
  echo "Go 程序已启动，PID=$GO_PID"
}

# 停止 Go 程序函数
stop_go_process() {
  if [ $GO_PID -ne 0 ] && kill -0 $GO_PID 2>/dev/null; then
    echo "🛑 停止上一个 Go 进程，PID=$GO_PID"
    kill $GO_PID
    wait $GO_PID 2>/dev/null
  fi
  GO_PID=0
}

echo "🎯 开始监控目录: $WATCH_DIR"
echo "💡 检测变化后自动重启 Go 程序"

while true; do
  CURRENT_HASH=$(find "$WATCH_DIR" -type f -exec md5sum {} + 2>/dev/null | sort | md5sum)

  if [[ "$CURRENT_HASH" != "$PREV_HASH" ]]; then
    echo "📦 检测到文件变化"

    stop_go_process
    start_go_process

    PREV_HASH="$CURRENT_HASH"
  fi

  sleep $INTERVAL
done
