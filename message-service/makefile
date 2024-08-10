BUILD_COMMAND := go build -o bin/message-service ./cmd/http
TOOLS_COMMAND := bash tools.sh

build:
		$(BUILD_COMMAND)

dev-stack:
		docker compose -f ./deployment/docker-compose.dev.yaml up

tools:
		bash tools.sh
