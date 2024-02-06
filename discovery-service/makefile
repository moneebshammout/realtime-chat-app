BUILD_COMMAND := go build -o bin/discovery-service ./cmd/http
TOOLS_COMMAND := bash tools.sh

build:
		$(BUILD_COMMAND)

dev-stack:
		docker compose -f ./deployment/docker-compose.dev.yaml up

proto:
 		protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		--validate_out=lang=go:. \
		--validate_opt=paths=source_relative \
		internal/gRPC/discovery.proto

tools:
		bash tools.sh
