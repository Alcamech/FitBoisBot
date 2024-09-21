package bot

import (
	"log"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/database/repository"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BotLoop() {
	log.Println("FitBoisBot Started")
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.DevToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = config.AppConfig.Telegram.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var botConfig tgbotapi.UpdateConfig
	botConfig = tgbotapi.NewUpdate(0)
	botConfig.Timeout = 60

	var updates tgbotapi.UpdatesChannel
	updates = bot.GetUpdatesChan(botConfig)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			messageText := update.Message.Text
			userID := update.Message.From.ID
			userName := update.Message.From.FirstName
			groupID := update.Message.Chat.ID

			user, err := repository.GetOrCreateUser(userID, userName, groupID)

			if err != nil {
				log.Printf("Error handling user: %v", err)
			}

			log.Printf("Message: %s", messageText)
			log.Printf("User: %+v", user)

			/*
				var msg tgbotapi.MessageConfig
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "go-rewrite: message received")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			*/
		}
	}
}
