package models

import (
	"group-service/internal/database"
	"group-service/pkg/utils"

	"gorm.io/gorm"
)

var logger = utils.GetLogger()

type baseModel struct {
	gorm.Model
	dbClient *database.DBClient `gorm:"-:all" json:"-"`
}

func (b *baseModel) Save() {
	b.dbClient.Session.Create(&b)
	logger.Infof("Model saved: %v", b)
}
