package bot

import (
	"fmt"
	"log/slog"

	"github.com/Alcamech/FitBoisBot/config"
	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"github.com/Alcamech/FitBoisBot/internal/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotService encapsulates all bot dependencies and state
type BotService struct {
	bot           *tgbotapi.BotAPI
	userStore     *store.UserStore
	activityStore *store.ActivityStore
	ggStore       *store.GGStore
	groupStore    *store.GroupStore
	tokenStore    *store.TokenStore
	ggStates      map[int64]*ggState
}

type ggState struct {
	isFastestGGAvailable bool
	activityPosterID     int64
}

func BotLoop() {
	slog.Info("FitBoisBot started")
	
	service, err := NewBotService()
	if err != nil {
		slog.Error("Failed to initialize bot service", "error", err)
		panic(err)
	}

	slog.Info("Bot authorized", "username", service.bot.Self.UserName)

	botConfig := tgbotapi.NewUpdate(0)
	botConfig.AllowedUpdates = []string{"message", "edited_message"}
	botConfig.Timeout = 60

	updates := service.bot.GetUpdatesChan(botConfig)

	service.processUpdates(updates)
}

// NewBotService creates a new bot service with initialized repositories
func NewBotService() (*BotService, error) {
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bot API: %w", err)
	}

	bot.Debug = config.AppConfig.Debug

	return &BotService{
		bot:           bot,
		userStore:     store.NewUserStore(database.DB),
		activityStore: store.NewActivityStore(database.DB),
		ggStore:       store.NewGGStore(database.DB),
		groupStore:    store.NewGroupStore(database.DB),
		tokenStore:    store.NewTokenStore(database.DB),
		ggStates:      make(map[int64]*ggState),
	}, nil
}

func (s *BotService) processUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if msg := extractMessage(update); msg != nil {
			s.processMessage(msg)
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

func (s *BotService) processMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	user, err := s.userStore.GetOrCreateUser(userID, msg.From.FirstName, chatID)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err, "user_id", userID, "chat_id", chatID)
		return
	}

	switch {
	case msg.IsCommand():
		s.handleCommand(msg)
	case isPhoto(msg):
	case isActivity(msg):
		s.handleActivity(msg, user)
	case isGG(msg.Text):
		s.handleGG(msg)
	default:
	}

	slog.Info("Message processed", "chat_id", chatID, "user", user.ToString(), "message", msg.Text)
}

func (s *BotService) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		s.onHelp(msg.Chat.ID)
	case "fastgg":
		s.onFastGG(msg.Chat.ID)
	case "tokens":
		s.onTokens(msg.Chat.ID)
	case "timezone":
		s.onSetTimezone(msg)
	}
}

func (s *BotService) handleActivity(msg *tgbotapi.Message, user *models.User) {
	chatID := msg.Chat.ID

	s.ggStates[chatID] = &ggState{
		isFastestGGAvailable: true,
		activityPosterID:     user.ID,
	}

	if err := s.onActivityPost(msg, user); err != nil {
		slog.Error("Failed to post activity", "error", err)
		return
	}

	if activityCountsMessage, err := s.getActivityCountsMessage(chatID); err == nil {
		s.sendHTMLText(chatID, activityCountsMessage)
	} else {
		slog.Error("Failed to retrieve activity counts", "error", err)
	}
}

func (s *BotService) handleGG(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	state, exists := s.ggStates[chatID]
	if !exists || !state.isFastestGGAvailable {
		return // GG is not available
	}

	// Users cannot GG their own activity
	if msg.From.ID == state.activityPosterID {
		return // Ignore self-GG attempts
	}

	// Get group timezone for proper year calculation
	group, err := s.groupStore.GetByID(chatID)
	if err != nil {
		slog.Error("Failed to fetch group for GG", "error", err, "chat_id", chatID)
		return
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = constants.DefaultTimezone
	}

	currentYear, err := GetCurrentYear(timezone)
	if err != nil {
		slog.Error("Failed to get current year for GG", "error", err, "timezone", timezone)
		return
	}

	if err := s.ggStore.CreateOrUpdateCount(msg.From.ID, chatID, currentYear); err != nil {
		slog.Error("Failed to record GG", "error", err)
		return
	}

	s.sendReply(chatID, constants.MsgFastestGG, msg.MessageID)
	state.isFastestGGAvailable = false
}

// GetAllGroups returns all groups from the database
func (s *BotService) GetAllGroups() ([]models.Group, error) {
	return s.groupStore.GetAll()
}

// SendAnnouncement sends a message to a specific group
func (s *BotService) SendAnnouncement(groupID int64, message string) error {
	_, err := s.bot.Send(tgbotapi.NewMessage(groupID, message))
	return err
}
