package DAO

import (
	"group-service/internal/database"
	"group-service/internal/database/models"
)

type GroupMemeberDAOStruct struct {
	BaseDAO[models.GroupMember]
}


func groupMember() *GroupMemeberDAOStruct {
	return &GroupMemeberDAOStruct{
		BaseDAO: BaseDAO[models.GroupMember]{
			dbClient: database.GetClient(),
		},
	}
}

var GroupMember *GroupMemeberDAOStruct

func init() {
	database.Connect()
	GroupMember = groupMember()
}