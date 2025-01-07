package models

type Activity struct {
	Activity string `gorm:"size:255"`
	Month    string `gorm:"size 2"`
	Day      string `gorm:"size 2"`
	Year     string `gorm:"size 2"`
	User     User   `gorm:"foreignKey:UserID;references:ID"`
	Group    Group  `gorm:"foreignKey:GroupID;references:ID"`
	ID       int64  `gorm:"primaryKey"`
	UserID   int64  `gorm:"not null"`
	GroupID  int64  `gorm:"not null"`
}
