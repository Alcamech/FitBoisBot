package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
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

func GetCurrentMonthInEST() string {
	est := time.FixedZone("EST", -5*60*60) // Fixed offset for EST (UTC-5)
	now := time.Now().In(est)
	return now.Format("01") // MM format
}

func GetCurrentYearInEST() string {
	est := time.FixedZone("EST", -5*60*60) // Fixed offset for EST (UTC-5)
	now := time.Now().In(est)
	return now.Format("2006") // YYYY format
}

func getUserFirstName(userID int64) string {
	user, err := userRepo.FindByID(userID)
	if err != nil {
		log.Printf("Failed to fetch user %d: %v", userID, err)
		return "Unknown"
	}
	return user.Name
}

func formatFastGGLeaderboard(leaderboard []models.Gg) string {
	if len(leaderboard) == 0 {
		return "No Fast GGs recorded yet! 🏆"
	}

	var builder strings.Builder
	builder.WriteString("Fastest GGs 😎 ")

	for i, entry := range leaderboard {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%s=%d", getUserFirstName(entry.UserID), entry.FastGGCount))
	}

	return builder.String()
}

func parseActivityMessage(msg *tgbotapi.Message) (string, string, string, string, error) {
	// format "activity-MM-DD-YYYY"
	parts := strings.Split(msg.Caption, "-")
	if len(parts) < 4 {
		fmt.Println("< 4")
		return "", "", "", "", fmt.Errorf("invalid message format")
	}

	activity := parts[0]
	month := parts[1]
	day := parts[2]
	year := parts[3]

	return activity, month, day, year, nil
}
