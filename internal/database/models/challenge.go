package models

import "time"

type Challenge struct {
	ID           int64                  `gorm:"primaryKey;autoIncrement"`
	GroupID      int64                  `gorm:"not null;index:idx_group_status"`
	CreatorID    int64                  `gorm:"not null;index:idx_creator"`
	Title        string                 `gorm:"size:255;not null"`
	Description  string                 `gorm:"type:text"`
	Difficulty   string                 `gorm:"size:20;not null"`
	Multiplier   float64                `gorm:"not null"`
	Status       string                 `gorm:"size:20;not null;default:'pending';index:idx_group_status"`
	StartDate    *time.Time             // When first non-creator participant joins (activates challenge)
	CreatedAt    time.Time              `gorm:"autoCreateTime"`
	UpdatedAt    time.Time              `gorm:"autoUpdateTime"`
	Group        Group                  `gorm:"foreignKey:GroupID;references:ID"`
	Creator      User                   `gorm:"foreignKey:CreatorID;references:ID"`
	Participants []ChallengeParticipant `gorm:"foreignKey:ChallengeID"`
}
