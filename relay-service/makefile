BUILD_COMMAND := go build -o bin/discovery-service ./cmd/http
TOOLS_COMMAND := bash tools.sh
include .env
export $(shell sed 's/=.*//' .env)

build:
		$(BUILD_COMMAND)

dev-stack:
		docker compose -f ./deployment/docker-compose.dev.yaml up
		
tools:
		bash tools.sh

migration:
		touch internal/database/migrations/$(shell date +"%Y%m%d%H%M%S")_$(name).cql


migrate:
		go build -o ./bin/migrate cmd/migrate/migrate.go &&  ./bin/migrate


schema:
	echo "Schema generated successfully $(DATABASE_CLUSTER_HOST)"
	./bin/schemagen -cluster=$(DATABASE_CLUSTER_HOST) -keyspace=$(DATABASE_KEYSPACE) -output="internal/database/models" -pkgname="models"