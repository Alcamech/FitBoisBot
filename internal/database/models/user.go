package models

import "fmt"

type User struct {
	Name      string     `gorm:"size:100;not null"`
	GGRecords []Gg       `gorm:"foreignKey:UserID"`
	Activites []Activity `gorm:"foreignKey:UserID"`
	ID        int64      `gorm:"primaryKey"`
	GroupID   int64      `gorm:"not null"`
}

func (u User) ToString() string {
	return fmt.Sprintf("User[ID=%d, Name=%s, GroupID=%d]", u.ID, u.Name, u.GroupID)
}
