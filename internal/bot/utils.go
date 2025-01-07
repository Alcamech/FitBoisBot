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

func GetPreviousMonth(timezone string) (string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Invalid timezone: %v", err)
		return "", err
	}
	now := time.Now().In(location)
	previousMonth := now.AddDate(0, -1, 0)
	return previousMonth.Format("01"), nil // MM format
}

func GetCurrentMonth(timezone string) (string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Invalid timezone: %v", err)
		return "", err
	}
	now := time.Now().In(location)
	return now.Format("01"), nil // MM format
}

func GetCurrentYear(timezone string) (string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Invalid timezone: %v", err)
		return "", err
	}
	now := time.Now().In(location)
	return now.Format("2006"), nil // YYYY format
}

func convertToFullYear(year string) (string, error) {
	if !isNumeric(year) {
		return "", fmt.Errorf("invalid year format, expected 'YY' or 'YYYY'")
	}

	if len(year) == 2 {
		currentYear := time.Now().Year()
		century := currentYear / 100 * 100
		fullYear := century + atoi(year)
		return fmt.Sprintf("%d", fullYear), nil // Always assume the current century
	}

	if len(year) == 4 {
		return year, nil
	}

	return "", fmt.Errorf("invalid year format, expected 'YY' or 'YYYY'")
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

func formatTokenLeaderboard(leaderboard []models.Token, year string) string {
	if len(leaderboard) == 0 {
		return fmt.Sprintf("No tokens awarded yet for %s! 🏆", year)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("🏆 Token Leaderboard for %s 🏆\n", year))

	for i, entry := range leaderboard {
		user, _ := userRepo.FindByID(entry.UserID)
		builder.WriteString(fmt.Sprintf("%d. %s - %d tokens\n", i+1, user.Name, entry.Balance))
	}

	return builder.String()
}

func parseActivityMessage(msg *tgbotapi.Message) (string, string, string, string, error) {
	// Expected format: "activity-MM-DD-YYYY" or "activity-MM-DD-YY"
	parts := strings.Split(msg.Caption, "-")
	if len(parts) != 4 {
		return "", "", "", "", fmt.Errorf("invalid message format, expected 'activity-MM-DD-YYYY' or 'activity-MM-DD-YY'")
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

	year, err := convertToFullYear(year)
	if err != nil {
		return "", "", "", "", err
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
