package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	DB *gorm.DB
}

func (r *ActivityRepository) CreateActivityRecord(userID int64, activity string, month string, day string, year string) error {
	record := models.Activity{
		UserID:   userID,
		Activity: activity,
		Month:    month,
		Day:      day,
		Year:     year,
	}

	return r.DB.Create(&record).Error
}

func (r *ActivityRepository) GetUserActivityByYear(userID int64, year string) ([]models.Activity, error) {
	var records []models.Activity
	err := r.DB.Where("user_id = ? AND year = ?", userID, year).Find(&records).Error
	return records, err
}

func (r *ActivityRepository) GetActivityCountByUserIdAndMonth(userID int64, month string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Activity{}).Where("user_id = ? AND month = ?", userID, month).Count(&count).Error
	return count, err
}

func (r *ActivityRepository) GetUsersWithActivities() ([]int64, error) {
	var userIDs []int64
	err := r.DB.Model(&models.Activity{}).Distinct("user_id").Pluck("user_id", &userIDs).Error
	return userIDs, err
}
