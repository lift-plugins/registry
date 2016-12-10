package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/hooklift/lift-registry/config"
	"github.com/hooklift/lift-registry/pkg/grpc/interceptors"
)

// ServiceEndpoint represents an endpoint in the GRPC server and HTTP Gateway muxer.
type ServiceEndpoint struct {
	GRPCServer     *grpc.Server
	GRPCClientConn *grpc.ClientConn
	HTTPGWMuxer    *runtime.ServeMux
}

// ServiceRegisterFn defines the type for registering grpc services.
type ServiceRegisterFn func(context.Context, ServiceEndpoint) error

// initGRPCServer initializes the GRPC server.
func initGRPCServer() *grpc.Server {
	cert, err := tls.X509KeyPair([]byte(config.TLSCert), []byte(config.TLSKey))
	if err != nil {
		glog.Fatalf("Unable to load GRPC server transport TLS cert and key")
	}

	serverOpts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(func() grpc.UnaryServerInterceptor {
			interceptor := interceptors.UnarySecurity(interceptors.DefaultUnary, config.ClientURI)
			interceptor = interceptors.UnaryMetrics(interceptor)
			return interceptor
		}()),
		grpc.StreamInterceptor(func() grpc.StreamServerInterceptor {
			interceptor := interceptors.StreamSecurity(interceptors.DefaultStream)
			interceptor = interceptors.StreamMetrics(interceptor)
			return interceptor
		}()),
	}
	return grpc.NewServer(serverOpts...)
}

// initGRPCLocalClient initilizes the GPRC client used by GRPC JSON Gateway
func initGRPCLocalClient() *grpc.ClientConn {
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(config.TLSCert))
	if !ok {
		glog.Fatal("grpc-gw: unable to append server TLS cert to cert pool")
	}

	clientCreds := credentials.NewClientTLSFromCert(certPool, config.PrimaryDomain)
	clientOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(clientCreds),
		//grpc.WithBackoffMaxDelay(1 * time.Second),
		grpc.WithUserAgent("grpc-gw"),
	}

	glog.V(2).Info("grpc-gw: connecting to local GRPC server...")

	address := "localhost:" + config.Port
	clientConn, err := grpc.Dial(address, clientOpts...)
	if err != nil {
		glog.Fatalf("grpc-gw: failed to connect to local server: %v", err)
	}
	return clientConn
}

// Handler responds to GRPC and JSON requests over HTTP2, or forwards the request
// to the next handler.
func Handler(h http.Handler, services []ServiceRegisterFn) http.Handler {
	server := initGRPCServer()
	client := initGRPCLocalClient()

	gwmux := runtime.NewServeMux()
	ctx := context.Background()

	serviceEndpoint := ServiceEndpoint{
		GRPCServer:     server,
		GRPCClientConn: client,
		HTTPGWMuxer:    gwmux,
	}

	glog.V(2).Info("Registering GRPC services...")
	for _, register := range services {
		if err := register(ctx, serviceEndpoint); err != nil {
			glog.Fatalf("failed to register service: %v", err)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if r.ProtoMajor == 2 && strings.Contains(contentType, "application/grpc") {
			server.ServeHTTP(w, r)
			return
		}

		// We need to skip /lib/api.swagger.json so that the browser can load it
		// from the static handler.
		accept := r.Header.Get("Accept")
		if strings.Contains(contentType, "application/json") ||
			strings.Contains(accept, "application/json") &&
				r.URL.Path != "/lib/api.swagger.json" {
			gwmux.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}
