package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
)

// getUserFirstName retrieves the first name of a user by ID.
func (s *BotService) getUserFirstName(userID int64) string {
	user, err := s.userStore.FindByID(userID)
	if err != nil {
		slog.Error("Failed to fetch user", "error", err, "user_id", userID)
		return "Unknown"
	}
	return user.Name
}

// formatFastGGLeaderboard formats the Fast GG leaderboard for display.
func (s *BotService) formatFastGGLeaderboard(leaderboard []models.Gg) string {
	if len(leaderboard) == 0 {
		return constants.MsgNoFastGGsRecorded
	}

	var builder strings.Builder
	builder.WriteString("⚡ <b>Fastest GG Leaderboard</b>\n\n")

	// Simple ranked list format
	for i, entry := range leaderboard {
		name := s.getUserFirstName(entry.UserID)
		builder.WriteString(fmt.Sprintf("%d. %s: <b>%d</b> GGs\n", i+1, name, entry.FastGGCount))
	}

	return builder.String()
}

// formatTokenLeaderboard formats the token leaderboard for display.
func (s *BotService) formatTokenLeaderboard(leaderboard []models.Token) string {
	if len(leaderboard) == 0 {
		return fmt.Sprint(constants.NoTokensTemplate)
	}

	var builder strings.Builder
	builder.WriteString("🏆 <b>Token Leaderboard </b>\n\n")

	// Simple ranked list format
	for i, entry := range leaderboard {
		user, _ := s.userStore.FindByID(entry.UserID)
		if user == nil {
			continue
		}
		builder.WriteString(fmt.Sprintf("%d. %s: <b>%d</b> tokens\n", i+1, user.Name, entry.Balance))
	}

	return builder.String()
}

// formatActivityCounts formats activity counts for display.
func formatActivityCounts(userCounts map[string]int64) string {
	if len(userCounts) == 0 {
		return constants.MsgNoActivitiesRecorded
	}

	var builder strings.Builder
	builder.WriteString("📊 <b>Monthly Activity Leaderboard</b>\n\n")

	// Simple clean format
	for name, count := range userCounts {
		builder.WriteString(fmt.Sprintf("%s: <b>%d</b>\n", name, count))
	}

	return builder.String()
}

// formatLeaderboardHeader creates a standard header for leaderboards
func formatLeaderboardHeader(icon, title string) string {
	return fmt.Sprintf("%s <b>%s</b>\n\n", icon, title)
}

// formatRankedEntry creates a standard ranked entry for leaderboards
func formatRankedEntry(rank int, name string, value interface{}, unit string) string {
	return fmt.Sprintf("%d. %s: <b>%v</b> %s\n", rank, name, value, unit)
}
