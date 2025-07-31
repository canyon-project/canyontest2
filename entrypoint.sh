#!/bin/sh

# 容器启动脚本
set -e

# 打印启动信息
echo "Starting application..."
echo "Current time: $(date)"

# 启动后端应用
exec /app/main