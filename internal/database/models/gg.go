package models

type Gg struct {
	ID          int64  `gorm:"primaryKey"`
	UserID      int64  `gorm:"not null"`
	GroupID     int64  `gorm:"not null"`
	Year        string `gorm:"size:4"`
	FastGGCount int    `gorm:"default:0"`
}
