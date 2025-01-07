package models

import "time"

type Token struct {
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Year      string    `gorm:"size:4;not null"`
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"`
	GroupID   int64     `gorm:"not null;index"`
	Balance   int       `gorm:"not null;default:0"`
}
