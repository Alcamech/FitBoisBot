package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
)

// formatChallengeDetails formats a challenge for detailed view.
func (s *BotService) formatChallengeDetails(challenge *models.Challenge) string {
	var builder strings.Builder

	// Title
	builder.WriteString(fmt.Sprintf("<b>%s</b>\n", challenge.Title))

	if challenge.Description != "" {
		builder.WriteString(fmt.Sprintf("<i>%s</i>\n", challenge.Description))
	}

	builder.WriteString("\n")

	// Status with time remaining if active/pending
	statusBadge := formatStatusBadge(challenge.Status)
	statusText := formatStatusText(challenge.Status)
	timeRemaining := s.getTimeRemaining(challenge)
	if timeRemaining != "" {
		builder.WriteString(fmt.Sprintf("Status: %s %s (%s)\n", statusBadge, statusText, timeRemaining))
	} else {
		builder.WriteString(fmt.Sprintf("Status: %s %s\n", statusBadge, statusText))
	}

	// Difficulty
	difficultyBadge := formatDifficultyBadge(challenge.Difficulty)
	builder.WriteString(fmt.Sprintf("Difficulty: %s %s (%.1fx Multiplier)\n", difficultyBadge, challenge.Difficulty, challenge.Multiplier))

	// Creator info
	creatorName := s.getUserFirstName(challenge.CreatorID)
	builder.WriteString(fmt.Sprintf("Creator: %s\n", creatorName))

	// Participants
	builder.WriteString("\n<b>Participants:</b>\n")
	if len(challenge.Participants) == 0 {
		builder.WriteString("No participants yet\n")
	} else {
		for _, p := range challenge.Participants {
			name := s.getUserFirstName(p.UserID)
			potential := float64(p.WagerAmount) + float64(p.WagerAmount)*challenge.Multiplier
			winnerMark := ""
			if p.IsWinner != nil {
				if *p.IsWinner {
					winnerMark = " 🏆"
				} else {
					winnerMark = " ❌"
				}
			}
			builder.WriteString(fmt.Sprintf("• %s: %d tokens (score: %d, potential: %.0f)%s\n",
				name, p.WagerAmount, p.Score, potential, winnerMark))
		}
	}

	return builder.String()
}

// formatChallengeTimeInfo returns time-related information based on challenge status.
func (s *BotService) formatChallengeTimeInfo(challenge *models.Challenge) string {
	now := time.Now()

	switch challenge.Status {
	case constants.ChallengeStatusPending:
		deadline := challenge.CreatedAt.Add(constants.ChallengeJoinWindow)
		remaining := deadline.Sub(now)
		if remaining > 0 {
			return fmt.Sprintf("⏳ Join window: %s remaining", formatDuration(remaining))
		}
		return "⏳ Join window expired"

	case constants.ChallengeStatusActive:
		if challenge.StartDate != nil {
			deadline := challenge.StartDate.Add(constants.ChallengeDuration)
			remaining := deadline.Sub(now)
			if remaining > 0 {
				return fmt.Sprintf("⏳ Time left: %s", formatDuration(remaining))
			}
			return "⏳ Challenge period ended"
		}

	case constants.ChallengeStatusCompleted:
		return constants.StatusBadgeCompleted + " Challenge completed"

	case constants.ChallengeStatusCancelled:
		return constants.StatusBadgeCancelled + " Challenge cancelled"
	}

	return ""
}

// getTimeRemaining returns just the time remaining string for status line.
func (s *BotService) getTimeRemaining(challenge *models.Challenge) string {
	now := time.Now()

	switch challenge.Status {
	case constants.ChallengeStatusPending:
		deadline := challenge.CreatedAt.Add(constants.ChallengeJoinWindow)
		remaining := deadline.Sub(now)
		if remaining > 0 {
			return formatDuration(remaining) + " to join"
		}
		return "join window expired"

	case constants.ChallengeStatusActive:
		if challenge.StartDate != nil {
			deadline := challenge.StartDate.Add(constants.ChallengeDuration)
			remaining := deadline.Sub(now)
			if remaining > 0 {
				return formatDuration(remaining) + " left"
			}
			return "time expired"
		}
	}

	return ""
}

// formatDuration formats a duration in a human-readable way.
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// formatStatusBadge returns an emoji badge for challenge status.
func formatStatusBadge(status string) string {
	switch status {
	case constants.ChallengeStatusPending:
		return constants.StatusBadgePending
	case constants.ChallengeStatusActive:
		return constants.StatusBadgeActive
	case constants.ChallengeStatusCompleted:
		return constants.StatusBadgeCompleted
	case constants.ChallengeStatusCancelled:
		return constants.StatusBadgeCancelled
	default:
		return ""
	}
}

// formatDifficultyBadge returns an emoji badge for difficulty.
func formatDifficultyBadge(difficulty string) string {
	switch difficulty {
	case constants.DifficultyEasy:
		return constants.DifficultyBadgeEasy
	case constants.DifficultyModerate:
		return constants.DifficultyBadgeModerate
	case constants.DifficultyHard:
		return constants.DifficultyBadgeHard
	default:
		return ""
	}
}

