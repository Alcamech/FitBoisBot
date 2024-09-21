package models

import "fmt"

type User struct {
	ID        int64      `gorm:"primaryKey"`
	Name      string     `gorm:"size:100;not null"`
	GroupID   int64      `gorm:"not null"`
	Activites []Activity `gorm:"foreignKey:UserID"`
	GGRecords []Gg       `gorm:"foreignKey:UserID"`
}

func (u User) ToString() string {
	return fmt.Sprintf("User[ID=%d, Name=%s, GroupID=%d]", u.ID, u.Name, u.GroupID)
}
