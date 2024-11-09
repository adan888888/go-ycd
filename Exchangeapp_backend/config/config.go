package config

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	global.AppConfig = &models.Config{}

	if err := viper.Unmarshal(global.AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	utils.NewLogger() //初始化log
	initDB()
	initRedis()
	/*已迁移致新的服务myserver*/
	//go server.TgRobot() //初始化机器人
}
