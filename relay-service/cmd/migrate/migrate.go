package main

import (
	"context"
	"os"
	"strings"
	"time"

	"relay-service/pkg/utils"

	"relay-service/internal/database/migrations"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
)

var (
	logger     = utils.InitLogger()
	_          = godotenv.Load()
	clusterUrl = os.Getenv("DATABASE_CLUSTER_HOST")
	keySpace   = os.Getenv("DATABASE_KEYSPACE")
)

func main() {
	logger.Infof("%s,,,,,,,,,,,,,%s", clusterUrl, keySpace)
	cluster := gocql.NewCluster(strings.Split(clusterUrl, ",")...)
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Timeout = time.Second * 10
	cluster.ConnectTimeout = time.Second * 10
	cluster.Keyspace = keySpace
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		logger.Fatal("Migrate:", err)
	}

	// Run migrations
	migrater(session)
	session.Close()
}

func migrater(session gocqlx.Session) {
	log := func(ctx context.Context, session gocqlx.Session, ev migrate.CallbackEvent, name string) error {
		logger.Infof("Migrate: %v %s", ev, name)
		return nil
	}
	reg := migrate.CallbackRegister{}
	reg.Add(migrate.BeforeMigration, "m1.cql", log)
	reg.Add(migrate.AfterMigration, "m1.cql", log)
	reg.Add(migrate.CallComment, "1", log)
	reg.Add(migrate.CallComment, "2", log)
	reg.Add(migrate.CallComment, "3", log)
	migrate.Callback = reg.Callback

	// Ensure the session is using the correct keyspace
	if err := migrate.FromFS(context.TODO(), session, migrations.Files); err != nil {
		logger.Fatal("Migrate:", err)
	}
}
