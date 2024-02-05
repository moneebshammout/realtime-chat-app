BUILD_COMMAND := go build -o bin/discovery-service ./cmd/http

build:
		$(BUILD_COMMAND)

dev-stack:
		docker compose -f ./deployment/docker-compose.dev.yaml up


migrations:
		$(PRISMA) migrate dev --schema=$(PRISMA_PATH)

migrate:
		$(PRISMA) migrate deploy --schema=$(PRISMA_PATH)


proto:
	protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    internal/gRPC/discovery.proto