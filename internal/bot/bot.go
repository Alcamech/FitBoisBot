package bot

import (
	"log"
	"sync"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/Alcamech/FitBoisBot/internal/database/repository"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	lastMessageUserID    int64
	isFastestGGAvailable bool
	stateMutex           sync.Mutex

	bot      *tgbotapi.BotAPI
	userRepo *repository.UserRepository
)

func InitBot() {
	log.Println("FitBoisBot Started")
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.DevToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = config.AppConfig.Telegram.Debug

	userRepo = &repository.UserRepository{DB: database.DB}

	log.Printf("Authorized on account %s", bot.Self.UserName)

}

func BotLoop() {
	InitBot()

	var botConfig tgbotapi.UpdateConfig
	botConfig = tgbotapi.NewUpdate(0)
	botConfig.Timeout = 60

	var updates tgbotapi.UpdatesChannel
	updates = bot.GetUpdatesChan(botConfig)

	for update := range updates {
		if update.Message != nil {
			handleUpdate(bot, update)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := update.Message
	userID := msg.From.ID
	userName := msg.From.FirstName
	chatID := msg.Chat.ID

	user, err := userRepo.GetOrCreateUser(userID, userName, chatID)
	if err != nil {
		log.Printf("Error getting or creating user: %v", err)
		return
	}

	if msg.IsCommand() {
		handleCommand(bot, msg)
	}

	log.Printf("User: %v, Message: %s", user.ToString(), msg.Text)

	/*
		var msg tgbotapi.MessageConfig
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "go-rewrite: message received")
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	*/
}

func handleCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		onHelp(bot, msg.Chat.ID)
	case "fastgg":
		onFastGG(bot, msg.Chat.ID)
	}
}

func onHelp(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Nothing to see here... yet"
	sendText(bot, chatID, message)
}

func onFastGG(bot *tgbotapi.BotAPI, chatID int64) {
	message := "Fastest GG: (Sample Data)"
	sendText(bot, chatID, message)
}

func sendText(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}
