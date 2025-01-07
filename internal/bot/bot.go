package bot

import (
	"log"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"github.com/Alcamech/FitBoisBot/internal/database/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	lastMessageUserID    int64
	isFastestGGAvailable bool

	bot          *tgbotapi.BotAPI
	userRepo     *repository.UserRepository
	activityRepo *repository.ActivityRepository
	ggRepo       *repository.GGRepository
	groupRepo    *repository.GroupRepository
	tokenRepo    *repository.TokenRepository
)

var ggStates = make(map[int64]*ggState)

type ggState struct {
	isFastestGGAvailable bool
	activityPosterID     int64
}

func BotLoop() {
	log.Println("FitBoisBot Started")
	var err error
	bot, err = tgbotapi.NewBotAPI(config.AppConfig.Telegram.DevToken)
	if err != nil {
		log.Panicf("Failed to initialize bot: %v", err)
	}

	bot.Debug = config.AppConfig.Telegram.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	initializeRepositories()

	botConfig := tgbotapi.NewUpdate(0)
	botConfig.AllowedUpdates = []string{"message", "edited_message"}
	botConfig.Timeout = 60

	updates := bot.GetUpdatesChan(botConfig)

	processUpdates(updates)
}

func initializeRepositories() {
	userRepo = &repository.UserRepository{DB: database.DB}
	activityRepo = &repository.ActivityRepository{DB: database.DB}
	ggRepo = &repository.GGRepository{DB: database.DB}
	groupRepo = &repository.GroupRepository{DB: database.DB}
	tokenRepo = &repository.TokenRepository{DB: database.DB}
}

func processUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if msg := extractMessage(update); msg != nil {
			processMessage(bot, msg)
		}
	}
}

func extractMessage(update tgbotapi.Update) *tgbotapi.Message {
	if msg := update.EditedMessage; msg != nil {
		return msg
	}
	if msg := update.Message; msg != nil {
		return msg
	}
	return nil
}

func processMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	user, err := userRepo.GetOrCreateUser(userID, msg.From.FirstName, chatID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return
	}

	switch {
	case msg.IsCommand():
		handleCommand(bot, msg)
	case isPhoto(msg):
		// sendText(bot, chatID, "Photo detected.")
	case isActivity(msg):
		handleActivity(bot, msg, user)
	case isGG(msg.Text):
		handleGG(bot, msg)
	default:
		// sendText(bot, chatID, "Message not recognized.")
	}

	log.Printf("Group: %d, User: %s, Message: %s", chatID, user.ToString(), msg.Text)
}

func handleCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		onHelp(bot, msg.Chat.ID)
	case "fastgg":
		onFastGG(bot, msg.Chat.ID)
	case "tokens":
		onTokens(bot, msg.Chat.ID)
	case "timezone":
		onSetTimezone(bot, msg)
	}
}

func handleActivity(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User) {
	chatID := msg.Chat.ID

	ggStates[chatID] = &ggState{
		isFastestGGAvailable: true,
		activityPosterID:     user.ID,
	}

	if err := onActivityPost(msg, user); err != nil {
		log.Printf("Error posting activity: %v", err)
		return
	}

	if activityCountsMessage, err := getActivityCountsMessage(chatID); err == nil {
		sendText(bot, chatID, activityCountsMessage)
	} else {
		log.Printf("Error retrieving activity counts: %v", err)
	}
}

func handleGG(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	state, exists := ggStates[chatID]
	if !exists || !state.isFastestGGAvailable {
		return // GG is not available
	}

	// if msg.From.ID == state.activityPosterID {
	// 	return
	// }

	currentYear := GetCurrentYearInEST()
	if err := ggRepo.CreateOrUpdateGGCount(msg.From.ID, chatID, currentYear); err != nil {
		log.Printf("Failed to record GG: %v", err)
		return
	}

	sendReply(bot, chatID, "Fastest gg in the west", msg.MessageID)
	state.isFastestGGAvailable = false
}
