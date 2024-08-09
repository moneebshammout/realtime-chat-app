package gRPC

// func (s DiscoveryServiceServer) UnaryServerInterceptor(ctx context.Context,
// 	req interface{},
// 	info *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (interface{}, error) {
// 	// interceptor logic
// 	return handler(ctx, req)
// }
import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"discovery-service/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ValidationInterceptor is a gRPC server interceptor that validates incoming requests.
func validationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	protoMsg, ok := req.(protoreflect.ProtoMessage)
	if !ok {
		return nil, status.Errorf(codes.Aborted, "request does not implement protoreflect.ProtoMessage")
	}

	if err := v.Validate(protoMsg); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return handler(ctx, req)
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// return handler(ctx, req)
	isValidSignature := func(payload []byte, signature, key string) bool {
		h := hmac.New(sha256.New, []byte(key))
		h.Write(payload)
		expectedSignature := hex.EncodeToString(h.Sum(nil))
		fmt.Printf("expectedSignature: %s\n", expectedSignature)
		return hmac.Equal([]byte(signature), []byte(expectedSignature))
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	signature := md["x-auth-signature"]
	if len(signature) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Signature header not found")
	}

	signatureKey := config.Env.SignatureKey
	payload, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to marshal request")
	}

	fmt.Printf("payload: %s\n", payload)
	fmt.Printf("signature: %s\n", signature[0])

	if !isValidSignature(payload, signature[0], signatureKey) {
		return nil, status.Error(codes.Unauthenticated, "invalid signature")
	}

	return handler(ctx, req)
}

func Interceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		authInterceptor,
		validationInterceptor,
	}
}
