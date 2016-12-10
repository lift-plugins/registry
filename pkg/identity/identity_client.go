package identity

import (
	"context"
	"crypto/x509"
	"log"
	"net/url"
	"sync"

	"github.com/golang/glog"
	idapi "github.com/hooklift/apis/go/identity"
	"github.com/hooklift/lift-registry/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	grpcClient *grpc.ClientConn
	once       sync.Once
	// tlsCert is defined when compiling with build tag "dev"
	tlsCert string
)

// Connection returns a connection to the Identity service.
// If a connection was previously established, it returns it and skips creating a new one.
func Connection(clientURI string) *grpc.ClientConn {
	once.Do(func() {
		clientOpts := []grpc.DialOption{
			// uses default backoff re-connections: https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md
			grpc.WithUserAgent(clientURI),
		}

		if tlsCert != "" {
			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM([]byte(tlsCert))
			if !ok {
				log.Fatalf("%+v", errors.New("unable to append certificate to cert pool"))
			}

			clientCreds := credentials.NewClientTLSFromCert(certPool, config.IdentityAddress)
			clientOpts = append(clientOpts, grpc.WithTransportCredentials(clientCreds))
		}

		u, err := url.Parse(config.IdentityAddress)
		if err != nil {
			log.Fatalln("Invalid Identity URL: %#v", err)
		}

		clientConn, err := grpc.Dial(u.Host, clientOpts...)
		if err != nil {
			log.Fatalf("%+v", errors.Wrapf(err, "failed connecting to Hooklift Identity service"))
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
