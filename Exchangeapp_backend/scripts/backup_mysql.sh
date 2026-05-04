#!/bin/bash

# MySQL 备份脚本
# 使用方法: ./backup_mysql.sh [选项]
# 选项:
#   -c, --compress    压缩备份文件
#   -d, --days N      保留N天的备份
#   -l, --list        列出备份文件
#   -r, --restore     恢复备份
#   -h, --help        显示帮助

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
DB_NAME="test"
DB_USER="root"
DB_PASS="mima123"
DB_HOST="127.0.0.1"
DB_PORT="3306"
BACKUP_DIR="../backups"
RETENTION_DAYS=30

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "MySQL 备份脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -c, --compress    压缩备份文件"
    echo "  -d, --days N      保留N天的备份（默认30天）"
    echo "  -l, --list        列出所有备份文件"
    echo "  -r, --restore     恢复最新的备份"
    echo "  -h, --help        显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                # 基本备份"
    echo "  $0 --compress     # 压缩备份"
    echo "  $0 --list         # 列出备份"
    echo "  $0 --restore      # 恢复备份"
}

# 检查 MySQL 客户端工具
check_mysql_tools() {
    if ! command -v mysqldump &> /dev/null; then
        log_error "mysqldump 命令未找到"
        log_info "请安装 MySQL 客户端: brew install mysql-client"
        exit 1
    fi

    if ! command -v mysql &> /dev/null; then
        log_error "mysql 命令未找到"
        log_info "请安装 MySQL 客户端: brew install mysql-client"
        exit 1
    fi
}

# 检查 MySQL 服务是否可连接
check_mysql_server() {
    if ! mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS -e "SELECT 1;" > /dev/null 2>&1; then
        log_error "无法连接到 MySQL 服务器"
        log_info "请检查 MySQL 服务是否运行: $DB_HOST:$DB_PORT"
        exit 1
    fi
}

# 检查备份源数据库是否存在
check_source_database() {
    if ! mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS --single-transaction --no-data $DB_NAME > /dev/null 2>&1; then
        log_error "无法访问数据库: $DB_NAME"
        log_info "请确认数据库存在且账号有权限访问"
        exit 1
    fi
}

# 确保恢复目标数据库存在
ensure_target_database() {
    if mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS \
        -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\`;" > /dev/null 2>&1; then
        log_info "已确认目标数据库存在: $DB_NAME"
    else
        log_error "创建目标数据库失败: $DB_NAME"
        exit 1
    fi
}

# 执行备份
backup_database() {
    local compress=$1
    
    # 创建备份目录
    mkdir -p $BACKUP_DIR
    
    # 生成备份文件名
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$BACKUP_DIR/backup_${timestamp}.sql"
    
    if [ "$compress" = true ]; then
        backup_file="${backup_file}.gz"
    fi
    
    log_info "开始备份数据库: $DB_NAME"
    log_info "备份文件: $backup_file"
    
    # 执行备份
    if [ "$compress" = true ]; then
        # 压缩备份
        if mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS --single-transaction --routines --triggers --databases $DB_NAME | gzip > $backup_file; then
            log_success "压缩备份成功: $backup_file"
        else
            log_error "压缩备份失败"
            exit 1
        fi
    else
        # 普通备份
        if mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS --single-transaction --routines --triggers --databases $DB_NAME > $backup_file; then
            log_success "备份成功: $backup_file"
        else
            log_error "备份失败"
            exit 1
        fi
    fi
    
    # 显示备份文件信息
    if [ -f "$backup_file" ]; then
        local file_size=$(du -h "$backup_file" | cut -f1)
        log_info "备份文件大小: $file_size"
    fi
}

# 列出备份文件
list_backups() {
    log_info "备份文件列表:"
    
    if [ ! -d "$BACKUP_DIR" ] || [ -z "$(ls -A $BACKUP_DIR 2>/dev/null)" ]; then
        log_warning "没有找到备份文件"
        return
    fi
    
    echo ""
    printf "%-30s %-10s %-20s\n" "文件名" "大小" "修改时间"
    printf "%-30s %-10s %-20s\n" "------------------------------" "----------" "--------------------"
    
    for file in $BACKUP_DIR/backup_*.sql*; do
        if [ -f "$file" ]; then
            local filename=$(basename "$file")
            local filesize=$(du -h "$file" | cut -f1)
            local modtime=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" "$file")
            printf "%-30s %-10s %-20s\n" "$filename" "$filesize" "$modtime"
        fi
    done
}

# 恢复备份
restore_backup() {
    log_info "查找最新的备份文件..."
    
    # 查找最新的备份文件
    local latest_backup=$(ls -t $BACKUP_DIR/backup_*.sql* 2>/dev/null | head -n1)
    
    if [ -z "$latest_backup" ]; then
        log_error "没有找到备份文件"
        exit 1
    fi
    
    log_info "找到最新备份: $latest_backup"
    
    # 确认恢复
    read -p "确定要恢复这个备份吗？这将覆盖当前数据库！(y/N): " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        log_info "取消恢复操作"
        exit 0
    fi
    
    log_warning "开始恢复数据库..."

    # 如果目标数据库不存在，先自动创建
    ensure_target_database
    
    # 执行恢复
    if [[ "$latest_backup" == *.gz ]]; then
        # 压缩文件恢复
        if gunzip -c "$latest_backup" | mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME; then
            log_success "数据库恢复成功"
        else
            log_error "数据库恢复失败"
            exit 1
        fi
    else
        # 普通文件恢复
        if mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME < "$latest_backup"; then
            log_success "数据库恢复成功"
        else
            log_error "数据库恢复失败"
            exit 1
        fi
    fi
}

# 清理旧备份
clean_old_backups() {
    local days=$1
    
    log_info "清理 $days 天前的备份文件..."
    
    if [ ! -d "$BACKUP_DIR" ]; then
        log_warning "备份目录不存在"
        return
    fi
    
    # 查找并删除旧文件
    local deleted_count=0
    for file in $BACKUP_DIR/backup_*.sql*; do
        if [ -f "$file" ]; then
            local file_age=$(($(date +%s) - $(stat -f "%m" "$file")))
            local days_seconds=$((days * 24 * 60 * 60))
            
            if [ $file_age -gt $days_seconds ]; then
                log_info "删除旧备份: $(basename "$file")"
                rm "$file"
                ((deleted_count++))
            fi
        fi
    done
    
    if [ $deleted_count -gt 0 ]; then
        log_success "已删除 $deleted_count 个旧备份文件"
    else
        log_info "没有需要删除的旧备份文件"
    fi
}

# 主函数
main() {
    local compress=false
    local list_only=false
    local restore_only=false
    local days=$RETENTION_DAYS
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -c|--compress)
                compress=true
                shift
                ;;
            -d|--days)
                days="$2"
                shift 2
                ;;
            -l|--list)
                list_only=true
                shift
                ;;
            -r|--restore)
                restore_only=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 执行相应操作
    if [ "$list_only" = true ]; then
        list_backups
    elif [ "$restore_only" = true ]; then
        check_mysql_tools
        check_mysql_server
        restore_backup
    else
        check_mysql_tools
        check_mysql_server
        check_source_database
        # 执行备份
        backup_database $compress
        
        # 清理旧备份
        clean_old_backups $days
        
        # 显示备份列表
        echo ""
        list_backups
    fi
}

# 执行主函数
main "$@"