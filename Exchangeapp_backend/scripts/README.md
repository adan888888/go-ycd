# MySQL（backup_mysql.sh） 备份脚本使用说明 （另外两个是用于mac电脑的定时任务，进入docker里面的命令进行的备份）

这个目录包含了 MySQL 数据库备份脚本，用于备份和恢复 MySQL 数据库。

## 文件说明

### `backup_mysql.sh` - MySQL 备份脚本
用于备份和恢复 MySQL 数据库，支持多种环境（Docker、本地、远程、云数据库）。采用的是直接在命令行里面，进行的命令进行的。

## 使用方法

### MySQL 备份脚本

```bash
# 基本备份
./scripts/backup_mysql.sh

# 压缩备份
./scripts/backup_mysql.sh --compress

# 列出所有备份文件
./scripts/backup_mysql.sh --list

# 恢复最新备份
./scripts/backup_mysql.sh --restore

# 设置保留天数（默认30天）
./scripts/backup_mysql.sh --days 7

# 显示帮助
./scripts/backup_mysql.sh --help
```

## 功能特性

### MySQL 备份脚本特性
- ✅ **直接使用 mysqldump** - 不依赖 Docker 命令
- ✅ **通用性强** - 适用于任何 MySQL 环境（Docker、本地、远程、云数据库）
- ✅ **支持压缩备份** - 节省存储空间
- ✅ **自动清理旧备份** - 可配置保留天数
- ✅ **备份文件列表** - 查看所有备份文件
- ✅ **一键恢复** - 恢复最新备份
- ✅ **安全确认** - 恢复前需要用户确认
- ✅ **彩色日志** - 清晰的输出信息

## 配置说明

### MySQL 配置
- **数据库名**: test
- **用户名**: root
- **密码**: mima123
- **主机地址**: 127.0.0.1
- **端口**: 3306
- **备份目录**: ../backups
- **默认保留天数**: 30天

## 定时任务设置

### 设置自动备份
```bash
# 编辑 crontab
crontab -e

# 每天凌晨2点压缩备份
0 2 * * * /path/to/scripts/backup_mysql.sh --compress

# 每周日凌晨3点全量备份，保留90天
0 3 * * 0 /path/to/scripts/backup_mysql.sh --compress --days 90
```

## 三种操作模式

### 1. 备份模式（默认）
```bash
./scripts/backup_mysql.sh
# 执行：备份 + 清理旧备份 + 显示备份列表
```

### 2. 恢复模式
```bash
./scripts/backup_mysql.sh --restore
# 执行：恢复最新备份
# 特点：需要用户确认，支持压缩文件
```

### 3. 列表模式
```bash
./scripts/backup_mysql.sh --list
# 执行：显示备份文件列表
# 显示：文件名、大小、修改时间
```

## 适用环境

### Docker 环境
```bash
DB_HOST="127.0.0.1"  # 连接 Docker 中的 MySQL
```

### 本地 MySQL
```bash
DB_HOST="127.0.0.1"  # 连接本地 MySQL
```

### 远程服务器
```bash
DB_HOST="192.168.1.100"  # 连接远程服务器
```

### 云数据库
```bash
DB_HOST="rds.aliyun.com"  # 连接阿里云 RDS
DB_HOST="cdb.tencent.com" # 连接腾讯云 CDB
```

## 实际命令示例

### 备份命令
```bash
# 普通备份
mysqldump -h 127.0.0.1 -P 3306 -u root -pmima123 --single-transaction --routines --triggers test > backup.sql

# 压缩备份
mysqldump -h 127.0.0.1 -P 3306 -u root -pmima123 --single-transaction --routines --triggers test | gzip > backup.sql.gz
```

### 恢复命令
```bash
# 普通文件恢复
mysql -h 127.0.0.1 -P 3306 -u root -pmima123 test < backup.sql

# 压缩文件恢复
gunzip -c backup.sql.gz | mysql -h 127.0.0.1 -P 3306 -u root -pmima123 test
```

## 注意事项

1. **权限**: 确保脚本有执行权限 (`chmod +x`)
2. **路径**: 脚本需要在项目根目录执行
3. **MySQL 客户端**: 确保已安装 mysqldump 和 mysql 命令
4. **网络连接**: 确保能够连接到 MySQL 服务器
5. **备份**: 定期检查备份文件是否正常生成

## 故障排除

### 常见问题

1. **mysqldump 命令未找到**
   ```bash
   # Mac 安装 MySQL 客户端
   brew install mysql-client
   
   # Ubuntu/Debian
   sudo apt-get install mysql-client
   
   # CentOS/RHEL
   sudo yum install mysql
   ```

2. **无法连接 MySQL**
   ```bash
   # 检查 MySQL 服务状态
   docker ps | grep mysql-test
   
   # 检查端口是否开放
   netstat -an | grep 3306
   
   # 测试连接
   mysql -h 127.0.0.1 -P 3306 -u root -pmima123 -e "SELECT 1;"
   ```

3. **权限不足**
   ```bash
   # 添加执行权限
   chmod +x scripts/*.sh
   ```

4. **备份失败**
   ```bash
   # 检查 MySQL 连接
   mysqldump -h 127.0.0.1 -P 3306 -u root -pmima123 --single-transaction --no-data test
   
   # 检查备份目录权限
   ls -la ../backups/
   ```

## 更新日志

- **v1.0** - 初始版本，基本功能
- **v1.1** - 添加压缩备份和恢复功能
- **v1.2** - 添加自动清理和列表功能
- **v1.3** - 添加彩色日志和错误处理
- **v1.4** - 改为直接使用 mysqldump，支持多种环境
