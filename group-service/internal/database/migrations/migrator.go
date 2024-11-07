package migrations

import (
	"group-service/internal/database"
	"group-service/internal/database/models"
	"group-service/pkg/utils"
)

var logger = utils.GetLogger()

func Migrate() {
	logger.Info("Migrating database...")
	err := database.GetClient().Session.AutoMigrate(&models.Group{},&models.GroupMember{})
	if err != nil {
		logger.Panic(err)
	}

	logger.Info("Database migrated successfully")
}
