package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Alcamech/FitBoisBot/internal/constants"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// onChallenge handles the /challenge command to create a new challenge.
func (s *BotService) onChallenge(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	args := msg.CommandArguments()
	if args == "" {
		s.sendHTMLText(chatID, constants.MsgChallengeInvalidFormat)
		return
	}

	// Parse command arguments
	params, err := parseChallengeCommand(args)
	if err != nil {
		if strings.Contains(err.Error(), "difficulty") {
			s.sendHTMLText(chatID, constants.MsgChallengeInvalidDifficulty)
		} else if strings.Contains(err.Error(), "wager") {
			s.sendHTMLText(chatID, constants.MsgChallengeInvalidWager)
		} else {
			s.sendHTMLText(chatID, constants.MsgChallengeInvalidFormat)
		}
		return
	}

	// Check for existing active/pending challenge
	existing, err := s.challengeStore.GetActiveOrPendingChallenge(chatID)
	if err == nil && existing != nil {
		s.sendHTMLText(chatID, constants.MsgChallengeAlreadyExists)
		return
	}

	// Check user has sufficient tokens (lifetime balance)
	hasTokens, err := s.checkUserTokens(userID, chatID, params.Wager)
	if err != nil || !hasTokens {
		s.sendHTMLText(chatID, constants.MsgChallengeInsufficientTokens)
		return
	}

	// Deduct wager from creator's balance
	if err := s.userBalanceStore.IncrementBalance(userID, chatID, -params.Wager); err != nil {
		slog.Error("Failed to deduct tokens", "error", err)
		s.sendHTMLText(chatID, "Failed to create challenge. Please try again.")
		return
	}

	// Create challenge
	challenge := &models.Challenge{
		GroupID:     chatID,
		CreatorID:   userID,
		Title:       params.Title,
		Description: params.Description,
		Difficulty:  params.Difficulty,
		Multiplier:  params.Multiplier,
		Status:      constants.ChallengeStatusPending,
	}

	if err := s.challengeStore.CreateChallenge(challenge); err != nil {
		slog.Error("Failed to create challenge", "error", err)
		// Refund tokens
		s.userBalanceStore.IncrementBalance(userID, chatID, params.Wager)
		s.sendHTMLText(chatID, "Failed to create challenge. Please try again.")
		return
	}

	// Add creator as first participant
	participant := &models.ChallengeParticipant{
		ChallengeID: challenge.ID,
		UserID:      userID,
		WagerAmount: params.Wager,
	}

	if err := s.participantStore.CreateParticipant(participant); err != nil {
		slog.Error("Failed to add creator as participant", "error", err)
		return
	}

	// Reload challenge with participants for formatting
	challenge, _ = s.challengeStore.GetChallengeByID(challenge.ID)
	message := s.formatChallengeCreated(challenge)
	s.sendHTMLText(chatID, message)
}

// onJoinChallenge handles the /joinchallenge command.
func (s *BotService) onJoinChallenge(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	args := msg.CommandArguments()
	wager, err := parseJoinChallengeCommand(args)
	if err != nil {
		s.sendHTMLText(chatID, constants.MsgChallengeInvalidWager)
		return
	}

	// Get active or pending challenge
	challenge, err := s.challengeStore.GetActiveOrPendingChallenge(chatID)
	if err != nil {
		s.sendHTMLText(chatID, constants.MsgChallengeNotFound)
		return
	}

	// Check if user is already a participant
	isParticipant, err := s.participantStore.IsUserParticipant(challenge.ID, userID)
	if err != nil {
		slog.Error("Failed to check participant", "error", err)
		return
	}
	if isParticipant {
		s.sendHTMLText(chatID, constants.MsgChallengeAlreadyParticipant)
		return
	}

	// Check user has sufficient tokens (lifetime balance)
	hasTokens, err := s.checkUserTokens(userID, chatID, wager)
	if err != nil || !hasTokens {
		s.sendHTMLText(chatID, constants.MsgChallengeInsufficientTokens)
		return
	}

	// Deduct wager from user's balance
	if err := s.userBalanceStore.IncrementBalance(userID, chatID, -wager); err != nil {
		slog.Error("Failed to deduct tokens", "error", err)
		return
	}

	// Add participant
	participant := &models.ChallengeParticipant{
		ChallengeID: challenge.ID,
		UserID:      userID,
		WagerAmount: wager,
	}

	if err := s.participantStore.CreateParticipant(participant); err != nil {
		slog.Error("Failed to add participant", "error", err)
		// Refund tokens
		s.userBalanceStore.IncrementBalance(userID, chatID, wager)
		return
	}

	// Check if this activates the challenge (first non-creator participant)
	isActivation := challenge.Status == constants.ChallengeStatusPending
	if isActivation {
		if err := s.challengeStore.ActivateChallenge(challenge.ID); err != nil {
			slog.Error("Failed to activate challenge", "error", err)
		}
		challenge.Status = constants.ChallengeStatusActive
	}

	message := s.formatChallengeJoined(challenge, wager, isActivation)
	s.sendHTMLText(chatID, message)
}

