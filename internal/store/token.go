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

// IncrementTokens adds tokens to a user's balance for a specific year.
func (s *TokenStore) IncrementTokens(userID, groupID int64, year string, tokens int) error {
	var token models.Token

	err := s.db.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).
		First(&token).Error

	if err == gorm.ErrRecordNotFound {
		token = models.Token{
			UserID:  userID,
			GroupID: groupID,
			Year:    year,
			Balance: tokens,
		}
		return s.db.Create(&token).Error
	}

	if err != nil {
		return err
	}

	token.Balance += tokens
	return s.db.Save(&token).Error
}

// GetYearlyLeaderboard returns the token leaderboard for a specific year.
func (s *TokenStore) GetYearlyLeaderboard(groupID int64, year string) ([]models.Token, error) {
	var leaderboard []models.Token
	err := s.db.Where("group_id = ? AND year = ?", groupID, year).
		Order("balance DESC").
		Find(&leaderboard).Error
	return leaderboard, err
}

// GetUserBalance returns a user's token balance for a specific year.
func (s *TokenStore) GetUserBalance(userID, groupID int64, year string) (int, error) {
	var token models.Token
	err := s.db.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).
		First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return token.Balance, nil
}

// GetUserLifetimeBalance returns a user's total token balance across all years.
func (s *TokenStore) GetUserLifetimeBalance(userID, groupID int64) (int, error) {
	var total int
	err := s.db.Model(&models.Token{}).
		Where("user_id = ? AND group_id = ?", userID, groupID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}