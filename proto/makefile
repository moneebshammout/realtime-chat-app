PROTO_COMMAND = buf generate

discovery-gen:
		$(PROTO_COMMAND) && mv 'discovery-grpc-gen' '../discovery-service/internal/gRPC'

chat-gen:
		$(PROTO_COMMAND) && cp -r ./discovery-grpc-gen '../chat-service/internal/gRPC' && cp -r ./websocket-manager-grpc-gen '../chat-service/internal/gRPC' && cp -r ./message-service-grpc-gen '../chat-service/internal/gRPC'

websocket-gen:
		$(PROTO_COMMAND) && cp -r ./websocket-manager-grpc-gen '../websocket-manager/internal/gRPC'

message-gen:
		$(PROTO_COMMAND) && cp -r ./message-service-grpc-gen '../message-service/internal/gRPC'

tools:
		bash tools.sh
