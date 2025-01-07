package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type TokenRepository struct {
	DB *gorm.DB
}

func (r *TokenRepository) IncrementTokens(userID, groupID int64, year string, tokens int) error {
	var token models.Token

	err := r.DB.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).
		First(&token).Error

	if err == gorm.ErrRecordNotFound {
		token = models.Token{
			UserID:  userID,
			GroupID: groupID,
			Year:    year,
			Balance: tokens,
		}
		return r.DB.Create(&token).Error
	}

	if err != nil {
		return err
	}

	token.Balance += tokens
	return r.DB.Save(&token).Error
}

func (r *TokenRepository) GetYearlyLeaderboard(groupID int64, year string) ([]models.Token, error) {
	var leaderboard []models.Token
	err := r.DB.Where("group_id = ? AND year = ?", groupID, year).
		Order("balance DESC").
		Find(&leaderboard).Error
	return leaderboard, err
}
