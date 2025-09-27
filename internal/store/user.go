// Package store provides data access functionality for the FitBoisBot application.
package store

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"github.com/Alcamech/FitBoisBot/internal/errors"
	"gorm.io/gorm"
)

// UserStore handles user data operations.
type UserStore struct {
	db *gorm.DB
}

// NewUserStore creates a new UserStore instance.
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

// GetOrCreateUser retrieves an existing user or creates a new one.
func (s *UserStore) GetOrCreateUser(userID int64, name string, groupID int64) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	
	// User exists, return it
	if err == nil {
		return &user, nil
	}
	
	// Handle other errors
	if err != gorm.ErrRecordNotFound {
		return nil, errors.NewDatabaseError("find", "users", err)
	}
	
	// User doesn't exist, ensure group exists first
	if err := s.ensureGroupExists(groupID); err != nil {
		return nil, err
	}
	
	// Create new user
	user = models.User{
		ID:      userID,
		Name:    name,
		GroupID: groupID,
	}
	
	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.NewDatabaseError("create", "users", err)
	}
	
	return &user, nil
}

// ensureGroupExists creates a group if it doesn't exist
func (s *UserStore) ensureGroupExists(groupID int64) error {
	var group models.Group
	err := s.db.First(&group, "id = ?", groupID).Error
	
	if err == nil {
		return nil // Group exists
	}
	
	if err != gorm.ErrRecordNotFound {
		return errors.NewDatabaseError("find", "groups", err)
	}
	
	// Group doesn't exist, create it
	group = models.Group{ID: groupID}
	if err := s.db.Create(&group).Error; err != nil {
		return errors.NewDatabaseError("create", "groups", err)
	}
	return nil
}

// FindByID retrieves a user by their ID.
func (s *UserStore) FindByID(userID int64) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}