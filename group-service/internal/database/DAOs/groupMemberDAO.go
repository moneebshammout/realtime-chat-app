package DAO

import (
	"group-service/internal/database"
	"group-service/internal/database/models"
)

type GroupMemeberDAOStruct struct {
	BaseDAO
}


func groupMember() *GroupMemeberDAOStruct {
	return &GroupMemeberDAOStruct{
		BaseDAO: BaseDAO{
			dbClient: database.GetClient(),
			model:    models.Group{},
		},
	}
}

var GroupMember *GroupMemeberDAOStruct

func init() {
	database.Connect()
	GroupMember = groupMember()
}