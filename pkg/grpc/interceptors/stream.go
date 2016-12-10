package interceptors

import "google.golang.org/grpc"

// DefaultStream is the default stream interceptor handler
var DefaultStream = grpc.StreamServerInterceptor(nil)

// StreamSecurity implements the security filter for stream GRPC calls.
func StreamSecurity(next grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO ...
		if next != nil {
			return next(srv, ss, info, handler)
		}
		return handler(srv, ss)
	}
}

// StreamMetrics reports metrics about grpc streaming calls.
func StreamMetrics(next grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if next != nil {
			next(srv, ss, info, handler)
		}
		// TODO ...
		return handler(srv, ss)
	}
}
