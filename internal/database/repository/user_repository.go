package repository

import (
	"github.com/Alcamech/FitBoisBot/internal/database/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) GetOrCreateUser(userID int64, name string, groupID int64) (*models.User, error) {
	var user models.User
	result := r.DB.First(&user, userID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user = models.User{
				ID:      userID,
				Name:    name,
				GroupID: groupID,
			}

			if err := r.DB.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}

	return &user, nil
}
