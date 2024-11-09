package rebot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"myserver/config"
	"time"
)

func StartTimer(tgBot *tgbotapi.BotAPI) {
	duration := GetDuration(config.AppConfig.TgBot.Hour, config.AppConfig.TgBot.Min, config.AppConfig.TgBot.Sec)
	// 使用一个无限循环进行倒计时
	for {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		//logrus.Infof("下班倒计时:%v %d小时%d分钟%d秒", duration, duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
		time.Sleep(2 * time.Second)
		duration = duration - 2*time.Second
		/**
		   if duration/time.Hour < 1 && duration/time.Minute%60 < 1 && duration/time.Second%60 <= 3 {
				for i := 0; i < 3; i++ {
					if duration/time.Second%60 >= 0 {
						tgBot.Send(tgbotapi.NewMessage(config.AppConfig.tgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)))
					}
					time.Sleep(100 * time.Millisecond)
					duration = duration - 1*time.Second
				}
			}

			// 当倒计时结束时退出循环
			if duration < 0 {
				tgBot.Send(tgbotapi.NewMessage(config.AppConfig.tgBot.ChatID, "下班时间到，全体起立，离开工位"))
				time.Sleep(5 * time.Second)
				startTimer(tgBot)
				break
			}
		*/

		///优化方案
		// 当倒计时结束时退出for循环
		if duration < 0 { //duration < time.Second  时间还可以这样对比
			go sendMsg(tgBot)
			StartTimer(tgBot)
			break
		}
	}
}

func sendMsg(tgBot *tgbotapi.BotAPI) {
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Millisecond * time.Duration(3-i) * 50)
		tgBot.Send(tgbotapi.NewMessage(config.AppConfig.TgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", 0, 0, 3-i)))
	}
	tgBot.Send(tgbotapi.NewMessage(config.AppConfig.TgBot.ChatID, "下班时间到，全体起立，离开工位"))
}
