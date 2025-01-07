package repository

import (
	"fmt"

	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	DB *gorm.DB
}

func (r *ActivityRepository) CreateActivityRecord(userID, groupID int64, activity, month, day, year string) error {
	record := models.Activity{
		UserID:   userID,
		GroupID:  groupID,
		Activity: activity,
		Month:    month,
		Day:      day,
		Year:     year,
	}

	return r.DB.Create(&record).Error
}

func (r *ActivityRepository) GetUserActivityByYear(userID, groupID int64, year string) ([]models.Activity, error) {
	var records []models.Activity
	err := r.DB.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).Find(&records).Error
	return records, err
}

func (r *ActivityRepository) GetActivityCountByUserIdAndMonth(userID, groupID int64, month string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Activity{}).Where("user_id = ? AND group_id = ? AND month = ?", userID, groupID, month).Count(&count).Error
	return count, err
}

func (r *ActivityRepository) GetUsersWithActivities(groupID int64) ([]int64, error) {
	var userIDs []int64
	err := r.DB.Model(&models.Activity{}).Where("group_id = ?", groupID).Distinct("user_id").Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *ActivityRepository) GetActivityCountByUserIdAndMonthAndYear(userID, groupID int64, month, year string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Activity{}).
		Where("user_id = ? AND group_id = ? AND month = ? AND year = ?", userID, groupID, month, year).
		Count(&count).Error
	return count, err
}

func (r *ActivityRepository) GetMostActiveUserForMonth(groupID int64, month, year string) (int64, int64, error) {
	var result struct {
		UserID int64
		Count  int64
	}
	err := r.DB.Model(&models.Activity{}).
		Select("user_id, COUNT(*) as count").
		Where("group_id = ? AND month = ? AND year = ?", groupID, month, year).
		Group("user_id").
		Order("count DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}

	if result.UserID == 0 {
		return 0, 0, fmt.Errorf("no activities found for group %d in %s/%s", groupID, month, year)
	}

	return result.UserID, result.Count, nil
}
