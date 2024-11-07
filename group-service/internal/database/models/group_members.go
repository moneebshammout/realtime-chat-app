package models

import (
	"gorm.io/gorm"
)

type GroupMember struct {
	gorm.Model
	SourceId string `gorm:"uniqueIndex:sourceId_groupId;not null" json:"sourceId"`
	GroupId  uint   `gorm:"uniqueIndex:sourceId_groupId;not null" json:"groupId"`
	Role     string `gorm:"not null" json:"role"`
	Group    Group  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" `
}

func (GroupMember) TableName() string {
	return "group_members"
}

const (
	OWNER  = "OWNER"
	Admin  = "ADMIN"
	Member = "MEMBER"
)

var MemberRoles = map[string]bool{
	OWNER:  true,
	Admin:  true,
	Member: true,
}
