package store

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// GroupStore handles group data operations.
type GroupStore struct {
	db *gorm.DB
}

// NewGroupStore creates a new GroupStore instance.
func NewGroupStore(db *gorm.DB) *GroupStore {
	return &GroupStore{db: db}
}

// SetTimezone updates the timezone for a group.
func (s *GroupStore) SetTimezone(groupID int64, timezone string) error {
	return s.db.Model(&models.Group{}).Where("id = ?", groupID).Update("timezone", timezone).Error
}

// GetByID retrieves a group by ID.
func (s *GroupStore) GetByID(groupID int64) (*models.Group, error) {
	var group models.Group
	err := s.db.First(&group, "id = ?", groupID).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetAll retrieves all groups.
func (s *GroupStore) GetAll() ([]models.Group, error) {
	var groups []models.Group
	err := s.db.Find(&groups).Error
	return groups, err
}