package config

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() {
	// 获取数据库配置
	dsn := global.AppConfig.Database.Dsn
	if dsn == "" {
		log.Fatal("数据库连接字符串不能为空")
	}

	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 连接数据库
	var err error
	global.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err := global.Db.DB()
	if err != nil {
		log.Fatal("获取数据库实例失败:", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(global.AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	fmt.Println("数据库连接成功")

	// 自动迁移数据库表结构
	err = global.Db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.ExchangeRate{},
		&models.Banner{},
		&models.TableYanchendao1{},
		&models.TableYanchendao2{},
		&models.BuyRecord{},     // 添加买币记录表的自动迁移
		&models.PasswordItem{},  // 添加密码本表的自动迁移
	)
	if err != nil {
		log.Fatal("数据库表迁移失败:", err)
	}

	fmt.Println("数据库表迁移完成")
}
