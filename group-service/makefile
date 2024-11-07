BUILD_COMMAND := go build -o bin/group-service ./cmd/http

build:
		$(BUILD_COMMAND)

dev-stack:
		docker compose -f ./deployment/docker-compose.dev.yaml up

tools:
		bash tools.sh