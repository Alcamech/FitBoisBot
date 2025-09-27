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

func (s *BotService) onHelp(chatID int64) {
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

func (s *BotService) onActivityPost(msg *tgbotapi.Message, user *models.User) error {
	activity, month, day, year, err := parseActivityMessage(msg)
	if err != nil {
		s.sendReply(msg.Chat.ID, constants.MsgActivityFormatError, msg.MessageID)
		return fmt.Errorf("failed to parse activity message: %w", err)
	}

	return s.activityStore.CreateOrUpdateRecord(user.ID, msg.Chat.ID, msg.MessageID, activity, month, day, year)
}
