package interceptors

import (
	"errors"
	"strings"

	"github.com/golang/glog"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	idapi "github.com/hooklift/apis/go/identity"
	"github.com/hooklift/lift-registry/server/config"
	"github.com/hooklift/lift-registry/server/pkg/identity"
)

// DefaultUnary is the default unary interceptor handler
var DefaultUnary = grpc.UnaryServerInterceptor(nil)

// UnarySecurity checks whether an access token was provided and decodes it. Otherwise, it yields to the next interceptor or the service function called.
func UnarySecurity(next grpc.UnaryServerInterceptor, clientID string) grpc.UnaryServerInterceptor {
	yieldFn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		if next != nil {
			return next(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}

	accounts := idapi.NewAccountsClient(identity.Connection(config.ClientID))

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, ok := metadata.FromContext(ctx)
		if !ok {
			glog.Warning("No metadata found")
			return yieldFn(ctx, req, info, handler)
		}

		values, ok := md["authorization"]
		if !ok {
			glog.Warning("No authorization header found")
			return yieldFn(ctx, req, info, handler)
		}

		authValue := values[0]

		tokenParts := strings.Fields(authValue)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return nil, errors.New("Unauthorized")
		}

		tokenValue := tokenParts[1]

		res, err := accounts.VerifyToken(ctx, &idapi.VerifyTokenRequest{
			ClientId: clientID,
			Token:    tokenValue,
		})

		if err != nil {
			return nil, errors.New("Unauthorized")
		}

		ctx = identity.NewContext(ctx, res.Token)
		return yieldFn(ctx, req, info, handler)
	}
}

// UnaryMetrics reports metrics about grpc unary calls.
func UnaryMetrics(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// TODO(c4milo): Report whatever metrics we define as valuable to know.
		if next != nil {
			return next(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}
}