// onViewChallenge handles the /viewchallenge command.
func (s *BotService) onViewChallenge(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	args := msg.CommandArguments()
	id, hasID := parseViewChallengeCommand(args)

	var challenge *models.Challenge
	var err error

	if hasID {
		challenge, err = s.challengeStore.GetChallengeByID(id)
	} else {
		challenge, err = s.challengeStore.GetCurrentChallenge(chatID)
	}

	if err != nil {
		s.sendHTMLText(chatID, constants.MsgChallengeNotFound)
		return
	}

	message := s.formatChallengeDetails(challenge)
	s.sendHTMLText(chatID, message)
}

// onListChallenges handles the /listchallenges command.
func (s *BotService) onListChallenges(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	args := msg.CommandArguments()
	page := parseListChallengesCommand(args)

	s.sendChallengeListPage(chatID, page, 0)
}

// sendChallengeListPage sends a paginated challenge list.
func (s *BotService) sendChallengeListPage(chatID int64, page int, editMessageID int) {
	limit := constants.ChallengePageSize
	offset := (page - 1) * limit

	total, err := s.challengeStore.CountChallengesByGroup(chatID)
	if err != nil {
		slog.Error("Failed to count challenges", "error", err)
		return
	}

	if total == 0 {
		s.sendHTMLText(chatID, constants.MsgChallengeNoChallenges)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	challenges, err := s.challengeStore.GetChallengesByGroup(chatID, limit, offset)
	if err != nil {
		slog.Error("Failed to get challenges", "error", err)
		return
	}

	message := s.formatChallengeList(challenges, page, totalPages)
	keyboard := createChallengeListKeyboard(challenges, page, totalPages)

	if editMessageID > 0 {
		// Edit existing message
		editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, editMessageID, message, keyboard)
		editMsg.ParseMode = "HTML"
		s.bot.Send(editMsg)
	} else {
		// Send new message
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = keyboard
		s.bot.Send(msg)
	}
}

// onScore handles the /score command.
func (s *BotService) onScore(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	// Get active challenge
	challenge, err := s.challengeStore.GetActiveOrPendingChallenge(chatID)
	if err != nil {
		s.sendHTMLText(chatID, constants.MsgChallengeNotFound)
		return
	}

	// Only active challenges can be scored
	if challenge.Status != constants.ChallengeStatusActive {
		s.sendHTMLText(chatID, constants.MsgChallengeNoLongerActive)
		return
	}

	// Parse score entries
	entries, err := parseScoreCommand(msg)
	if err != nil {
		s.sendHTMLText(chatID, fmt.Sprintf("Invalid score format: %s", err.Error()))
		return
	}

	// The creator can score anyone; other participants can only score themselves.
	if !canScoreEntries(challenge.CreatorID == userID, userID, entries) {
		s.sendHTMLText(chatID, constants.MsgChallengeOnlyScoreSelf)
		return
	}

	// Build scores map and verify all users are participants
	scores := make(map[int64]int)
	names := make(map[int64]string)
	for _, entry := range entries {
		isParticipant, err := s.participantStore.IsUserParticipant(challenge.ID, entry.UserID)
		if err != nil || !isParticipant {
			s.sendHTMLText(chatID, constants.MsgChallengeUserNotParticipant)
			return
		}
		scores[entry.UserID] = entry.Points
		names[entry.UserID] = s.getUserFirstName(entry.UserID)
	}

	// Update scores
	if err := s.participantStore.UpdateScores(challenge.ID, scores); err != nil {
		slog.Error("Failed to update scores", "error", err)
		return
	}

	message := formatScoreUpdate(entries, names)
	s.sendHTMLText(chatID, message)
}

// onCancelChallenge handles the /cancelchallenge command.
func (s *BotService) onCancelChallenge(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	challenge, err := s.challengeStore.GetActiveOrPendingChallenge(chatID)
	if err != nil {
		s.sendHTMLText(chatID, constants.MsgChallengeNotFound)
		return
	}

	// Only creator can cancel
	if challenge.CreatorID != userID {
		s.sendHTMLText(chatID, constants.MsgChallengeOnlyCreator)
		return
	}

	// Can only cancel pending challenges
	if challenge.Status != constants.ChallengeStatusPending {
		s.sendHTMLText(chatID, constants.MsgChallengeCannotCancel)
		return
	}

	// Refund creator's wager to their balance
	creatorParticipant, err := s.participantStore.GetCreatorParticipant(challenge.ID, userID)
	if err == nil && creatorParticipant != nil {
		s.userBalanceStore.IncrementBalance(userID, chatID, creatorParticipant.WagerAmount)
	}

	// Update status
	if err := s.challengeStore.UpdateChallengeStatus(challenge.ID, constants.ChallengeStatusCancelled); err != nil {
		slog.Error("Failed to cancel challenge", "error", err)
		return
	}

	s.sendHTMLText(chatID, constants.MsgChallengeCancelled)
}

