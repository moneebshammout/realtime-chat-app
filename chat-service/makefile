BUILD_COMMAND := go build -o bin/chat-service ./cmd/http
PRISMA := go run github.com/steebchen/prisma-client-go
PRISMA_PATH := internal/db/schema.prisma
# server:
# 	CompileDaemon -build="$(BUILD_COMMAND)" -command="godotenv chat-service" 

build:
		$(BUILD_COMMAND)

migrations:
		$(PRISMA) migrate dev --schema=$(PRISMA_PATH)

migrate:
		$(PRISMA) migrate deploy --schema=$(PRISMA_PATH)