# chat-app

# gateway-proxy ---> locahost:80
# api-gateway-1 & 2 ---> locahost:8080
# service-proxy ---> locahost:3000
# grpc-proxy ---> locahost:3001

# user-service ---> locahost:5000 & service-proxy/user-service 
# relay-service ---> locahost:7400 & service-proxy/relay-service 
# group-service ---> locahost:6000 & service-proxy/group-service 
# discovery-service ---> locahost:7101 & service-proxy/discovery-service
# chat-service ---> locahost:7000 & service-proxy/chat-service


### GRPC

# discovery-service ---> locahost:7100 & grpc-proxy/discovery-service
# websocket_manager-service ---> locahost:7200 & grpc-proxy/websocket_manager-service
# message-service ---> locahost:7300 & grpc-proxy/message-service
# group-service locahost:6001 & grpc-proxy/group-service
# group-message-service locahost:7500 & grpc-proxy/group--message-service




### monitoring

# promethues -> http://localhost:9090
# zookeeper-UI -> http://localhost:9001
# Queues-UI -> http://localhost:9000
