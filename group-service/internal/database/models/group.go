package models

import (
	"group-service/internal/database"
)

type Group struct {
	Name string `gorm:"not null" json:"name"`
	baseModel
}

func (Group) TableName() string {
	return "groups"
}

func GroupModel() *Group {
	return &Group{
		baseModel: baseModel{
			dbClient: database.GetClient(),
		},
	}
}
