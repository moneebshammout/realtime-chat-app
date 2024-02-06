package discovery

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// unaryInterceptor is a gRPC server interceptor that validates incoming requests
func ValidationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Perform validation
	if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			// Return validation error to the client
			return nil, status.Errorf(status.Code(err), err.Error())
		}
	}

	// Continue to the handler if validation is successful
	return handler(ctx, req)
}
