package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC interceptor that logs, tracks request time,
// and optionally supports request ID or future auth.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Extract request ID (if provided via metadata)
		md, ok := metadata.FromIncomingContext(ctx)
		var requestID string
		if ok {
			if val := md.Get("x-request-id"); len(val) > 0 {
				requestID = val[0]
			}
		}

		// Log request method
		log.Printf("📥 [%s] incoming gRPC: %s", requestID, info.FullMethod)

		// 👇 You can inject auth here (stubbed for now)
		// user, err := authenticate(md)
		// if err != nil {
		//     return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
		// }

		// Call the actual handler
		resp, err := handler(ctx, req)

		// Log result
		duration := time.Since(start)
		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("❌ [%s] %s failed (%s): %s", requestID, info.FullMethod, duration, st.Message())
		} else {
			log.Printf("✅ [%s] %s completed in %s", requestID, info.FullMethod, duration)
		}

		return resp, err
	}
}
