#!/bin/bash

# MySQL 一次性备份脚本
# 用途:
#   手动执行一次数据库备份，将压缩后的 .sql.gz 文件保存到 /Users/a123123/Documents
#   备份文件命名规则为 backup_YYYYMMDD_HHMMSS.sql.gz
#
# 使用命令:
#   /bin/bash "/Users/a123123/GolandProjects/Web003Gin-01_gingormtutorials/Exchangeapp_backend/scripts/backup_once_to_documents.sh"
#
# 说明:
#   1. 该脚本只执行“备份一次”，不会自动恢复数据
#   2. 该脚本会进入 Docker 容器 mysql-test 内执行 mysqldump
#   3. 备份文件保存在 Mac 本机的 Documents 目录，不会堆积在容器内部

set -euo pipefail

BACKUP_DIR="/Users/a123123/Documents"
CONTAINER_NAME="mysql-test"
DB_NAME="test"
DB_USER="root"
DB_PASS="mima123"

TIMESTAMP="$(date '+%Y%m%d_%H%M%S')"
BACKUP_FILE="${BACKUP_DIR}/backup_${TIMESTAMP}.sql.gz"

mkdir -p "${BACKUP_DIR}"

echo "开始备份数据库 ${DB_NAME} ..."

if ! command -v docker >/dev/null 2>&1; then
  echo "错误: 未找到 docker 命令"
  exit 1
fi

if ! docker inspect -f '{{.State.Running}}' "${CONTAINER_NAME}" >/dev/null 2>&1; then
  echo "错误: 未找到容器 ${CONTAINER_NAME}"
  exit 1
fi

if [ "$(docker inspect -f '{{.State.Running}}' "${CONTAINER_NAME}")" != "true" ]; then
  echo "错误: 容器 ${CONTAINER_NAME} 当前未运行"
  exit 1
fi

if docker exec "${CONTAINER_NAME}" mysqldump \
  -u"${DB_USER}" \
  -p"${DB_PASS}" \
  --single-transaction \
  --routines \
  --triggers \
  "${DB_NAME}" | gzip > "${BACKUP_FILE}"; then
  FILE_SIZE="$(du -h "${BACKUP_FILE}" | awk '{print $1}')"
  echo "备份成功: ${BACKUP_FILE}"
  echo "文件大小: ${FILE_SIZE}"
else
  rm -f "${BACKUP_FILE}"
  echo "错误: 备份失败"
  exit 1
fi
