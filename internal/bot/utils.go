package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func isActivity(msg *tgbotapi.Message) bool {
	return (msg.Photo != nil && msg.Caption != "")
}

func isGG(msgText string) bool {
	return strings.ToLower(msgText) == "gg"
}

func isPhoto(msg *tgbotapi.Message) bool {
	return (msg.Photo != nil && msg.Caption == "")
}

func isEdited(msg *tgbotapi.Message) bool {
	return msg.EditDate != 0
}

func isActivityValidFormat(messageCaption string) bool {
	// format "activity-MM-DD-YYYY"
	parts := strings.Split(messageCaption, "-")

	return len(parts) == 4
}

func getTimeInEST() (time.Time, error) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone: %v", err)
	}
	return time.Now().In(location), nil
}
