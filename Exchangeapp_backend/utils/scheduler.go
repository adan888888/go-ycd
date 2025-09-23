package utils

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	stopChan chan bool
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		stopChan: make(chan bool),
	}
}

// StartWeeklyBackup 启动每周备份任务
func (s *Scheduler) StartWeeklyBackup() {
	log.Info("启动每周数据库备份定时任务")

	// 计算到下一个周日的时间
	now := time.Now()
	nextSunday := s.getNextSunday(now)

	// 计算等待时间
	waitDuration := nextSunday.Sub(now)
	log.Infof("下次备份时间: %s (等待 %v)", nextSunday.Format("2006-01-02 15:04:05"), waitDuration)

	// 启动定时器
	ticker := time.NewTicker(3 * 24 * time.Hour) // 每3天执行一次

	// 立即执行一次备份（可选）
	go func() {
		log.Info("执行初始数据库备份")
		if err := BackupDatabase(); err != nil {
			log.Errorf("初始备份失败: %v", err)
		} else {
			log.Info("初始备份完成")
		}
	}()

	// 等待到下一个周日
	time.Sleep(waitDuration)

	// 开始定期备份
	for {
		select {
		case <-ticker.C: //定时器通道 监听器 每7天触发一次
			log.Info("开始执行定期数据库备份")
			if err := BackupDatabase(); err != nil {
				log.Errorf("数据库备份失败: %v", err)
			} else {
				log.Info("定期备份完成")
			}

			// 清理旧备份文件
			if err := CleanOldBackups(); err != nil {
				log.Errorf("清理旧备份失败: %v", err)
			}

		case <-s.stopChan: //停止管道监听器
			ticker.Stop() // 手动停止定时器
			log.Info("停止定时备份任务")
			return
		}
	}
}

// StartDailyBackup 启动每日备份任务（可选）
func (s *Scheduler) StartDailyBackup() {
	log.Info("启动每日数据库备份定时任务")

	// 计算到下一个凌晨2点的时间
	now := time.Now()
	nextBackup := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if nextBackup.Before(now) {
		nextBackup = nextBackup.Add(24 * time.Hour)
	}

	waitDuration := nextBackup.Sub(now)
	log.Infof("下次备份时间: %s (等待 %v)", nextBackup.Format("2006-01-02 15:04:05"), waitDuration)

	// 等待到指定时间
	time.Sleep(waitDuration)

	// 启动每日定时器
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Info("开始执行每日数据库备份")
			if err := BackupDatabase(); err != nil {
				log.Errorf("数据库备份失败: %v", err)
			} else {
				log.Info("每日备份完成")
			}

		case <-s.stopChan:
			log.Info("停止每日备份任务")
			return
		}
	}
}

// StartCustomBackup 启动自定义间隔备份任务
func (s *Scheduler) StartCustomBackup(interval time.Duration) {
	log.Infof("启动自定义间隔备份任务，间隔: %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Info("开始执行自定义间隔数据库备份")
			if err := BackupDatabase(); err != nil {
				log.Errorf("数据库备份失败: %v", err)
			} else {
				log.Info("自定义间隔备份完成")
			}

			// 清理旧备份文件
			if err := CleanOldBackups(); err != nil {
				log.Errorf("清理旧备份失败: %v", err)
			}

		case <-s.stopChan:
			log.Info("停止自定义间隔备份任务")
			return
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	select {
	case s.stopChan <- true:
		 // 成功发送停止信号
		 log.Info("发送停止信号成功")
	default:
		 // 通道已满，忽略
		 log.Info("通道已满，忽略")
	}
}

// getNextSunday 获取下一个周日的时间
func (s *Scheduler) getNextSunday(now time.Time) time.Time {
	// 计算距离下一个周日的天数
	daysUntilSunday := (7 - int(now.Weekday())) % 7
	if daysUntilSunday == 0 {
		daysUntilSunday = 7 // 如果今天是周日，则等到下周日
	}

	// 计算下一个周日的日期
	nextSunday := now.AddDate(0, 0, daysUntilSunday)

	// 设置为凌晨2点
	return time.Date(nextSunday.Year(), nextSunday.Month(), nextSunday.Day(), 2, 0, 0, 0, nextSunday.Location())
}

// GetNextBackupTime 获取下次备份时间
func (s *Scheduler) GetNextBackupTime() time.Time {
	now := time.Now()
	return s.getNextSunday(now)
}