// formatChallengeList formats a list of challenges for pagination display.
func (s *BotService) formatChallengeList(challenges []models.Challenge, page, totalPages int) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("📋 <b>Challenges</b> (Page %d/%d)\n\n", page, totalPages))

	for i, c := range challenges {
		statusBadge := formatStatusBadge(c.Status)
		difficultyBadge := formatDifficultyBadge(c.Difficulty)
		dateStr := c.CreatedAt.Format("Jan 02")
		statusText := formatStatusText(c.Status)

		builder.WriteString(fmt.Sprintf("<b>%d.</b> %s <b>%s</b>\n", i+1, difficultyBadge, c.Title))
		builder.WriteString(fmt.Sprintf("    %s %s · %s\n", statusBadge, statusText, dateStr))

		if i < len(challenges)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// formatStatusText returns a human-readable status text.
func formatStatusText(status string) string {
	switch status {
	case constants.ChallengeStatusPending:
		return "Pending"
	case constants.ChallengeStatusActive:
		return "Active"
	case constants.ChallengeStatusCompleted:
		return "Completed"
	case constants.ChallengeStatusCancelled:
		return "Cancelled"
	default:
		return status
	}
}

// formatChallengeResult formats the completion results.
func (s *BotService) formatChallengeResult(challenge *models.Challenge, winners []models.ChallengeParticipant) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("<b>%s</b> - Results\n\n", challenge.Title))

	if len(winners) == 0 {
		builder.WriteString("No winners - all wagers burned.\n")
		return builder.String()
	}

	builder.WriteString("<b>Winners:</b>\n")
	for _, w := range winners {
		name := s.getUserFirstName(w.UserID)
		payout := float64(w.WagerAmount) + float64(w.WagerAmount)*challenge.Multiplier
		profit := float64(w.WagerAmount) * challenge.Multiplier
		builder.WriteString(fmt.Sprintf("🏆 %s: +%.0f tokens (%.0f profit)\n", name, payout, profit))
	}

	return builder.String()
}

// formatChallengeCreated formats the challenge creation confirmation.
func (s *BotService) formatChallengeCreated(challenge *models.Challenge) string {
	var builder strings.Builder

	builder.WriteString("<b>Challenge Created!</b> \n\n")
	builder.WriteString(fmt.Sprintf("<b>%s</b>\n", challenge.Title))
	if challenge.Description != "" {
		builder.WriteString(fmt.Sprintf("<i>%s</i>\n", challenge.Description))
	}
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("<b>Difficulty</b>: %s (%.1fx)\n", challenge.Difficulty, challenge.Multiplier))
	builder.WriteString(fmt.Sprintf("<b>Creator wager</b>: %d tokens\n", challenge.Participants[0].WagerAmount))
	builder.WriteString("\n")
	builder.WriteString("⏳ Waiting for participants (12 hours to join)\n")
	builder.WriteString("\nUse <code>/joinchallenge [wager]</code> to join!")

	return builder.String()
}

// formatChallengeJoined formats the join confirmation.
func (s *BotService) formatChallengeJoined(challenge *models.Challenge, wager int, isActivation bool) string {
	var builder strings.Builder

	if isActivation {
		builder.WriteString("<b>Challenge Activated!</b> 🔥\n\n")
		builder.WriteString(fmt.Sprintf("<b>%s</b>\n", challenge.Title))
		builder.WriteString("14-day countdown has begun!")
	} else {
		builder.WriteString("<b>Joined Challenge!</b>\n\n")
		builder.WriteString(fmt.Sprintf("Challenge: <b>%s</b>\n", challenge.Title))
		builder.WriteString(fmt.Sprintf("Your wager: %d tokens\n", wager))
		potential := float64(wager) + float64(wager)*challenge.Multiplier
		builder.WriteString(fmt.Sprintf("Potential payout: %.0f tokens\n", potential))
	}

	return builder.String()
}

// formatScoreUpdate formats the score update confirmation.
func formatScoreUpdate(entries []ScoreEntry, names map[int64]string) string {
	var builder strings.Builder

	builder.WriteString("<b>Scores Updated!</b>\n\n")
	for _, e := range entries {
		name := names[e.UserID]
		if name == "" {
			name = "Unknown"
		}
		sign := "+"
		if e.Points < 0 {
			sign = ""
		}
		builder.WriteString(fmt.Sprintf("• %s: %s%d points\n", name, sign, e.Points))
	}

	return builder.String()
}

// formatChallengeSummary formats a brief challenge summary.
func (s *BotService) formatChallengeSummary(challenge *models.Challenge) string {
	statusBadge := formatStatusBadge(challenge.Status)
	difficultyBadge := formatDifficultyBadge(challenge.Difficulty)
	return fmt.Sprintf("%s %s <b>%s</b> (%.1fx)", difficultyBadge, statusBadge, challenge.Title, challenge.Multiplier)
}
