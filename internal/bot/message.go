package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendText(bot *tgbotapi.BotAPI, chatID int64, message string) {
	bot.Send(tgbotapi.NewMessage(chatID, message))
}

func sendReply(bot *tgbotapi.BotAPI, chatID int64, message string, replyToMessageID int) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = replyToMessageID
	bot.Send(msg)
}

func getActivityCountsMessage(chatID int64) (string, error) {
	month, err := getCurrentMonthInEST()
	if err != nil {
		return "", fmt.Errorf("failed to get current month: %w", err)
	}

	userIDs, err := activityRepo.GetUsersWithActivities(chatID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user activities: %w", err)
	}

	userCounts := make(map[string]int64, len(userIDs))
	for _, userID := range userIDs {
		count, err := activityRepo.GetActivityCountByUserIdAndMonth(userID, chatID, month)
		if err != nil {
			return "", fmt.Errorf("failed to fetch activity count for user %d: %w", userID, err)
		}

		user, err := userRepo.FindByID(userID)
		if err != nil {
			return "", fmt.Errorf("failed to fetch user %d: %w", userID, err)
		}

		userCounts[user.Name] = count
	}

	return formatActivityCounts(userCounts), nil
}

func formatActivityCounts(userCounts map[string]int64) string {
	if len(userCounts) == 0 {
		return "No activity recorded."
	}

	var builder strings.Builder
	builder.WriteString("Activity counts updated: ")

	for name, count := range userCounts {
		fmt.Fprintf(&builder, "%s=%d, ", name, count)
	}

	return strings.TrimSuffix(builder.String(), ", ")
}
