package store

import (
	"fmt"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

// ActivityPost represents activity data for creating or updating records.
type ActivityPost struct {
	UserID    int64
	GroupID   int64
	MessageID int
	Activity  string
	Month     string
	Day       string
	Year      string
}

// ActivityStore handles activity data operations.
type ActivityStore struct {
	db *gorm.DB
}

// NewActivityStore creates a new ActivityStore instance.
func NewActivityStore(db *gorm.DB) *ActivityStore {
	return &ActivityStore{db: db}
}

// CreateRecord creates a new activity record.
func (s *ActivityStore) CreateRecord(post ActivityPost) error {
	record := models.Activity{
		UserID:    post.UserID,
		GroupID:   post.GroupID,
		MessageID: post.MessageID,
		Activity:  post.Activity,
		Month:     post.Month,
		Day:       post.Day,
		Year:      post.Year,
	}
	return s.db.Create(&record).Error
}

// UpdateRecord updates an existing activity record by message ID.
func (s *ActivityStore) UpdateRecord(post ActivityPost) error {
	var existingRecord models.Activity
	err := s.db.Where("user_id = ? AND group_id = ? AND message_id = ?",
		post.UserID, post.GroupID, post.MessageID).First(&existingRecord).Error

	if err != nil {
		return err
	}

	existingRecord.Activity = post.Activity
	existingRecord.Month = post.Month
	existingRecord.Day = post.Day
	existingRecord.Year = post.Year
	return s.db.Save(&existingRecord).Error
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