package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"myserver/config"
	"myserver/rebot"
	"os"
	"os/signal"
)

func main() {
	config.AppConfig = &config.Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../config") //./表示main所在的文件夹 ../main上一层的文件夹

	//读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	//解释成结构体
	if err := viper.Unmarshal(config.AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	go rebot.TgRobot()
	fmt.Println("启动机器人...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt) //os.Interrupt:ctrl+c , os.Kill:杀死进程
	<-quit
	fmt.Println("机器人停止 ...")

}
