
# Distrebuted Real-Time Chat App

A scalable, distributed real-time chat application built with **Go**, **gRPC**, **Nginx**, **PostgreSQL**, **ScyllaDB**, **Redis**, and **Apache Zookeeper** for service discovery. Dockerized for easy deployment and management.

## Features
- **Real-Time Messaging**: Instant messaging through WebSockets.
- **Microservice Architecture**: Go microservices with gRPC for fast communication.
- **Service Discovery**: Apache Zookeeper for discovering chat servers in a distributed environment.
- **Caching**: Redis for caching to improve performance.
- **Databases**: PostgreSQL for relational data and ScyllaDB for scalable NoSQL storage.
- **Proxying**: Nginx for traffic management and reverse proxying.

## Technologies
- **Go** for web server and microservices.
- **gRPC** for efficient communication.
- **PostgreSQL** for relational data storage.
- **ScyllaDB** for fast, scalable NoSQL storage.
- **Redis** for caching.
- **Apache Zookeeper** for service discovery.
- **Nginx** for reverse proxy.
- **Docker** for containerization.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/moneebshammout/realtime-chat-app.git
   ```

2. Build And Run the infrastructure stack:
   ```bash
   docker-compose -f docker.compose.stack.yaml up
   ```

3. Then build and Run the whole system:
   ```bash
   docker-compose up
   ```

![chat-app-design drawio](https://github.com/user-attachments/assets/9390fbba-059b-4be1-8044-19af766fcce5)

## Services Overview

### 1. **User Service**
Handles user management, authentication, and authorization. Ensures users can register, log in, and access their profiles securely.

### 2. **Group Service**
Manages groups and the addition/removal of members. This service allows users to create groups and manage group memberships.

### 3. **Last Seen Service**
Tracks the last online status of each user. It helps to show whether a user is currently online or when they were last active.

### 4. **Chat Service**
The backbone of the system. This WebSocket server facilitates communication between users. Users connect to it to start chats, send/receive messages, and maintain communication.

### 5. **Discovery Service**
This service handles service discovery, using **Apache Zookeeper** to track available chat servers. It ensures that users are always routed to an available server in the distributed system.

### 6. **Relay Service**
Manages unsent messages when a user is offline. It temporarily stores messages and delivers them when the recipient is online.

### 7. **Group Message Service**
Manages and stores messages sent within groups. This service ensures that all group members receive messages and stores them in the database.

### 8. **Message Service**
Handles one-to-one messaging between users. It is responsible for sending, receiving, and storing private messages.

### 9. **Media Service**
Handles uploading, storing, and serving media files like images, audio, and video within chats.

### 10. **Notification Service**
Sends notifications to users about new messages, user status changes, and other relevant events.

### 11. **Websocket Manager Service**
Manage Websocket sessions and keeps references for each connection alive


## URLs
- **Gateway Proxy**: `http://localhost`
- **API Gateway**: `http://localhost:8080`
- **Service Proxy**: `http://localhost:3000`
- **GRPC Proxy**: `http://localhost:3001`

- **User Service (REST API)**: `http://localhost/user-service`
- **Relay Service (REST API)**: `http://localhost/relay-service`
- **Group Service (REST API)**: `http://localhost/group-service`
- **Discovery Service (REST API)**: `http://localhost/discovery-service`
- **Chat Service (REST API)**: `http://localhost/chat-service`
- **Last Seen Service (REST API)**: `http://localhost/last-seen-service`

- **Discovery Service (gRPC)**: `http://localhost:3001/Discovery`
- **Websocket Manager Service (gRPC)**: `http://localhost:3001/WebsocketManager.WebsocketManager`
- **Message Service (gRPC)**: `http://localhost:3001/MessageService.MessageService`
- **Group Service (gRPC)**: `http://localhost:3001/GroupService.GroupService`
- **Group Message Service (gRPC)**: `http://localhost:3001/GroupMessageService.GroupMessageService`

- **Redis**: `http://localhost:6379`
- **PostgreSQL**: `http://localhost:5432`
- **ScyllaDB**: `http://localhost:9042`

## Monitoring and Logs
- **Prometheus UI**: [localhost:9090](http://localhost:9090)
- **Zookeeper UI**: [localhost:9001](http://localhost:9001)
- **Queues UI**: [localhost:9000](http://localhost:9000)

## Conclusion
This real-time chat application utilizes modern technologies to provide an efficient and scalable messaging solution. With features like service discovery via Apache Zookeeper, WebSocket-based messaging, and microservices architecture, it is designed to handle high traffic and provide a seamless user experience.
