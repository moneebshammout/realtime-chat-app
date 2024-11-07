package DAO

import (
	"group-service/internal/database"
	"group-service/internal/database/models"
)

type GroupDAOStruct struct {
	BaseDAO
}


func group() *GroupDAOStruct {
	return &GroupDAOStruct{
		BaseDAO: BaseDAO{
			dbClient: database.GetClient(),
			model:    models.Group{},
		},
	}
}

var Group *GroupDAOStruct

func init() {
	database.Connect()
	Group = group()
}