package models

type Activity struct {
	ID       int64  `gorm:"primaryKey"`
	UserID   int64  `gorm:"not null"`
	Activity string `gorm:"size:255"`
	Month    string `gorm:"size 2"`
	Day      string `gorm:"size 2"`
	Year     string `gorm:"size 2"`
}
