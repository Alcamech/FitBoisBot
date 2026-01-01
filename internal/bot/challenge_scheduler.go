package bot

import (
	"log/slog"
	"time"

	"github.com/Alcamech/FitBoisBot/internal/constants"
)

// ScheduleChallengeChecks runs periodic checks for challenge auto-actions.
func (s *BotService) ScheduleChallengeChecks() {
	ticker := time.NewTicker(constants.ChallengeSchedulerInterval)
	defer ticker.Stop()

	slog.Info("Challenge scheduler started", "interval", constants.ChallengeSchedulerInterval)

	for range ticker.C {
		s.processPendingChallenges()
		s.processActiveChallenges()
	}
}

// processPendingChallenges auto-cancels pending challenges that have exceeded the join window.
func (s *BotService) processPendingChallenges() {
	challenges, err := s.challengeStore.GetPendingChallengesForCancellation(constants.ChallengeJoinWindow)
	if err != nil {
		slog.Error("Failed to get pending challenges for cancellation", "error", err)
		return
	}

	for _, challenge := range challenges {
		// Only cancel if still only one participant (the creator)
		count, err := s.participantStore.GetParticipantCount(challenge.ID)
		if err != nil {
			slog.Error("Failed to get participant count", "error", err, "challenge_id", challenge.ID)
			continue
		}

		if count <= 1 {
			// Refund creator's wager
			creatorParticipant, err := s.participantStore.GetCreatorParticipant(challenge.ID, challenge.CreatorID)
			if err == nil && creatorParticipant != nil {
				year, _ := s.getCurrentYear(challenge.GroupID)
				if err := s.tokenStore.IncrementTokens(challenge.CreatorID, challenge.GroupID, year, creatorParticipant.WagerAmount); err != nil {
					slog.Error("Failed to refund creator", "error", err, "challenge_id", challenge.ID)
				}
			}

			// Cancel the challenge
			if err := s.challengeStore.UpdateChallengeStatus(challenge.ID, constants.ChallengeStatusCancelled); err != nil {
				slog.Error("Failed to cancel challenge", "error", err, "challenge_id", challenge.ID)
				continue
			}

			// Notify the group
			s.sendHTMLText(challenge.GroupID, constants.MsgChallengeAutoCancelled)
			slog.Info("Challenge auto-cancelled", "challenge_id", challenge.ID, "group_id", challenge.GroupID)
		}
	}
}

// processActiveChallenges auto-completes active challenges that have exceeded the duration.
func (s *BotService) processActiveChallenges() {
	challenges, err := s.challengeStore.GetActiveChallengesForCompletion(constants.ChallengeDuration)
	if err != nil {
		slog.Error("Failed to get active challenges for completion", "error", err)
		return
	}

	for _, challenge := range challenges {
		// Mark as completed (no automatic token distribution)
		if err := s.challengeStore.UpdateChallengeStatus(challenge.ID, constants.ChallengeStatusCompleted); err != nil {
			slog.Error("Failed to complete challenge", "error", err, "challenge_id", challenge.ID)
			continue
		}

		// Notify the group
		s.sendHTMLText(challenge.GroupID, constants.MsgChallengeAutoCompleted)
		slog.Info("Challenge auto-completed", "challenge_id", challenge.ID, "group_id", challenge.GroupID)
	}
}
