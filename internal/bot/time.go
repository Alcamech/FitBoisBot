package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// GetPreviousMonth returns the previous month in MM format for the given timezone.
func GetPreviousMonth(timezone string) (string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Error("Invalid timezone", "error", err)
		return "", fmt.Errorf("failed to load timezone %s: %w", timezone, err)
	}
	now := time.Now().In(location)
	previousMonth := now.AddDate(0, -1, 0)
	return previousMonth.Format("01"), nil // MM format
}

// GetCurrentMonth returns the current month in MM format for the given timezone.
func GetCurrentMonth(timezone string) (string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Error("Invalid timezone", "error", err)
		return "", fmt.Errorf("failed to load timezone %s: %w", timezone, err)
	}
	now := time.Now().In(location)
	return now.Format("01"), nil // MM format
}

// GetCurrentYear returns the current year in YYYY format for the given timezone.
func GetCurrentYear(timezone string) (string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Error("Invalid timezone", "error", err)
		return "", fmt.Errorf("failed to load timezone %s: %w", timezone, err)
	}
	now := time.Now().In(location)
	return now.Format("2006"), nil // YYYY format
}

// convertToFullYear converts a 2 or 4 digit year to full YYYY format.
func convertToFullYear(yearStr string) (string, error) {
	yearValue, err := strconv.Atoi(yearStr)
	if err != nil {
		return "", fmt.Errorf("year must be numeric (got %q)", yearStr)
	}

	switch len(yearStr) {
	case 2:
		// Convert YY to YYYY by assuming current century
		currentYear := time.Now().Year()
		century := (currentYear / 100) * 100
		fullYear := century + yearValue
		return fmt.Sprintf("%d", fullYear), nil
	case 4:
		// Already full year
		return yearStr, nil
	default:
		return "", fmt.Errorf("year must be 2 or 4 digits (got %q)", yearStr)
	}
}