// onCompleteChallenge handles the /completechallenge command.
func (s *BotService) onCompleteChallenge(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	// Get challenge (can be active or completed but not yet awarded)
	challenge, err := s.challengeStore.GetActiveOrPendingChallenge(chatID)
	if err != nil {
		// Try to get a completed challenge that might need awarding
		challenge, err = s.challengeStore.GetCurrentChallenge(chatID)
		if err != nil || (challenge.Status != constants.ChallengeStatusActive && challenge.Status != constants.ChallengeStatusCompleted) {
			s.sendHTMLText(chatID, constants.MsgChallengeNotFound)
			return
		}
	}

	// Only creator can complete
	if challenge.CreatorID != userID {
		s.sendHTMLText(chatID, constants.MsgChallengeOnlyCreator)
		return
	}

	// Check if already awarded (any participant has is_winner set)
	for _, p := range challenge.Participants {
		if p.IsWinner != nil {
			s.sendHTMLText(chatID, constants.MsgChallengeAlreadyAwarded)
			return
		}
	}

	// Parse winners from message
	winnerIDs := parseCompleteChallengeCommand(msg)

	year, _ := s.getCurrentYear(chatID)

	if len(winnerIDs) == 0 {
		// No winners - all wagers burned
		if err := s.participantStore.SetWinners(challenge.ID, []int64{}); err != nil {
			slog.Error("Failed to set no winners", "error", err)
		}
		s.challengeStore.UpdateChallengeStatus(challenge.ID, constants.ChallengeStatusCompleted)
		s.sendHTMLText(chatID, constants.MsgChallengeNoWinners)
		return
	}

	// Verify all winners are participants
	for _, winnerID := range winnerIDs {
		isParticipant, err := s.participantStore.IsUserParticipant(challenge.ID, winnerID)
		if err != nil || !isParticipant {
			s.sendHTMLText(chatID, constants.MsgChallengeUserNotParticipant)
			return
		}
	}

	// Get winner participants for payout calculation
	winners, err := s.participantStore.GetParticipantsByUserIDs(challenge.ID, winnerIDs)
	if err != nil {
		slog.Error("Failed to get winners", "error", err)
		return
	}

	// Award tokens to winners
	for _, winner := range winners {
		// Payout = wager + (wager * multiplier)
		// Bonus (wager * multiplier) is tracked as earnings
		bonus := int(float64(winner.WagerAmount) * challenge.Multiplier)
		payout := winner.WagerAmount + bonus

		// Track the bonus as yearly earnings
		if err := s.tokenStore.AddEarnings(winner.UserID, chatID, year, bonus); err != nil {
			slog.Error("Failed to record earnings", "error", err, "user_id", winner.UserID)
		}

		// Add full payout to spendable balance
		if err := s.userBalanceStore.IncrementBalance(winner.UserID, chatID, payout); err != nil {
			slog.Error("Failed to add to balance", "error", err, "user_id", winner.UserID)
		}
	}

	// Mark winners
	if err := s.participantStore.SetWinners(challenge.ID, winnerIDs); err != nil {
		slog.Error("Failed to set winners", "error", err)
	}

	// Update status
	s.challengeStore.UpdateChallengeStatus(challenge.ID, constants.ChallengeStatusCompleted)

	message := s.formatChallengeResult(challenge, winners)
	s.sendHTMLText(chatID, message)
}

// handleChallengeCallback handles inline keyboard callbacks for challenges.
func (s *BotService) handleChallengeCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	action, value := parseCallbackData(callback.Data)

	switch action {
	case "challenge_list":
		page, _ := strconv.Atoi(value)
		if page < 1 {
			page = 1
		}
		s.sendChallengeListPage(chatID, page, messageID)

	case "challenge_view":
		id, _ := strconv.ParseInt(value, 10, 64)
		challenge, err := s.challengeStore.GetChallengeByID(id)
		if err != nil {
			return
		}
		message := s.formatChallengeDetails(challenge)
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, message)
		editMsg.ParseMode = "HTML"
		// Add back button
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(constants.ButtonBackToList, "challenge_list_1"),
			),
		)
		editMsg.ReplyMarkup = &keyboard
		s.bot.Send(editMsg)
	}

	// Answer callback to remove loading state
	s.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
}

// Helper functions

// getCurrentYear gets the current year for the group's timezone.
func (s *BotService) getCurrentYear(chatID int64) (string, error) {
	group, err := s.groupStore.GetByID(chatID)
	if err != nil {
		return "", err
	}

	timezone := group.Timezone
	if timezone == "" {
		timezone = constants.DefaultTimezone
	}

	return GetCurrentYear(timezone)
}

// checkUserTokens checks if a user has sufficient tokens in their spendable balance.
func (s *BotService) checkUserTokens(userID, groupID int64, amount int) (bool, error) {
	return s.userBalanceStore.HasSufficientBalance(userID, groupID, amount)
}

// canScoreEntries reports whether the sender is authorized to apply the given
// score entries. The challenge creator may score any participant; every other
// participant may only score themselves.
func canScoreEntries(isCreator bool, senderID int64, entries []ScoreEntry) bool {
	if isCreator {
		return true
	}
	for _, entry := range entries {
		if entry.UserID != senderID {
			return false
		}
	}
	return true
}
