package models

import "time"

type ChallengeParticipant struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	ChallengeID int64     `gorm:"not null;uniqueIndex:unique_challenge_user"`
	UserID      int64     `gorm:"not null;index:idx_user;uniqueIndex:unique_challenge_user"`
	WagerAmount int       `gorm:"not null"`
	Score       int       `gorm:"not null;default:0"`
	IsWinner    *bool     `gorm:"default:null"`
	JoinedAt    time.Time `gorm:"not null;autoCreateTime"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Challenge   Challenge `gorm:"foreignKey:ChallengeID;references:ID"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
}
