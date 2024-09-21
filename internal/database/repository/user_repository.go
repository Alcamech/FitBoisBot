package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database"
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

func GetOrCreateUser(userID int64, name string, groupID int64) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, userID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user = models.User{
				ID:      userID,
				Name:    name,
				GroupID: groupID,
			}

			if err := database.DB.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}

	return &user, nil
}
