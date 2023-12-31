package config

import (
	"context"
	"fmt"

	"user-service/internal/database/db"
)

type PrismaDB struct {
	Client  *db.PrismaClient
	Context context.Context
}

var clientDB = &PrismaDB{}

func DBConnect() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(fmt.Sprintf("Faild To Connect To Postgress Database: %s", err))
	}

	clientDB.Client = client
	clientDB.Context = context.Background()
}

func DB() *PrismaDB {
	return clientDB
}

func KillDBConnection() {
	clientDB.Client.Disconnect()
}
