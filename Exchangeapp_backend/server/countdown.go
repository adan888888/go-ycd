package server

import (
	"exchangeapp/global"
	"exchangeapp/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func Countdown(TgBot *tgbotapi.BotAPI) {
	duration := utils.GetDuration(global.AppConfig.TgBot.Hour, global.AppConfig.TgBot.Min, global.AppConfig.TgBot.Sec)
	// 使用一个无限循环进行倒计时
	for {
		//utils.Logger.Errorf("下班倒计时:%v %d小时%d分钟%d秒", duration, duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
		time.Sleep(2 * time.Second)
		duration = duration - 2*time.Second
		if duration/time.Hour < 1 && duration/time.Minute%60 < 1 && duration/time.Second%60 <= 3 {
			for i := 0; i < 3; i++ {
				if duration/time.Second%60 >= 0 {
					TgBot.Send(tgbotapi.NewMessage(global.AppConfig.TgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)))
				}
				time.Sleep(100 * time.Millisecond)
				duration = duration - 1*time.Second
			}
		}

		// 当倒计时结束时退出循环
		if duration <= 0 {
			TgBot.Send(tgbotapi.NewMessage(global.AppConfig.TgBot.ChatID, "下班时间到，全体起立，离开工位"))
			time.Sleep(5 * time.Second)
			Countdown(TgBot)
			break
		}
	}
}
