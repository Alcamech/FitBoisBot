package store

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// UserBalanceStore handles user balance data operations.
type UserBalanceStore struct {
	db *gorm.DB
}

// NewUserBalanceStore creates a new UserBalanceStore instance.
func NewUserBalanceStore(db *gorm.DB) *UserBalanceStore {
	return &UserBalanceStore{db: db}
}

// GetBalance returns a user's spendable token balance.
func (s *UserBalanceStore) GetBalance(userID, groupID int64) (int, error) {
	var balance models.UserBalance
	err := s.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&balance).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return balance.Balance, nil
}

// IncrementBalance adds (or subtracts if negative) tokens from a user's balance.
func (s *UserBalanceStore) IncrementBalance(userID, groupID int64, amount int) error {
	var balance models.UserBalance

	err := s.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&balance).Error

	if err == gorm.ErrRecordNotFound {
		balance = models.UserBalance{
			UserID:  userID,
			GroupID: groupID,
			Balance: amount,
		}
		return s.db.Create(&balance).Error
	}

	if err != nil {
		return err
	}

	balance.Balance += amount
	return s.db.Save(&balance).Error
}

// HasSufficientBalance checks if a user has at least the specified amount of tokens.
func (s *UserBalanceStore) HasSufficientBalance(userID, groupID int64, amount int) (bool, error) {
	balance, err := s.GetBalance(userID, groupID)
	if err != nil {
		return false, err
	}
	return balance >= amount, nil
}
