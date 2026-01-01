package store

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// ParticipantStore handles challenge participant data operations.
type ParticipantStore struct {
	db *gorm.DB
}

// NewParticipantStore creates a new ParticipantStore instance.
func NewParticipantStore(db *gorm.DB) *ParticipantStore {
	return &ParticipantStore{db: db}
}

// CreateParticipant adds a participant to a challenge.
func (s *ParticipantStore) CreateParticipant(participant *models.ChallengeParticipant) error {
	return s.db.Create(participant).Error
}

// GetParticipants returns all participants for a challenge.
func (s *ParticipantStore) GetParticipants(challengeID int64) ([]models.ChallengeParticipant, error) {
	var participants []models.ChallengeParticipant
	err := s.db.Preload("User").
		Where("challenge_id = ?", challengeID).
		Order("joined_at ASC").
		Find(&participants).Error
	return participants, err
}

// GetParticipantCount returns the number of participants in a challenge.
func (s *ParticipantStore) GetParticipantCount(challengeID int64) (int64, error) {
	var count int64
	err := s.db.Model(&models.ChallengeParticipant{}).
		Where("challenge_id = ?", challengeID).
		Count(&count).Error
	return count, err
}

// UpdateScores updates scores for multiple participants. The scores map contains userID -> points delta.
func (s *ParticipantStore) UpdateScores(challengeID int64, scores map[int64]int) error {
	for userID, points := range scores {
		err := s.db.Model(&models.ChallengeParticipant{}).
			Where("challenge_id = ? AND user_id = ?", challengeID, userID).
			Update("score", gorm.Expr("score + ?", points)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// SetWinners marks the specified users as winners and others as losers.
func (s *ParticipantStore) SetWinners(challengeID int64, winnerUserIDs []int64) error {
	// Mark all participants as losers first
	falseVal := false
	err := s.db.Model(&models.ChallengeParticipant{}).
		Where("challenge_id = ?", challengeID).
		Update("is_winner", &falseVal).Error
	if err != nil {
		return err
	}

	// Mark winners
	if len(winnerUserIDs) > 0 {
		trueVal := true
		err = s.db.Model(&models.ChallengeParticipant{}).
			Where("challenge_id = ? AND user_id IN ?", challengeID, winnerUserIDs).
			Update("is_winner", &trueVal).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// IsUserParticipant checks if a user is already a participant in a challenge.
func (s *ParticipantStore) IsUserParticipant(challengeID, userID int64) (bool, error) {
	var count int64
	err := s.db.Model(&models.ChallengeParticipant{}).
		Where("challenge_id = ? AND user_id = ?", challengeID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetParticipantsByUserIDs returns participants matching the given user IDs for a challenge.
func (s *ParticipantStore) GetParticipantsByUserIDs(challengeID int64, userIDs []int64) ([]models.ChallengeParticipant, error) {
	var participants []models.ChallengeParticipant
	err := s.db.Preload("User").
		Where("challenge_id = ? AND user_id IN ?", challengeID, userIDs).
		Find(&participants).Error
	return participants, err
}

// GetCreatorParticipant returns the participant record for the challenge creator.
func (s *ParticipantStore) GetCreatorParticipant(challengeID, creatorID int64) (*models.ChallengeParticipant, error) {
	var participant models.ChallengeParticipant
	err := s.db.Where("challenge_id = ? AND user_id = ?", challengeID, creatorID).
		First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}
