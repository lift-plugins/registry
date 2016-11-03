package grpc

import (
	"errors"
	"strings"

	"github.com/golang/glog"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrInvalidTokenType is returned if a Bearer prefix is not found in the authorization header
	ErrInvalidTokenType = errors.New("security:invalid-token-type")
)

// securityUnaryInterceptor checks whether an access token was provided and decodes it. Otherwise, it yields to the next interceptor or the service function called.
func securityUnary(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	yield := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if next != nil {
			return next(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		glog.V(4).Infof("GRPC request: %#v", req)

		md, ok := metadata.FromContext(ctx)
		if !ok {
			glog.Warning("No metadata found")
			return yield(ctx, req, info, handler)
		}

		values, ok := md["authorization"]
		if !ok {
			glog.Warning("No authorization metadata found")
			return yield(ctx, req, info, handler)
		}

		glog.Info("Decoding and verifying JWT token...")
		authValue := values[0]
		glog.Infof("Token: %s", authValue)

		tokenParts := strings.Fields(authValue)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return nil, errors.New("Unauthorized")
		}

		// TODO(c4milo): Reach out to authorization server via gRPC to verify token, we have to do it this way
		// in order to properly support token or grant revocations.
		// token, err := jwt.Decode(tokenParts[1])
		// if err != nil {
		// 	return nil, errors.New("Unauthorized")
		// }

		//ctx = context.WithValue(ctx, jwt.TokenCtxKey, token)
		return yield(ctx, req, info, handler)
	}
}

// securityStreamInterceptor implements the security filter for stream GRPC calls.
func securityStream(next grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO ...
		if next != nil {
			return next(srv, ss, info, handler)
		}
		return handler(srv, ss)
	}
}

// metricsUnary reports metrics about grpc unary calls.
func metricsUnary(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// TODO(c4milo): Report whatever metrics we define as valuable to know.
		if next != nil {
			return next(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}
}

// metricsStream reports metrics about grpc streaming calls.
func metricsStream(next grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if next != nil {
			next(srv, ss, info, handler)
		}
		// TODO ...
		return handler(srv, ss)
	}
}

// unaryInterceptor chains all unary interceptors used.
func unaryInterceptors() grpc.UnaryServerInterceptor {
	rack := securityUnary(nil)
	rack = metricsUnary(rack)
	return rack
}

// StreamInterceptor chains all stream interceptors used.
func streamInterceptors() grpc.StreamServerInterceptor {
	rack := securityStream(nil)
	rack = metricsStream(rack)
	return rack
}
