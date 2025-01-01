package bot

import (
	"fmt"
	"log"
	"strconv"
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
	// Expected format: "activity-MM-DD-YYYY"
	parts := strings.Split(msg.Caption, "-")
	if len(parts) != 4 {
		return "", "", "", "", fmt.Errorf("invalid message format, expected 'activity-MM-DD-YYYY'")
	}

	activity := parts[0]
	month := parts[1]
	day := parts[2]
	year := parts[3]

	if len(month) != 2 || !isNumeric(month) || atoi(month) < 1 || atoi(month) > 12 {
		return "", "", "", "", fmt.Errorf("invalid month format, expected 'MM'")
	}

	if len(day) != 2 || !isNumeric(day) || atoi(day) < 1 || atoi(day) > 31 {
		return "", "", "", "", fmt.Errorf("invalid day format, expected 'DD'")
	}

	if len(year) != 4 || !isNumeric(year) {
		return "", "", "", "", fmt.Errorf("invalid year format, expected 'YYYY'")
	}

	return activity, month, day, year, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func atoi(s string) int {
	value, _ := strconv.Atoi(s)
	return value
}
