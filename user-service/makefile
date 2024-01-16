BUILD_COMMAND := go build -o bin/user-service ./cmd/http
PRISMA := go run github.com/steebchen/prisma-client-go
PRISMA_PATH := internal/db/schema.prisma
# server:
# 	CompileDaemon -build="$(BUILD_COMMAND)" -command="godotenv user-service" 

build:
		$(BUILD_COMMAND)

dev-stack:
		docker compose -f ./deployment/docker-compose.dev.yaml up


migrations:
		$(PRISMA) migrate dev --schema=$(PRISMA_PATH)

migrate:
		$(PRISMA) migrate deploy --schema=$(PRISMA_PATH)