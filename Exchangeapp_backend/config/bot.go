package config

import (
	"exchangeapp/global"
	"exchangeapp/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func initTgBot() {
	token := AppConfig.TgBot.Token
	chatID := int64(-4500054486)
	var bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	msg := tgbotapi.NewMessage(chatID, "你好")
	// 发送消息
	_, err = bot.Send(msg)

	log.Printf("Message sent to chat ID %d", chatID)
	if err != nil {
		panic(err)
	}
	// 设置更新配置
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 获取更新通道
	updates := bot.GetUpdatesChan(u)

	// 处理更新
	for update := range updates {
		if update.Message == nil { // 忽略任何非消息更新
			continue
		}
		// 打印收到的消息
		log.Printf("收到消息==>[%s] %s", update.Message.From.UserName, update.Message.Text)

		// 检查消息是否提到了机器人
		if update.Message.IsCommand() || strings.Contains(update.Message.Text, "@"+bot.Self.UserName) {
			//if strings.HasPrefix(update.Message.Text, "@bx_xia_Bot") {
			// 回复消息
			responseText := "你提到我了吗？我在这里！大佬请指教！"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			switch update.Message.Command() {
			case "start":
				msg.Text = "Hello! I am your friendly Telegram bot."
			case "help":
				msg.Text = "You can control me by sending these commands:\n/start - to start the bot\n/help - to get this help message"
			default:
				if utils.IsNumber(update.Message.Command()) {
					number, _ := strconv.Atoi(update.Message.Command())
					if number >= 1 && number <= 100 {
						msg.Text = update.Message.Command()
					} else {
						msg.Text = "是数字"
					}
				} else {
					if strings.Contains(update.Message.Text, "@bx_xia_Bot ") {
						msg.Text = "不要@我，我没空..."
					} else {
						msg.Text = "请重新输入..."
					}
				}
			}
			msg.ReplyToMessageID = update.Message.MessageID
			// 发送回复消息
			bot.Send(msg)

		}
	}
	global.TgBot = bot
}
