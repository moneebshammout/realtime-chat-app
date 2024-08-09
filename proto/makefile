PROTO_COMMAND = buf generate

discovery-gen:
		$(PROTO_COMMAND) && mv 'discovery-grpc-gen' '../discovery-service/internal/gRPC'

chat-gen:
		$(PROTO_COMMAND) && cp -r ./discovery-grpc-gen '../chat-service/internal/gRPC'

tools:
		bash tools.sh
