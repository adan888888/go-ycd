package controllers

import (
	"exchangeapp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// BackupController 备份控制器
type BackupController struct{}

// NewBackupController 创建备份控制器
func NewBackupController() *BackupController {
	return &BackupController{}
}

// ManualBackup 手动触发备份
func (bc *BackupController) ManualBackup(c *gin.Context) {
	log.Info("收到手动备份请求")

	if err := utils.BackupDatabase(); err != nil {
		log.Errorf("手动备份失败: %v", err)
		ServerFailMsg(c, "备份失败: "+err.Error())
		return
	}

	OkMsg(c, "数据库备份完成", gin.H{
		"time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetBackupList 获取备份文件列表
func (bc *BackupController) GetBackupList(c *gin.Context) {
	backupFiles, err := utils.GetBackupList()
	if err != nil {
		log.Errorf("获取备份列表失败: %v", err)
		ServerFailMsg(c, "获取备份列表失败: "+err.Error())
		return
	}

	OkMsg(c, "查询成功", gin.H{
		"list":  backupFiles,
		"count": len(backupFiles),
	})
}

// CleanOldBackups 清理旧备份
func (bc *BackupController) CleanOldBackups(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		FailMsg(c, "无效的天数参数")
		return
	}

	if days < 1 {
		FailMsg(c, "保留天数必须大于0")
		return
	}

	if err := utils.CleanOldBackups(); err != nil {
		log.Errorf("清理旧备份失败: %v", err)
		ServerFailMsg(c, "清理旧备份失败: "+err.Error())
		return
	}

	OkMsg(c, "旧备份清理完成", gin.H{"retain_days": days})
}

// GetBackupStatus 获取备份状态信息
func (bc *BackupController) GetBackupStatus(c *gin.Context) {
	backupFiles, err := utils.GetBackupList()
	if err != nil {
		log.Errorf("获取备份状态失败: %v", err)
		ServerFailMsg(c, "获取备份状态失败: "+err.Error())
		return
	}

	var latestBackup string
	if len(backupFiles) > 0 {
		latestBackup = backupFiles[len(backupFiles)-1]
	}

	OkMsg(c, "查询成功", gin.H{
		"total_backups": len(backupFiles),
		"latest_backup": latestBackup,
		"backup_files":  backupFiles,
	})
}
