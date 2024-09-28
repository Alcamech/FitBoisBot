package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type GGRepository struct {
	DB *gorm.DB
}

func (r *GGRepository) CreateOrUpdateGGCount(userID, groupID int64, year string) error {
	var gg models.Gg
	err := r.DB.Where("user_id = ? AND group_id = ? AND year = ?", userID, groupID, year).First(&gg).Error
	if err == gorm.ErrRecordNotFound {
		gg = models.Gg{
			UserID:      userID,
			GroupID:     groupID,
			Year:        year,
			FastGGCount: 1,
		}
		return r.DB.Create(&gg).Error
	}

	if err != nil {
		return err
	}

	gg.FastGGCount++
	return r.DB.Save(&gg).Error
}

func (r *GGRepository) GetGGCount(userID int64, year string) (int, error) {
	var gg models.Gg
	err := r.DB.Where("user_id = ? AND year = ?", userID, year).First(&gg).Error
	return gg.FastGGCount, err
}
