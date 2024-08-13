package config

import (
	"context"

	"user-service/internal/database/db"
)

type PrismaDB struct {
	Client  *db.PrismaClient
	Context context.Context
	Cancel  context.CancelFunc
}

var clientDB = &PrismaDB{}

func DBConnect() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		logger.Panicf("Faild To Connect To Postgress Database: %s", err)
	}

	clientDB.Client = client
	clientDB.Context, clientDB.Cancel = context.WithCancel(context.Background())
}

func DB() *PrismaDB {
	return clientDB
}

func KillDBConnection() {
	clientDB.Client.Prisma.Disconnect()
	clientDB.Client.Disconnect()
	clientDB.Cancel()
}
