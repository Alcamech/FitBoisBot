package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"github.com/Alcamech/FitBoisBot/internal/version"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *BotService) onFastGG(chatID int64) {
	group, err := s.groupStore.GetByID(chatID)
	if err != nil {
		slog.Error("Failed to fetch group", "error", err, "chat_id", chatID)
		s.sendText(chatID, constants.MsgGroupTimezoneError)
		return
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = constants.DefaultTimezone
	}

	year, err := GetCurrentYear(timezone)
	if err != nil {
		slog.Error("Failed to get current year", "error", err, "timezone", timezone)
		s.sendText(chatID, constants.MsgInvalidTimezone)
		return
	}

	leaderboard, err := s.ggStore.GetLeaderboard(chatID, year)
	if err != nil {
		slog.Error("Failed to fetch GG leaderboard", "error", err, "chat_id", chatID)
		return
	}

	message := s.formatFastGGLeaderboard(leaderboard)
	s.sendHTMLText(chatID, message)
}

func (s *BotService) onHelp(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	args := strings.ToLower(strings.TrimSpace(msg.CommandArguments()))

	// Check for topic-specific help
	if args == "challenge" || args == "c" {
		s.sendHTMLText(chatID, constants.MsgChallengeHelp)
		return
	}

	// Default help
	versionInfo := strings.ReplaceAll(version.GetVersionInfo(), ".", "\\.")
	helpText := fmt.Sprintf(constants.MsgHelpText, versionInfo)
	s.sendMarkdownText(chatID, helpText)
}

func (s *BotService) onSetTimezone(msg *tgbotapi.Message) {
	args := msg.CommandArguments()

	// No arguments - show current timezone
	if args == "" {
		s.showCurrentTimezone(msg.Chat.ID)
		return
	}

	// Arguments provided - set new timezone
	if err := validateTimezone(args); err != nil {
		s.sendText(msg.Chat.ID, constants.MsgTimezoneInvalid)
		return
	}

	err := s.groupStore.SetTimezone(msg.Chat.ID, args)
	if err != nil {
		slog.Error("Failed to set timezone", "error", err, "chat_id", msg.Chat.ID)
		s.sendText(msg.Chat.ID, constants.MsgTimezoneSetFailed)
		return
	}

	s.sendText(msg.Chat.ID, constants.MsgTimezoneSetSuccess+args)
}

func (s *BotService) showCurrentTimezone(chatID int64) {
	group, err := s.groupStore.GetByID(chatID)
	if err != nil {
		slog.Error("Failed to fetch group", "error", err, "chat_id", chatID)
		s.sendText(chatID, constants.MsgGroupTimezoneError)
		return
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = constants.DefaultTimezone
	}

	s.sendText(chatID, fmt.Sprintf("Current timezone: %s", timezone))
}

func (s *BotService) onTokens(chatID int64) {
	group, err := s.groupStore.GetByID(chatID)
	if err != nil {
		slog.Error("Failed to fetch group", "error", err, "chat_id", chatID)
		s.sendText(chatID, constants.MsgGroupTimezoneError)
		return
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = constants.DefaultTimezone
	}

	year, err := GetCurrentYear(timezone)
	if err != nil {
		slog.Error("Failed to get current year", "error", err, "timezone", timezone)
		s.sendText(chatID, constants.MsgInvalidTimezone)
		return
	}

	leaderboard, err := s.tokenStore.GetYearlyLeaderboard(chatID, year)
	if err != nil {
		slog.Error("Failed to fetch leaderboard", "error", err, "chat_id", chatID)
		s.sendText(chatID, constants.MsgLeaderboardFailed)
		return
	}

	message := s.formatTokenLeaderboard(leaderboard)
	s.sendHTMLText(chatID, message)
}

func (s *BotService) onBalance(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	balance, err := s.tokenStore.GetUserLifetimeBalance(userID, chatID)
	if err != nil {
		slog.Error("Failed to fetch balance", "error", err, "user_id", userID)
		s.sendText(chatID, "Failed to fetch balance.")
		return
	}

	name := msg.From.FirstName
	s.sendHTMLText(chatID, fmt.Sprintf("<b>%s</b>: %d tokens", name, balance))
}

func (s *BotService) onLeaderboard(chatID int64) {
	message, err := s.getActivityCountsMessage(chatID)
	if err != nil {
		slog.Error("Failed to fetch leaderboard", "error", err, "chat_id", chatID)
		s.sendText(chatID, constants.MsgLeaderboardFailed)
		return
	}

	s.sendHTMLText(chatID, message)
}

func (s *BotService) onActivityPost(msg *tgbotapi.Message, user *models.User) (bool, error) {
	post, err := parseActivityMessage(msg)
	if err != nil {
		s.sendReply(msg.Chat.ID, constants.MsgActivityFormatError, msg.MessageID)
		return false, fmt.Errorf("failed to parse activity message: %w", err)
	}

	// Try to update existing record first (for edited captions)
	err = s.activityStore.UpdateRecord(*post)
	if err == nil {
		// Record was updated, not a new activity
		return false, nil
	}

	// Record doesn't exist, create new one
	err = s.activityStore.CreateRecord(*post)
	if err != nil {
		return false, err
	}

	// New activity was created
	return true, nil
}
