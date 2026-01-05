package store

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// TokenStore handles token data operations.
type TokenStore struct {
	db *gorm.DB
}

// NewTokenStore creates a new TokenStore instance.
func NewTokenStore(db *gorm.DB) *TokenStore {
	return &TokenStore{db: db}
}

// AddEarnings adds tokens to a user's earnings for a specific year.
// This tracks how many tokens a user has earned (always positive amounts).
func (s *TokenStore) AddEarnings(userID, groupID int64, year string, tokens int) error {
	var token models.Token

	err := s.db.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).
		First(&token).Error

	if err == gorm.ErrRecordNotFound {
		token = models.Token{
			UserID:  userID,
			GroupID: groupID,
			Year:    year,
			Earned:  tokens,
		}
		return s.db.Create(&token).Error
	}

	if err != nil {
		return err
	}

	token.Earned += tokens
	return s.db.Save(&token).Error
}

// GetYearlyLeaderboard returns the token earnings leaderboard for a specific year.
func (s *TokenStore) GetYearlyLeaderboard(groupID int64, year string) ([]models.Token, error) {
	var leaderboard []models.Token
	err := s.db.Where("group_id = ? AND year = ?", groupID, year).
		Order("earned DESC").
		Find(&leaderboard).Error
	return leaderboard, err
}

// GetUserEarnings returns a user's token earnings for a specific year.
func (s *TokenStore) GetUserEarnings(userID, groupID int64, year string) (int, error) {
	var token models.Token
	err := s.db.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).
		First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return token.Earned, nil
}