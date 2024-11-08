package DAO

import (
	"group-service/internal/database"
	"group-service/internal/database/models"
)

type GroupDAOStruct struct {
	BaseDAO[models.Group]
}

func group() *GroupDAOStruct {
	return &GroupDAOStruct{
		BaseDAO: BaseDAO[models.Group]{
			dbClient: database.GetClient(),
		},
	}
}

var Group *GroupDAOStruct

func init() {
	database.Connect()
	Group = group()
}
