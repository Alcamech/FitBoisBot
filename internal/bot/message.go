package bot

import (
	"fmt"
	"log/slog"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func (s *BotService) sendText(chatID int64, message string) {
	if _, err := s.bot.Send(tgbotapi.NewMessage(chatID, message)); err != nil {
		telegramErr := errors.NewTelegramError("sendMessage", chatID, err)
		slog.Error("Failed to send text message", "error", telegramErr, "chat_id", chatID)
	}
}

func (s *BotService) sendMarkdownText(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "MarkdownV2"
	if _, err := s.bot.Send(msg); err != nil {
		telegramErr := errors.NewTelegramError("sendMessage", chatID, err)
		slog.Error("Failed to send markdown message", "error", telegramErr, "chat_id", chatID)
	}
}

func (s *BotService) sendHTMLText(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	if _, err := s.bot.Send(msg); err != nil {
		telegramErr := errors.NewTelegramError("sendMessage", chatID, err)
		slog.Error("Failed to send HTML message", "error", telegramErr, "chat_id", chatID)
	}
}

func (s *BotService) sendReply(chatID int64, message string, replyToMessageID int) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = replyToMessageID
	if _, err := s.bot.Send(msg); err != nil {
		telegramErr := errors.NewTelegramError("sendMessage", chatID, err)
		slog.Error("Failed to send reply message", "error", telegramErr, "chat_id", chatID, "reply_to", replyToMessageID)
	}
}

func (s *BotService) getActivityCountsMessage(chatID int64) (string, error) {
	group, err := s.groupStore.GetByID(chatID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch group %d: %w", chatID, err)
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = constants.DefaultTimezone
	}

	month, err := GetCurrentMonth(timezone)
	if err != nil {
		return "", fmt.Errorf("failed to get current month for timezone %s: %w", timezone, err)
	}

	year, err := GetCurrentYear(timezone)
	if err != nil {
		return "", fmt.Errorf("failed to get current year for timezone %s: %w", timezone, err)
	}

	userIDs, err := s.activityStore.GetUsersWithActivities(chatID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user activities: %w", err)
	}

	userCounts := make(map[string]int64, len(userIDs))
	for _, userID := range userIDs {
		count, err := s.activityStore.GetCountByUserMonthYear(userID, chatID, month, year)
		if err != nil {
			return "", fmt.Errorf("failed to fetch activity count for user %d: %w", userID, err)
		}

		user, err := s.userStore.FindByID(userID)
		if err != nil {
			return "", fmt.Errorf("failed to fetch user %d: %w", userID, err)
		}

		userCounts[user.Name] = count
	}

	return formatActivityCounts(userCounts), nil
}


func (s *BotService) sendMonthlyAwardMessage(chatID int64, userName, month, year string, rewardAmount int) {
	message := fmt.Sprintf(constants.MonthlyAwardTemplate, userName, month, year, rewardAmount)
	s.sendText(chatID, message)
}
