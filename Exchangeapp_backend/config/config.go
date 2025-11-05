package config

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	//尝试读取配置文件。如果找到指定名称和类型的配置文件，viper 会将其内容加载到内存中。
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	
	// 启用配置文件监听，自动重新加载配置（无需重启）
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已更新: %s，配置已自动重新加载", e.Name)
	})
	
	//初始化全局配置结构体
	global.AppConfig = &models.Config{}
	//将配置信息解析到结构体中
	if err := viper.Unmarshal(global.AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	//InitRabbitMQ()
	utils.NewLogger() //初始化log
	initDB()
	initRedis()
	/*已迁移致新的服务myserver*/
	//go server.TgRobot() //初始化机器人
}
