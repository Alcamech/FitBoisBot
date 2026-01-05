package models

import "time"

type UserBalance struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;uniqueIndex:idx_user_group"`
	GroupID   int64     `gorm:"not null;uniqueIndex:idx_user_group"`
	Balance   int       `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
