PROTO_COMMAND = buf generate

discovery-gen:
	$(PROTO_COMMAND) 
	$(PROTO_COMMAND) --template buf.gateway.gen.yaml --path ./protos/discovery.proto --output ./discovery-grpc-gen
	cp -r 'discovery-grpc-gen/protos/discovery.pb.gw.go' 'discovery-grpc-gen'
	cp -r 'discovery-grpc-gen/swagger/protos/discovery.swagger.json' 'discovery-grpc-gen'
	rm -rf 'discovery-grpc-gen/protos'
	rm -rf 'discovery-grpc-gen/swagger'
	cp -r './discovery-grpc-gen' '../discovery-service/internal/gRPC'

chat-gen:
	$(MAKE) discovery-gen
	cp -r ./discovery-grpc-gen '../chat-service/internal/gRPC'
	cp -r ./websocket-manager-grpc-gen '../chat-service/internal/gRPC'
	cp -r ./message-service-grpc-gen '../chat-service/internal/gRPC'

websocket-gen:
	$(PROTO_COMMAND)
	cp -r ./websocket-manager-grpc-gen '../websocket-manager/internal/gRPC'

message-gen:
	$(PROTO_COMMAND)
	cp -r ./message-service-grpc-gen '../message-service/internal/gRPC'

tools:
	bash tools.sh
