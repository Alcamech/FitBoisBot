package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendText(bot *tgbotapi.BotAPI, chatID int64, message string) {
	bot.Send(tgbotapi.NewMessage(chatID, message))
}

func sendMarkdownText(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func sendReply(bot *tgbotapi.BotAPI, chatID int64, message string, replyToMessageID int) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = replyToMessageID
	bot.Send(msg)
}

func getActivityCountsMessage(chatID int64) (string, error) {
	group, err := groupRepo.GetGroupByID(chatID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch group %d: %w", chatID, err)
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = "America/New_York"
	}

	month, err := GetCurrentMonth(timezone)
	if err != nil {
		return "", fmt.Errorf("failed to get current month for timezone %s: %w", timezone, err)
	}

	year, err := GetCurrentYear(timezone)
	if err != nil {
		return "", fmt.Errorf("failed to get current year for timezone %s: %w", timezone, err)
	}

	userIDs, err := activityRepo.GetUsersWithActivities(chatID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user activities: %w", err)
	}

	userCounts := make(map[string]int64, len(userIDs))
	for _, userID := range userIDs {
		count, err := activityRepo.GetActivityCountByUserIdAndMonthAndYear(userID, chatID, month, year)
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

func sendMonthlyAwardMessage(bot *tgbotapi.BotAPI, chatID int64, userName, month, year string, rewardAmount int) {
	message := fmt.Sprintf(
		"Month counts have been reset\n\nCongratulations ⭐️ %s ⭐️ for being the most active user for %s/%s 🏆\n\nHere's your reward. 💰 You've won %d FitBoi Tokens! 💰",
		userName, month, year, rewardAmount,
	)
	sendText(bot, chatID, message)
}
