package models

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	GroupMembers []GroupMember
}

func (Group) TableName() string {
	return "groups"
}
