package models

type User struct {
	ID      int64  `gorm:"primaryKey"`
	Name    string `gorm:"size:100;not null"`
	GroupID int64  `gorm:"not null"`
}
