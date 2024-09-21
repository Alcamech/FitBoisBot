package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	DB *gorm.DB
}

func (r *ActivityRepository) CreateActivityRecord(userID int64, activity, month, day, year string) error {
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
