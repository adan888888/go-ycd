#!/bin/bash

# Mac 唤醒/开盖时执行一次数据库备份
# 该脚本由 SleepWatcher 在系统唤醒时调用

set -euo pipefail

PATH="/usr/local/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin:${PATH:-}"

BACKUP_SCRIPT="/Users/a123123/GolandProjects/Web003Gin-01_gingormtutorials/Exchangeapp_backend/scripts/backup_once_to_documents.sh"
LOG_FILE="/Users/a123123/Documents/test_db_backup_wakeup.log"

mkdir -p "/Users/a123123/Documents"

{
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] 检测到 Mac 唤醒，开始执行一次备份"
  /bin/bash "${BACKUP_SCRIPT}"
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] 唤醒备份执行完成"
  echo ""
} >> "${LOG_FILE}" 2>&1
