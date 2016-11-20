package identity

import (
	"context"
	"crypto/x509"
	"errors"
	"log"
	"net/url"
	"sync"

	"github.com/golang/glog"
	idapi "github.com/hooklift/apis/go/identity"
	"github.com/hooklift/lift-registry/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	grpcClient *grpc.ClientConn
	once       sync.Once
)

// Connection returns a connection to the Identity service.
// If a connection was previously established, it returns it and skips creating a new one.
func Connection(clientID string) *grpc.ClientConn {
	once.Do(func() {
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(config.TLSCert))
		if !ok {
			log.Fatalf("%+v", errors.New("unable to append certificate to cert pool"))
		}

		clientCreds := credentials.NewClientTLSFromCert(certPool, config.IdentityAddress)
		clientOpts := []grpc.DialOption{
			grpc.WithTransportCredentials(clientCreds),
			// uses default backoff re-connections: https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md
			grpc.WithUserAgent(clientID),
		}

		u, err := url.Parse(config.IdentityAddress)
		if err != nil {
			log.Fatalln("Invalid Identity URL: %#v", err)
		}

		clientConn, err := grpc.Dial(u.Host, clientOpts...)
		if err != nil {
			log.Fatalf("%+v", errors.New("failed connecting to Hooklift Identity service"))
		}
		grpcClient = clientConn

		glog.Infof("Connected to identity service at %q...", config.IdentityAddress)
	})

	return grpcClient
}

type tokenKey struct{}

// FromContext returns a registry connection from context if found.
func FromContext(ctx context.Context) (token *idapi.Token, ok bool) {
	token, ok = ctx.Value(tokenKey{}).(*idapi.Token)
	return
}

// NewContext creates a new context with the GRPC client connection attached.
func NewContext(ctx context.Context, token *idapi.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}
