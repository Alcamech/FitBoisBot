package bot

import (
	"fmt"
	"strings"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// parseActivityMessage parses activity message in format "activity-MM-DD-YYYY" or "activity-MM-DD-YY".
func parseActivityMessage(msg *tgbotapi.Message) (string, string, string, string, error) {
	// Expected format: "activity-MM-DD-YYYY" or "activity-MM-DD-YY"
	parts := strings.Split(msg.Caption, "-")
	if len(parts) != 4 {
		return "", "", "", "", fmt.Errorf("invalid format: expected 'activity-MM-DD-YYYY' or 'activity-MM-DD-YY'")
	}

	activity, month, day, yearStr := parts[0], parts[1], parts[2], parts[3]

	// Validate activity name
	if !isValidActivityName(activity) {
		return "", "", "", "", fmt.Errorf("invalid activity name: must be 1-%d characters", constants.ActivityNameMaxLength)
	}

	// Validate month
	if err := validateDatePart(month, constants.MinMonth, constants.MaxMonth, "month"); err != nil {
		return "", "", "", "", err
	}

	// Validate day
	if err := validateDatePart(day, constants.MinDay, constants.MaxDay, "day"); err != nil {
		return "", "", "", "", err
	}

	// Convert year to full format
	year, err := convertToFullYear(yearStr)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid year: %w", err)
	}

	return activity, month, day, year, nil
}

// isActivity checks if message is an activity post (photo with caption).
func isActivity(msg *tgbotapi.Message) bool {
	// Must have both photo and caption
	if msg.Photo == nil || msg.Caption == "" {
		return false
	}

	// Check if caption follows activity format (basic validation)
	return strings.Contains(msg.Caption, "-")
}

// isGG checks if message is a "GG" response.
func isGG(msgText string) bool {
	// Allow case-insensitive "gg" responses
	return strings.ToLower(strings.TrimSpace(msgText)) == "gg"
}

// isPhoto checks if message is just a photo without caption.
func isPhoto(msg *tgbotapi.Message) bool {
	return (msg.Photo != nil && msg.Caption == "")
}

