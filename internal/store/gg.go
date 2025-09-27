package store

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// GGStore handles "fastest GG" data operations.
type GGStore struct {
	db *gorm.DB
}

// NewGGStore creates a new GGStore instance.
func NewGGStore(db *gorm.DB) *GGStore {
	return &GGStore{db: db}
}

// CreateOrUpdateCount creates a new GG record or increments existing count.
func (s *GGStore) CreateOrUpdateCount(userID, groupID int64, year string) error {
	var gg models.Gg
	err := s.db.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).First(&gg).Error
	if err == gorm.ErrRecordNotFound {
		gg = models.Gg{
			UserID:      userID,
			GroupID:     groupID,
			Year:        year,
			FastGGCount: 1,
		}
		return s.db.Create(&gg).Error
	}

	if err != nil {
		return err
	}

	gg.FastGGCount++
	return s.db.Save(&gg).Error
}

// GetLeaderboard returns the GG leaderboard for a group and year.
func (s *GGStore) GetLeaderboard(groupID int64, year string) ([]models.Gg, error) {
	var leaderboard []models.Gg
	err := s.db.Where("group_id = ? AND year = ?", groupID, year).
		Order("fast_gg_count DESC").
		Find(&leaderboard).Error
	return leaderboard, err
}