package store

import (
	"fmt"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// ActivityStore handles activity data operations.
type ActivityStore struct {
	db *gorm.DB
}

// NewActivityStore creates a new ActivityStore instance.
func NewActivityStore(db *gorm.DB) *ActivityStore {
	return &ActivityStore{db: db}
}

// CreateOrUpdateRecord creates a new activity record or updates existing one if message is edited.
func (s *ActivityStore) CreateOrUpdateRecord(userID, groupID int64, messageID int, activity, month, day, year string) error {
	// For new messages, try to find existing record by message_id first
	var existingRecord models.Activity
	err := s.db.Where("user_id = ? AND group_id = ? AND message_id = ?", userID, groupID, messageID).First(&existingRecord).Error
	
	if err == nil {
		// Record exists, update it
		existingRecord.Activity = activity
		existingRecord.Month = month
		existingRecord.Day = day
		existingRecord.Year = year
		return s.db.Save(&existingRecord).Error
	}
	
	// Record doesn't exist, create new one
	record := models.Activity{
		UserID:    userID,
		GroupID:   groupID,
		MessageID: messageID, // Direct int value, no pointer needed
		Activity:  activity,
		Month:     month,
		Day:       day,
		Year:      year,
	}

	return s.db.Create(&record).Error
}

// GetUserActivityByYear retrieves all activities for a user in a specific year.
func (s *ActivityStore) GetUserActivityByYear(userID, groupID int64, year string) ([]models.Activity, error) {
	var records []models.Activity
	err := s.db.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).Find(&records).Error
	return records, err
}

// GetCountByUserAndMonth returns the activity count for a user in a specific month.
func (s *ActivityStore) GetCountByUserAndMonth(userID, groupID int64, month string) (int64, error) {
	var count int64
	err := s.db.Model(&models.Activity{}).Where("user_id = ? AND group_id = ? AND month = ?", userID, groupID, month).Count(&count).Error
	return count, err
}

// GetUsersWithActivities returns all user IDs that have activities in a group.
func (s *ActivityStore) GetUsersWithActivities(groupID int64) ([]int64, error) {
	var userIDs []int64
	err := s.db.Model(&models.Activity{}).Where("group_id = ?", groupID).Distinct("user_id").Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// GetCountByUserMonthYear returns activity count for a user in a specific month and year.
func (s *ActivityStore) GetCountByUserMonthYear(userID, groupID int64, month, year string) (int64, error) {
	var count int64
	err := s.db.Model(&models.Activity{}).
		Where("user_id = ? AND group_id = ? AND month = ? AND year = ?", userID, groupID, month, year).
		Count(&count).Error
	return count, err
}

// GetMostActiveUsersForMonth returns all users tied for the highest activity count for a specific month and year.
func (s *ActivityStore) GetMostActiveUsersForMonth(groupID int64, month, year string) ([]int64, int64, error) {
	// First, get the maximum activity count
	var maxCount int64
	err := s.db.Model(&models.Activity{}).
		Select("COUNT(*) as count").
		Where("group_id = ? AND month = ? AND year = ?", groupID, month, year).
		Group("user_id").
		Order("count DESC").
		Limit(1).
		Scan(&maxCount).Error
	if err != nil {
		return nil, 0, err
	}

	if maxCount == 0 {
		return nil, 0, fmt.Errorf("no activities found for group %d in %s/%s", groupID, month, year)
	}

	// Now get all users with that maximum count
	var userIDs []int64
	err = s.db.Model(&models.Activity{}).
		Select("user_id").
		Where("group_id = ? AND month = ? AND year = ?", groupID, month, year).
		Group("user_id").
		Having("COUNT(*) = ?", maxCount).
		Pluck("user_id", &userIDs).Error
	if err != nil {
		return nil, 0, err
	}

	return userIDs, maxCount, nil
}