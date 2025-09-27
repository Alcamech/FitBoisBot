package models

import "time"

type Group struct {
	Timezone  string    `gorm:"size:255;default:'America/New_York'"`
	ID        int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
