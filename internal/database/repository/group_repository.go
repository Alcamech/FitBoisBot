package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	DB *gorm.DB
}

func (r *GroupRepository) SetGroupTimezone(groupID int64, timezone string) error {
	return r.DB.Model(&models.Group{}).Where("id = ?", groupID).Update("timezone", timezone).Error
}

func (r *GroupRepository) GetGroupByID(groupID int64) (*models.Group, error) {
	var group models.Group
	err := r.DB.First(&group, "id = ?", groupID).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) GetAllGroups() ([]models.Group, error) {
	var groups []models.Group
	err := r.DB.Find(&groups).Error
	return groups, err
}
