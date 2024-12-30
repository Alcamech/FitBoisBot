package models

type Gg struct {
	Year        string `gorm:"size:4"`
	ID          int64  `gorm:"primaryKey"`
	UserID      int64  `gorm:"not null"`
	GroupID     int64  `gorm:"not null"`
	FastGGCount int    `gorm:"default:0"`
}
