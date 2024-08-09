PROTO_COMMAND = buf generate

discovery-gen:
		$(PROTO_COMMAND) && mv 'discovery-grpc-gen' '../discovery-service/internal/gRPC'

chat-gen:
		$(PROTO_COMMAND) && cp -r ./discovery-grpc-gen '../chat-service/internal/gRPC' && cp -r ./websocket-manager-grpc-gen '../chat-service/internal/gRPC'

websocket-gen:
		$(PROTO_COMMAND) && cp -r ./websocket-manager-grpc-gen '../websocket-manager/internal/gRPC'

tools:
		bash tools.sh
