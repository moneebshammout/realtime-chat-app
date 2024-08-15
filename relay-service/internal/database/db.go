// package database

// import (
// 	"context"
// "github.com/scylladb/gocqlx/v3"
// "github.com/gocql/gocql"

// 	dbConfig "relay-service/config/db"
// 	"relay-service/pkg/utils"

// )
// var logger = utils.GetLogger()
// type DBClient struct {
// 	Session *gocqlx.Session
// }


// func Connect() (*DBClient, error) {
// 	cluster := gocql.NewCluster( dbConfig.Env.ClusterUrl)
// 	session, err := gocqlx.WrapSession(cluster.CreateSession())

// 	if err != nil {
// 		logger.Panicf("Faild To Connect To Postgress Database: %s", err)
// 	}

	
// }

// func Disconnect() {

// }
