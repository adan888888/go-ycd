package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"exchangeapp/global"
	log "github.com/sirupsen/logrus"
)

// BackupDatabase 执行数据库备份
func BackupDatabase() error {
	// 创建备份目录
	backupDir := "./backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}

	// 生成备份文件名（包含时间戳）
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.sql", timestamp))

	// 从DSN中解析数据库连接信息
	dsn := global.AppConfig.Database.Dsn
	host, user, password, database := parseDSN(dsn)

	// 构建mysqldump命令
	cmd := exec.Command("mysqldump",
		"-h", host,
		"-u", user,
		"-p"+password,
		"--single-transaction",
		"--routines",
		"--triggers",
		database)

	// 创建备份文件
	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %v", err)
	}
	defer file.Close()

	// 设置命令输出
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	// 执行备份命令
	log.Infof("开始备份数据库: %s", database)
	if err := cmd.Run(); err != nil {
		// 如果备份失败，删除空文件
		os.Remove(backupFile)
		return fmt.Errorf("数据库备份失败: %v", err)
	}

	log.Infof("数据库备份完成: %s", backupFile)
	return nil
}

// parseDSN 解析数据库连接字符串
func parseDSN(dsn string) (host, user, password, database string) {
	// 解析格式: user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
	
	// 提取用户和密码部分
	parts := strings.Split(dsn, "@")
	if len(parts) < 2 {
		return "localhost", "root", "", "test"
	}
	
	userPass := parts[0]
	tcpPart := parts[1]
	
	// 解析用户名和密码
	userPassParts := strings.Split(userPass, ":")
	if len(userPassParts) >= 2 {
		user = userPassParts[0]
		password = userPassParts[1]
	} else {
		user = userPass
		password = ""
	}
	
	// 解析主机和数据库
	// 查找tcp(host:port)部分
	tcpStart := strings.Index(tcpPart, "tcp(")
	if tcpStart != -1 {
		tcpEnd := strings.Index(tcpPart[tcpStart:], ")")
		if tcpEnd != -1 {
			hostPort := tcpPart[tcpStart+4 : tcpStart+tcpEnd]
			hostParts := strings.Split(hostPort, ":")
			if len(hostParts) > 0 {
				host = hostParts[0]
			}
		}
	}
	
	// 解析数据库名
	dbStart := strings.Index(tcpPart, "/")
	if dbStart != -1 {
		dbEnd := strings.Index(tcpPart[dbStart:], "?")
		if dbEnd != -1 {
			database = tcpPart[dbStart+1 : dbStart+dbEnd]
		} else {
			database = tcpPart[dbStart+1:]
		}
	}
	
	// 设置默认值
	if host == "" {
		host = "localhost"
	}
	if user == "" {
		user = "root"
	}
	if database == "" {
		database = "test"
	}
	
	return host, user, password, database
}

// CleanOldBackups 清理旧的备份文件（保留最近30天的备份）
func CleanOldBackups() error {
	backupDir := "./backups"
	
	// 读取备份目录
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("读取备份目录失败: %v", err)
	}
	
	// 计算30天前的时间
	cutoffTime := time.Now().AddDate(0, 0, -30)
	
	// 删除30天前的备份文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// 检查文件名是否匹配备份文件格式
		if !strings.HasPrefix(file.Name(), "backup_") || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		
		// 获取文件信息
		fileInfo, err := file.Info()
		if err != nil {
			log.Warnf("无法获取文件信息: %s", file.Name())
			continue
		}
		
		// 如果文件超过30天，删除它
		if fileInfo.ModTime().Before(cutoffTime) {
			filePath := filepath.Join(backupDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Warnf("删除旧备份文件失败: %s, 错误: %v", filePath, err)
			} else {
				log.Infof("已删除旧备份文件: %s", filePath)
			}
		}
	}
	
	return nil
}

// GetBackupList 获取备份文件列表
func GetBackupList() ([]string, error) {
	backupDir := "./backups"
	
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %v", err)
	}
	
	var backupFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "backup_") && strings.HasSuffix(file.Name(), ".sql") {
			backupFiles = append(backupFiles, file.Name())
		}
	}
	
	return backupFiles, nil
}
