package server

import (
	"exchangeapp/global"
	"exchangeapp/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func Countdown(TgBot *tgbotapi.BotAPI) {
	duration := utils.GetDuration(global.AppConfig.TgBot.Hour, global.AppConfig.TgBot.Min)
	// 使用一个无限循环进行倒计时
	for {
		hour := duration / time.Hour
		mine := duration / time.Minute % 60
		second := duration / time.Second % 60
		utils.Logger.Errorf("下班倒计时: %d小时%d分钟%d秒", hour, mine, second)
		if hour < 1 && mine < 10 {
			if second < 10 {
				time.Sleep(100 * time.Millisecond)
			} else {
				time.Sleep(5 * time.Second)
			}
			TgBot.Send(tgbotapi.NewMessage(global.AppConfig.TgBot.ChatID, fmt.Sprintf("下班倒计时: %d小时%d分钟%d秒", hour, mine, second)))
		} else {
			time.Sleep(5 * time.Second)
		}

		duration = duration - time.Second*1
		// 当倒计时结束时退出循环
		if duration <= 0 {
			fmt.Println("下班倒计时结束！")
			TgBot.Send(tgbotapi.NewMessage(global.AppConfig.TgBot.ChatID, "下班时间到，全体起立，离开工位"))
			Countdown(TgBot)
			break
		}
	}
}
