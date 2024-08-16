package database

import (
	"fmt"
	"strings"
	"time"

	dbConfig "relay-service/config/db"
	"relay-service/pkg/utils"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

var (
	logger   = utils.GetLogger()
	dbClient *DBClient
)

type DBClient struct {
	Session gocqlx.Session
}

func getCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster(strings.Split(dbConfig.Env.ClusterUrl, ",")...)
	cluster.NumConns = 5
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Timeout = time.Second * 10
	cluster.ConnectTimeout = time.Second * 10

	return cluster
}

func Connect() {
	logger.Infof("DBClient: Connecting to database %v", len(strings.Split(dbConfig.Env.ClusterUrl, ",")))

	for {
		cluster := getCluster()
		initialSession, err := cluster.CreateSession()
		if err != nil {
			logger.Errorf("DBClient: Error connecting to database: %s", err)
			logger.Infof("DBClient: Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Create keyspace if it does not exist
		keyspace := dbConfig.Env.KeySpace
		createKeyspaceQuery := fmt.Sprintf(`
			CREATE KEYSPACE IF NOT EXISTS %s
			WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
		`, keyspace)

		logger.Infof("DBClient: Creating keyspace %s", keyspace)
		if err := initialSession.Query(createKeyspaceQuery).Exec(); err != nil {
			logger.Panicf("Failed to create keyspace %s: %v", keyspace, err)
		}

		// Close initial session
		initialSession.Close()

		logger.Infof("DBClient: initializing session %s", keyspace)
		// Now create a session with the keyspace specified
		cluster = getCluster()
		cluster.Keyspace = keyspace
		session, err := gocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			logger.Errorf("DBClient: Error connecting to database with keyspace: %s", err)
			logger.Infof("DBClient: Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Info("DBClient: Connected to database with keyspace")
		dbClient = &DBClient{
			Session: session,
		}

		break
	}
}

func Disconnect() {
	dbClient.Session.Close()
}

func Client() *DBClient {
	return dbClient
}
