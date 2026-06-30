#!/usr/bin/env bash
# 将 Docker 前端入口 (默认 8080) 通过 ngrok 暴露到公网。
# H5 在容器内经 Nginx 反代 /api，无需单独指 ngrok API；Flutter 真机需公网 API 地址。
set -euo pipefail

PORT="${1:-8080}"
NGROK_CONFIG="${HOME}/Library/Application Support/ngrok/ngrok.yml"

if ! command -v ngrok >/dev/null 2>&1; then
  echo "未找到 ngrok，请先安装：brew install ngrok/ngrok/ngrok"
  exit 1
fi

if ! curl -s -o /dev/null --connect-timeout 2 "http://127.0.0.1:${PORT}/"; then
  echo "本地 ${PORT} 端口无响应。请先启动 Docker 前后端，例如："
  echo "  cd Exchangeapp_backend && docker compose up -d --build"
  echo "  cd Exchangeapp_frontend && docker compose up -d --build"
  exit 1
fi

if [[ ! -f "${NGROK_CONFIG}" ]]; then
  echo "ngrok 尚未配置 authtoken（只需做一次）。"
  echo ""
  echo "1. 注册：https://dashboard.ngrok.com/signup"
  echo "2. 复制 token：https://dashboard.ngrok.com/get-started/your-authtoken"
  echo "3. 执行：ngrok config add-authtoken <你的token>"
  echo ""
  echo "配置完成后重新运行本脚本。"
  exit 1
fi

echo "Docker 前端已就绪 (http://127.0.0.1:${PORT})"
echo "启动 ngrok 隧道（暴露 Nginx，H5 + /api 一体）..."
echo ""
echo "启动后看终端 Forwarding 里的 https 地址，手机浏览器直接打开即可。"
echo "ngrok 本地管理面板：http://127.0.0.1:4040"
echo "按 Ctrl+C 停止"
echo ""

exec ngrok http "${PORT}"
