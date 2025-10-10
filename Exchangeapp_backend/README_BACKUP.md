# 数据库备份功能说明

## 功能概述

本项目已集成自动数据库备份功能，支持每周自动备份和手动备份。

## 功能特性

- ✅ 每周自动备份数据库（每周日凌晨2点执行）
- ✅ 手动触发备份
- ✅ 自动清理30天前的旧备份文件
- ✅ 备份文件管理
- ✅ RESTful API接口

## 文件结构

```
utils/
├── backup.go      # 备份核心功能
└── scheduler.go   # 定时任务调度器

controllers/
└── backup_controller.go  # 备份API控制器
```

## API接口

### 1. 手动触发备份
```http
POST /api/backup/manual
```

**响应示例：**
```json
{
  "success": true,
  "message": "数据库备份完成",
  "time": "2024-01-15 14:30:25"
}
```

### 2. 获取备份文件列表
```http
GET /api/backup/list
```

**响应示例：**
```json
{
  "success": true,
  "data": [
    "backup_20240115_143025.sql",
    "backup_20240108_140000.sql"
  ],
  "count": 2
}
```

### 3. 获取备份状态
```http
GET /api/backup/status
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "total_backups": 2,
    "latest_backup": "backup_20240115_143025.sql",
    "backup_files": [
      "backup_20240115_143025.sql",
      "backup_20240108_140000.sql"
    ]
  }
}
```

### 4. 清理旧备份
```http
DELETE /api/backup/clean?days=30
```

**参数：**
- `days`: 保留天数，默认为30天

## 备份文件存储

- 备份文件存储在 `./backups/` 目录下
- 文件命名格式：`backup_YYYYMMDD_HHMMSS.sql`
- 自动清理30天前的备份文件

## 定时任务配置

### 每周备份
- 执行时间：每周日凌晨2点
- 自动清理旧备份文件
- 支持优雅停止

### 自定义间隔备份（可选）
```go
// 启动每日备份
scheduler.StartDailyBackup()

// 启动自定义间隔备份（如每6小时）
scheduler.StartCustomBackup(6 * time.Hour)
```

## 数据库连接配置

备份功能会自动从配置文件 `config.yml` 中读取数据库连接信息：

```yaml
database:
  dsn: root:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local
```

## 系统要求

- MySQL数据库
- `mysqldump` 命令可用
- 足够的磁盘空间存储备份文件

## 日志记录

所有备份操作都会记录在应用日志中：

```
INFO[2024-01-15T14:30:25+08:00] 开始备份数据库: database_name
INFO[2024-01-15T14:30:30+08:00] 数据库备份完成: ./backups/backup_20240115_143025.sql
```

## 故障排除

### 1. 备份失败
- 检查数据库连接配置
- 确认 `mysqldump` 命令可用
- 检查磁盘空间是否充足
- 查看应用日志获取详细错误信息

### 2. 定时任务未执行
- 确认应用正常启动
- 检查系统时间设置
- 查看日志中的定时任务启动信息

### 3. 权限问题
- 确保应用有创建 `./backups/` 目录的权限
- 确保有执行 `mysqldump` 的权限

## 安全建议

1. **备份文件安全**：定期将备份文件复制到安全的存储位置
2. **访问控制**：考虑为备份API添加认证中间件
3. **加密存储**：对敏感数据的备份文件进行加密
4. **监控告警**：设置备份失败的监控和告警机制

## 扩展功能

### 1. 备份到云存储
可以扩展备份功能，将备份文件上传到云存储服务（如AWS S3、阿里云OSS等）。

### 2. 备份压缩
可以添加备份文件压缩功能，减少存储空间占用。

### 3. 增量备份
可以实现增量备份功能，只备份变更的数据。

### 4. 备份验证
可以添加备份文件完整性验证功能。

### 5. 使用SQL语句创建数据库
```sql
-- 删除数据库（如果存在）
DROP DATABASE IF EXISTS [数据库名字];

-- 创建数据库
CREATE DATABASE [数据库名字]
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE [数据库名字];
```