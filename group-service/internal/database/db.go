package database

import (
	"time"

	"group-service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"group-service/pkg/utils"
)

type DBClient struct {
	Session *gorm.DB
}

var (
	logger   = utils.GetLogger()
	dbClient *DBClient
)

func Connect() {
	logger.Infof("DBClient: Connecting to database ")

	for {
		db, err := gorm.Open(postgres.Open(config.Env.PostgresAddr), &gorm.Config{
			// ConnPool: ,
		})
		if err != nil {
			logger.Errorf("DBClient: Error connecting to database: %s", err)
			logger.Infof("DBClient: Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Info("DBClient: Connected to database with keyspace")
		dbClient = &DBClient{
			Session: db,
		}

		break
	}
}

func GetClient() *DBClient {
	return dbClient
}

func Disconnect() {
	if db, err := dbClient.Session.DB(); err != nil {
		logger.Errorf("DBClient: Error disconnecting from database: %s", err)
	} else {
		logger.Info("DBClient: Disconnected from database")
		db.Close()
	}
}
