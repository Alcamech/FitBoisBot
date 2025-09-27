package models

import "time"

type Activity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null"`
	GroupID   int64     `gorm:"not null"`
	MessageID int       `gorm:"uniqueIndex:unique_message,priority:3;not null"`
	Activity  string    `gorm:"size:255;not null"`
	Month     string    `gorm:"size:2;not null"`
	Day       string    `gorm:"size:2;not null"`
	Year      string    `gorm:"size:4;not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Group     Group     `gorm:"foreignKey:GroupID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
