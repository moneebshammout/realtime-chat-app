// package gRPC

// import (
// 	"context"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/status"
// 	"github.com/bufbuild/protovalidate-go"
// )

// // unaryInterceptor is a gRPC server interceptor that validates incoming requests
// func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	// Perform validation
// 	v, err := protovalidate.New()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if _, ok := req.(interface{ Validate() error }); ok {
// 		if err := v.Validate(req); err != nil {
// 			// Return validation error to the client
// 			return nil, status.Errorf(status.Code(err), err.Error())
// 		}
// 	}

// 	// Continue to the handler if validation is successful
// 	return handler(ctx, req)
// }

package gRPC

import (
	"context"

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
		return nil, status.Errorf(status.Code(err), err.Error())
	}

	// Proceed to the handler if validation is successful
	return handler(ctx, req)
}
