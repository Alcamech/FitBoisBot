package bot

import (
	"github.com/Alcamech/FitBoisBot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StartBot() {
	log.Println("FitBoisBot Started")
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.DevToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var botConfig tgbotapi.UpdateConfig
	botConfig = tgbotapi.NewUpdate(0)
	botConfig.Timeout = 60

	var updates tgbotapi.UpdatesChannel
	updates = bot.GetUpdatesChan(botConfig)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var msg tgbotapi.MessageConfig
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "go-rewrite: message received")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
