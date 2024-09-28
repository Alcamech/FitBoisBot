package bot

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/Alcamech/FitBoisBot/internal/database/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	lastMessageUserID    int64
	isFastestGGAvailable bool
	stateMutex           sync.Mutex

	bot          *tgbotapi.BotAPI
	userRepo     *repository.UserRepository
	activityRepo *repository.ActivityRepository
)

func BotLoop() {
	log.Println("FitBoisBot Started")
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.DevToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = config.AppConfig.Telegram.Debug

	userRepo = &repository.UserRepository{DB: database.DB}
	activityRepo = &repository.ActivityRepository{DB: database.DB}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var botConfig tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	botConfig.AllowedUpdates = []string{"message", "edited_message"}
	botConfig.Timeout = 60

	var updates tgbotapi.UpdatesChannel = bot.GetUpdatesChan(botConfig)

	for update := range updates {
		if update.EditedMessage != nil {
			handleUpdate(bot, update)
		} else if update.Message != nil {
			handleUpdate(bot, update)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := getMessageFromUpdate(update)

	if msg == nil {
		return
	}

	from := msg.From
	chat := msg.Chat
	userId, userName, chatID := from.ID, from.FirstName, chat.ID

	user, err := userRepo.GetOrCreateUser(userId, userName, chatID)
	if err != nil {
		log.Printf("Error getting or creating user: %v", err)
		return
	}

	if isPhoto(msg) {
		sendText(bot, chatID, "go-rewrite: photo detected.")
	} else if isActivity(msg) && isActivityValidFormat(msg.Caption) {
		sendText(bot, chatID, "go-rewrite: activity detected.")
		if err := onActivityPost(msg, user); err == nil {
			activityCountsMessage, err := getActivityCountsMessage()
			if err == nil {
				sendText(bot, chatID, activityCountsMessage)
			} else {
				log.Printf("Error getting activity counts: %v", err)
			}
		} else {
			log.Printf("Error posting activity: %v", err)
		}

	} else if isGG(msg.Text) {
		sendText(bot, chatID, "go-rewrite: gg detected.")
	} else {
		sendText(bot, chatID, "go-rewrite: non-activity detected.")
	}

	if msg.IsCommand() {
		handleCommand(bot, msg)
	}

	log.Printf("User: %v, Message: %s", user.ToString(), msg.Text)
}

func getMessageFromUpdate(update tgbotapi.Update) *tgbotapi.Message {
	var msg *tgbotapi.Message

	if update.EditedMessage != nil {
		msg = update.EditedMessage
	} else if update.Message != nil {
		msg = update.Message

		if msg.Photo != nil {
			stateMutex.Lock()
			lastMessageUserID = msg.From.ID
			isFastestGGAvailable = true
			stateMutex.Unlock()
		}
	}

	return msg
}

func handleCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		onHelp(bot, msg.Chat.ID)
	case "fastgg":
		onFastGG(bot, msg.Chat.ID)
	}
}

func parseActivityMessage(msg *tgbotapi.Message) (string, string, string, string, error) {
	// format "activity-MM-DD-YYYY" or "activity-MMDDYYYY"
	parts := strings.Split(msg.Caption, "-")
	if len(parts) < 4 {
		return "", "", "", "", fmt.Errorf("invalid message format")
	}

	activity := parts[0]
	month := parts[1]
	day := parts[2]
	year := parts[3]

	return activity, month, day, year, nil
}
