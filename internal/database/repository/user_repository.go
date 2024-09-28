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
	err := r.DB.First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{
				ID:      userID,
				Name:    name,
				GroupID: groupID,
			}

			if err := r.DB.Create(&user).Error; err != nil {
				return nil, err
			}

			return &user, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(userID int64) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
