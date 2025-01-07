package models

type Gg struct {
	Year        string `gorm:"size:4"`
	User        User   `gorm:"foreignKey:UserID;references:ID"`
	Group       Group  `gorm:"foreignKey:GroupID;references:ID"`
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	UserID      int64  `gorm:"not null;index"`
	GroupID     int64  `gorm:"not null"`
	FastGGCount int    `gorm:"default:0"`
}
