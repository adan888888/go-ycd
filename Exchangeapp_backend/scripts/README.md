# scripts 备份脚本说明

本目录包含 MySQL 数据库备份相关脚本。

## 文件一览

| 文件 | 用途 |
|------|------|
| `backup_mysql.sh` | 通用备份/恢复，直连 MySQL，适合 Docker、本地、远程 |
| `backup_once_to_documents.sh` | Mac 一次性备份，输出到 `~/Documents`；也被 launchd 定时和唤醒脚本调用 |
| `backup_on_wakeup.sh` | Mac 开盖/唤醒时触发，内部调用 `backup_once_to_documents.sh` |
| `com.a123123.test-db-backup.daily.plist` | Mac `launchd` 配置，每天 0:00 调用 `backup_once_to_documents.sh` |

---

## `backup_mysql.sh`

通过 `mysqldump` 直连 MySQL，不依赖 `docker exec`。在项目根目录执行。

**默认配置**（可在脚本内修改）：

| 项 | 值 |
|----|-----|
| 数据库 | test |
| 用户/密码 | root / mima123 |
| 地址/端口 | 127.0.0.1:3306 |
| 备份目录 | `../backups` |
| 保留天数 | 30 |

**常用命令**：

```bash
./scripts/backup_mysql.sh                  # 备份 + 清理旧文件 + 列出列表
./scripts/backup_mysql.sh --compress       # 压缩备份
./scripts/backup_mysql.sh --list           # 仅列出备份文件
./scripts/backup_mysql.sh --restore        # 恢复最新备份（需确认）
./scripts/backup_mysql.sh --days 7         # 指定保留天数
./scripts/backup_mysql.sh --help           # 查看帮助
```

**定时备份（crontab 示例）**：

```bash
# 每天凌晨 2 点压缩备份
0 2 * * * /path/to/scripts/backup_mysql.sh --compress

# 每周日凌晨 3 点备份，保留 90 天
0 3 * * 0 /path/to/scripts/backup_mysql.sh --compress --days 90
```

**注意**：需安装 `mysqldump` 和 `mysql`（Mac: `brew install mysql-client`）；脚本需有执行权限（`chmod +x scripts/*.sh`）。

---

## Mac 自动备份（开发机）

三个文件配合使用，备份写入 `~/Documents/backup_*.sql.gz`：

| 触发方式 | 文件 | 说明 |
|----------|------|------|
| 手动 | `backup_once_to_documents.sh` | 进入容器 `mysql-test` 执行 mysqldump |
| 每天 0:00 | `com.a123123.test-db-backup.daily.plist` | 调用 `backup_once_to_documents.sh`；复制到 `~/Library/LaunchAgents/` 后 `launchctl load` |
| 开盖/唤醒 | `backup_on_wakeup.sh` | 由 SleepWatcher 调用 `backup_once_to_documents.sh` |

```bash
# 手动执行一次
/bin/bash scripts/backup_once_to_documents.sh
```

日志：`~/Documents/test_db_backup_daily.out.log`、`~/Documents/test_db_backup_wakeup.log`
