package database

import (
	"sync"
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
	once     sync.Once
)

func Connect() {
	once.Do(func() {
		logger.Infof("DBClient: Connecting to database ")

		for {
			db, err := gorm.Open(postgres.Open(config.Env.PostgresAddr), &gorm.Config{})
			if err != nil {
				logger.Errorf("DBClient: Error connecting to database: %s", err)
				logger.Infof("DBClient: Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}

			err = db.Exec("CREATE SCHEMA IF NOT EXISTS group_service").Error
			if err != nil {
				logger.Errorf("DBClient: Error creating schema: %s", err)
				logger.Infof("DBClient: Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}

			sqlDB, err := db.DB()
			if err != nil {
				logger.Errorf("DBClient: Error getting database handle: %s", err)
				logger.Infof("DBClient: Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}

			if err = sqlDB.Ping(); err != nil {
				logger.Errorf("DBClient: Error pinging database: %s", err)
				logger.Infof("DBClient: Retrying in 5 seconds...")
				time.Sleep(5 * time.Second)
				continue
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)

			logger.Info("DBClient: Connected to database with keyspace")
			dbClient = &DBClient{
				Session: db,
			}

			break
		}
	})
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
