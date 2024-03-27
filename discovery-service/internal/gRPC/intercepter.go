package gRPC

import (
	"context"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ValidationInterceptor is a gRPC server interceptor that validates incoming requests.
func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Initialize the protovalidate validator
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	// Assert that req is a ProtoMessage
	protoMsg, ok := req.(protoreflect.ProtoMessage)
	if !ok {
		// Handle the case where req is not a ProtoMessage
		return nil, status.Errorf(codes.Aborted, "request does not implement protoreflect.ProtoMessage")
	}

	// Perform validation
	if err := v.Validate(protoMsg); err != nil {
		// Convert the error to a gRPC status error
		// fmt.Printf("error: %v\n",  err.Error())
		responseMap := map[string]interface{}{
			"message": err.Error(),
			"fields":  make(map[string]string),
		}
		splitErr := strings.Split(err.Error(), "-")
		responseMap["message"] = splitErr[0]
		for i := 1; i < len(splitErr); i++ {
			fieldError := strings.Split(splitErr[i], ":")
			responseMap["fields"].(map[string]string)[fieldError[0]] = fieldError[1]
		}

		// fmt.Print(responseMap)
		// responseStr,_ := json.Marshal(responseMap)
		// return nil, status.Errorf(codes.InvalidArgument, string(responseStr))
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Proceed to the handler if validation is successful
	return handler(ctx, req)
}
