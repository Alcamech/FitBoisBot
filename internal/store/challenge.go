package store

import (
	"time"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// ChallengeStore handles challenge data operations.
type ChallengeStore struct {
	db *gorm.DB
}

// NewChallengeStore creates a new ChallengeStore instance.
func NewChallengeStore(db *gorm.DB) *ChallengeStore {
	return &ChallengeStore{db: db}
}

// CreateChallenge creates a new challenge.
func (s *ChallengeStore) CreateChallenge(challenge *models.Challenge) error {
	return s.db.Create(challenge).Error
}

// GetCurrentChallenge returns the most recent challenge for a group (any status).
func (s *ChallengeStore) GetCurrentChallenge(groupID int64) (*models.Challenge, error) {
	var challenge models.Challenge
	err := s.db.Preload("Participants").Preload("Participants.User").Preload("Creator").
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		First(&challenge).Error
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

// GetActiveOrPendingChallenge returns an active or pending challenge for a group if one exists.
func (s *ChallengeStore) GetActiveOrPendingChallenge(groupID int64) (*models.Challenge, error) {
	var challenge models.Challenge
	err := s.db.Preload("Participants").Preload("Participants.User").Preload("Creator").
		Where("group_id = ? AND status IN ?", groupID, []string{"active", "pending"}).
		First(&challenge).Error
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

// GetChallengeByID retrieves a specific challenge by ID with participants preloaded.
func (s *ChallengeStore) GetChallengeByID(id int64) (*models.Challenge, error) {
	var challenge models.Challenge
	err := s.db.Preload("Participants").Preload("Participants.User").Preload("Creator").
		First(&challenge, id).Error
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

// GetChallengesByGroup returns paginated challenges for a group sorted by created_at DESC.
func (s *ChallengeStore) GetChallengesByGroup(groupID int64, limit, offset int) ([]models.Challenge, error) {
	var challenges []models.Challenge
	err := s.db.Where("group_id = ?", groupID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&challenges).Error
	return challenges, err
}

// CountChallengesByGroup returns the total count of challenges for a group.
func (s *ChallengeStore) CountChallengesByGroup(groupID int64) (int64, error) {
	var count int64
	err := s.db.Model(&models.Challenge{}).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return count, err
}

// UpdateChallengeStatus updates the status of a challenge.
func (s *ChallengeStore) UpdateChallengeStatus(id int64, status string) error {
	return s.db.Model(&models.Challenge{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ActivateChallenge sets the challenge to active and records the start date.
func (s *ChallengeStore) ActivateChallenge(id int64) error {
	now := time.Now()
	return s.db.Model(&models.Challenge{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "active",
			"start_date": now,
		}).Error
}

// GetPendingChallengesForCancellation returns pending challenges older than the given duration.
func (s *ChallengeStore) GetPendingChallengesForCancellation(joinWindow time.Duration) ([]models.Challenge, error) {
	var challenges []models.Challenge
	cutoff := time.Now().Add(-joinWindow)
	err := s.db.Preload("Participants").Preload("Creator").
		Where("status = ? AND created_at < ?", "pending", cutoff).
		Find(&challenges).Error
	return challenges, err
}

// GetActiveChallengesForCompletion returns active challenges that have exceeded the challenge duration.
func (s *ChallengeStore) GetActiveChallengesForCompletion(duration time.Duration) ([]models.Challenge, error) {
	var challenges []models.Challenge
	cutoff := time.Now().Add(-duration)
	err := s.db.Preload("Participants").Preload("Creator").
		Where("status = ? AND start_date < ?", "active", cutoff).
		Find(&challenges).Error
	return challenges, err
}
