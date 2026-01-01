package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ChallengeParams holds parsed parameters for challenge creation.
type ChallengeParams struct {
	Difficulty  string
	Multiplier  float64
	Wager       int
	Title       string
	Description string
}

// parseChallengeCommand parses the /challenge command arguments.
// Format: /challenge [difficulty] [wager] [title] [description]
func parseChallengeCommand(args string) (*ChallengeParams, error) {
	parts := strings.Fields(args)
	if len(parts) < 3 {
		return nil, fmt.Errorf("not enough arguments")
	}

	difficulty := strings.ToLower(parts[0])
	multiplier, err := getDifficultyMultiplier(difficulty)
	if err != nil {
		return nil, err
	}

	wager, err := strconv.Atoi(parts[1])
	if err != nil || wager <= 0 {
		return nil, fmt.Errorf("invalid wager amount")
	}

	// Title is the third part
	title := parts[2]

	// Description is everything after title (optional)
	var description string
	if len(parts) > 3 {
		description = strings.Join(parts[3:], " ")
	}

	return &ChallengeParams{
		Difficulty:  difficulty,
		Multiplier:  multiplier,
		Wager:       wager,
		Title:       title,
		Description: description,
	}, nil
}

// getDifficultyMultiplier returns the multiplier for a difficulty level.
func getDifficultyMultiplier(difficulty string) (float64, error) {
	switch difficulty {
	case constants.DifficultyEasy:
		return constants.MultiplierEasy, nil
	case constants.DifficultyModerate:
		return constants.MultiplierModerate, nil
	case constants.DifficultyHard:
		return constants.MultiplierHard, nil
	default:
		return 0, fmt.Errorf("invalid difficulty")
	}
}

// parseJoinChallengeCommand parses the /joinchallenge command arguments.
// Format: /joinchallenge [wager]
func parseJoinChallengeCommand(args string) (int, error) {
	args = strings.TrimSpace(args)
	if args == "" {
		return 0, fmt.Errorf("wager amount required")
	}

	wager, err := strconv.Atoi(args)
	if err != nil || wager <= 0 {
		return 0, fmt.Errorf("invalid wager amount")
	}

	return wager, nil
}

// parseViewChallengeCommand parses the /viewchallenge command arguments.
// Format: /viewchallenge [id] (optional)
func parseViewChallengeCommand(args string) (int64, bool) {
	args = strings.TrimSpace(args)
	if args == "" {
		return 0, false // No ID provided, show current challenge
	}

	id, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return 0, false
	}

	return id, true
}

// parseListChallengesCommand parses the /listchallenges command arguments.
// Format: /listchallenges [page] (optional, defaults to 1)
func parseListChallengesCommand(args string) int {
	args = strings.TrimSpace(args)
	if args == "" {
		return 1
	}

	page, err := strconv.Atoi(args)
	if err != nil || page < 1 {
		return 1
	}

	return page
}

// ScoreEntry represents a user's score update.
type ScoreEntry struct {
	UserID int64
	Points int
}

// parseScoreCommand parses the /score command from a message.
// Supports two formats:
// - Individual scores: /score @john 5 @rob 6
// - Shared score: /score @john @rob @mark 10
func parseScoreCommand(msg *tgbotapi.Message) ([]ScoreEntry, error) {
	if len(msg.Entities) == 0 {
		return nil, fmt.Errorf("no users mentioned")
	}

	args := msg.CommandArguments()
	if args == "" {
		return nil, fmt.Errorf("no arguments provided")
	}

	// Extract all mentioned users
	mentionedUsers := extractMentionedUserIDs(msg)
	if len(mentionedUsers) == 0 {
		return nil, fmt.Errorf("no valid users mentioned")
	}

	// Parse the arguments to extract numbers
	parts := strings.Fields(args)
	numbers := []int{}
	for _, part := range parts {
		if !strings.HasPrefix(part, "@") {
			num, err := strconv.Atoi(part)
			if err == nil {
				numbers = append(numbers, num)
			}
		}
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("no point values provided")
	}

	var entries []ScoreEntry

	// If we have one number and multiple users, it's a shared score
	if len(numbers) == 1 {
		for _, userID := range mentionedUsers {
			entries = append(entries, ScoreEntry{
				UserID: userID,
				Points: numbers[0],
			})
		}
		return entries, nil
	}

	// If we have equal numbers to users, it's individual scores
	if len(numbers) == len(mentionedUsers) {
		for i, userID := range mentionedUsers {
			entries = append(entries, ScoreEntry{
				UserID: userID,
				Points: numbers[i],
			})
		}
		return entries, nil
	}

	// Try to match by position in the argument string
	// Parse in order: @user1 num1 @user2 num2 ...
	entries = parsePositionalScores(args, msg)
	if len(entries) > 0 {
		return entries, nil
	}

	return nil, fmt.Errorf("could not parse score format")
}

// parsePositionalScores parses scores in positional format: @user1 num1 @user2 num2
func parsePositionalScores(args string, msg *tgbotapi.Message) []ScoreEntry {
	var entries []ScoreEntry
	parts := strings.Fields(args)

	var currentUserID int64
	for _, part := range parts {
		if strings.HasPrefix(part, "@") {
			// Find user ID for this mention
			username := strings.TrimPrefix(part, "@")
			userID := findUserIDByUsername(msg, username)
			if userID != 0 {
				currentUserID = userID
			}
		} else if currentUserID != 0 {
			// Try to parse as number
			points, err := strconv.Atoi(part)
			if err == nil {
				entries = append(entries, ScoreEntry{
					UserID: currentUserID,
					Points: points,
				})
				currentUserID = 0
			}
		}
	}

	return entries
}

// findUserIDByUsername finds a user ID from message entities by username.
func findUserIDByUsername(msg *tgbotapi.Message, username string) int64 {
	if len(msg.Entities) == 0 {
		return 0
	}

	for _, entity := range msg.Entities {
		if entity.Type == "text_mention" && entity.User != nil {
			// Check if the mention text matches
			mentionText := msg.Text[entity.Offset : entity.Offset+entity.Length]
			if strings.TrimPrefix(mentionText, "@") == username {
				return entity.User.ID
			}
		} else if entity.Type == "mention" {
			// For @username mentions, we need to look up by username
			mentionText := msg.Text[entity.Offset+1 : entity.Offset+entity.Length] // Skip @
			if mentionText == username && entity.User != nil {
				return entity.User.ID
			}
		}
	}

	return 0
}

// extractMentionedUserIDs extracts user IDs from message mentions.
func extractMentionedUserIDs(msg *tgbotapi.Message) []int64 {
	var userIDs []int64
	seen := make(map[int64]bool)

	if len(msg.Entities) == 0 {
		return userIDs
	}

	for _, entity := range msg.Entities {
		if entity.Type == "text_mention" && entity.User != nil {
			if !seen[entity.User.ID] {
				userIDs = append(userIDs, entity.User.ID)
				seen[entity.User.ID] = true
			}
		}
	}

	return userIDs
}

// parseCompleteChallengeCommand extracts winner user IDs from the message.
func parseCompleteChallengeCommand(msg *tgbotapi.Message) []int64 {
	return extractMentionedUserIDs(msg)
}

// parseCallbackData parses inline keyboard callback data.
// Returns action type and value.
func parseCallbackData(data string) (action string, value string) {
	parts := strings.SplitN(data, "_", 3)
	if len(parts) < 2 {
		return "", ""
	}

	if len(parts) == 2 {
		return parts[0] + "_" + parts[1], ""
	}

	return parts[0] + "_" + parts[1], parts[2]
}
