package models

import (
	"reflect"

	"relay-service/internal/database"
)

func MessageDAO() *model {
	return &model{
		dbClient:  database.Client(),
		model:     Messages,
		modelType: reflect.TypeOf(MessagesStruct{}),
	}
}
