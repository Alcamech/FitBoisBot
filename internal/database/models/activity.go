package models

type Activity struct {
	Activity string `gorm:"size:255"`
	Month    string `gorm:"size 2"`
	Day      string `gorm:"size 2"`
	Year     string `gorm:"size 2"`
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	UserID   int64  `gorm:"not null;index"`
	GroupID  int64  `gorm:"not null"`
}
