package controllers

import (
	"exchangeapp/utils"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "备份失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据库备份完成",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetBackupList 获取备份文件列表
func (bc *BackupController) GetBackupList(c *gin.Context) {
	backupFiles, err := utils.GetBackupList()
	if err != nil {
		log.Errorf("获取备份列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取备份列表失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    backupFiles,
		"count":   len(backupFiles),
	})
}

// CleanOldBackups 清理旧备份
func (bc *BackupController) CleanOldBackups(c *gin.Context) {
	// 获取保留天数参数，默认为30天
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的天数参数",
		})
		return
	}
	
	if days < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "保留天数必须大于0",
		})
		return
	}
	
	// 这里可以扩展为支持自定义保留天数的清理逻辑
	if err := utils.CleanOldBackups(); err != nil {
		log.Errorf("清理旧备份失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "清理旧备份失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "旧备份清理完成",
		"retain_days": days,
	})
}

// GetBackupStatus 获取备份状态信息
func (bc *BackupController) GetBackupStatus(c *gin.Context) {
	// 获取备份文件列表
	backupFiles, err := utils.GetBackupList()
	if err != nil {
		log.Errorf("获取备份状态失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取备份状态失败: " + err.Error(),
		})
		return
	}
	
	// 计算备份文件总大小（暂时不实现）
	// totalSize := int64(0)
	// for _, fileName := range backupFiles {
	// 	// 这里可以添加获取文件大小的逻辑
	// }
	
	// 获取最新的备份文件
	var latestBackup string
	if len(backupFiles) > 0 {
		latestBackup = backupFiles[len(backupFiles)-1] // 假设按时间排序
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_backups": len(backupFiles),
			"latest_backup": latestBackup,
			"backup_files":  backupFiles,
		},
	})
}
