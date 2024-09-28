package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendText(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func getActivityCountsMessage() (string, error) {
	nowInEST, err := getTimeInEST()
	if err != nil {
		return "", err
	}

	month := nowInEST.Format("01") // MM in Go is "01"

	userIDs, err := activityRepo.GetUsersWithActivities()
	if err != nil {
		return "", fmt.Errorf("failed to fetch record count: %v", err)
	}

	userCounts := make(map[string]int64)

	for _, userID := range userIDs {
		count, err := activityRepo.GetActivityCountByUserIdAndMonth(userID, month)
		if err != nil {
			return "", fmt.Errorf("failed to fetch record count: %v", err)
		}

		user, err := userRepo.FindByID(userID)
		if err != nil {
			return "", fmt.Errorf("failed to fetch user: %v", err)
		}

		userCounts[user.Name] = count
	}

	var content strings.Builder
	for name, count := range userCounts {
		content.WriteString(fmt.Sprintf("%s=%d, ", name, count))
	}

	contentStr := strings.TrimRight(content.String(), ", ")

	return fmt.Sprintf("Activity counts updated: %s", contentStr), nil
}
