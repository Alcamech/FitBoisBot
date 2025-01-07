package models

import "fmt"

type User struct {
	Name      string     `gorm:"size:100;not null"`
	Group     Group      `gorm:"foreignKey:GroupID;references:ID"` // Reference to Group
	GGRecords []Gg       `gorm:"foreignKey:UserID"`
	Activites []Activity `gorm:"foreignKey:UserID"`
	ID        int64      `gorm:"primaryKey;autoIncrement:false"`
	GroupID   int64      `gorm:"not null;index"`
}

func (u User) ToString() string {
	return fmt.Sprintf("User[ID=%d, Name=%s, GroupID=%d]", u.ID, u.Name, u.GroupID)
}
