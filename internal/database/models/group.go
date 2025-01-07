package models

type Group struct {
	Timezone string `gorm:"size:255;default:'America/New_York'"`
	ID       int64  `gorm:"primaryKey"`
}
