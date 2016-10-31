package registry

import (
	api "github.com/hooklift/apis/go"
	"github.com/hooklift/lift-registry/server/grpc"
	context "golang.org/x/net/context"
)

// Service implements Lift Registry service.
type Service struct{}

// Search finds Lift plugins in the registry.
func (s *Service) Search(ctx context.Context, r *api.SearchRequest) (*api.SearchResponse, error) {
	return nil, nil
}

// Publish indexes plugin metadata.
func (s *Service) Publish(cxt context.Context, r *api.PublishRequest) (*api.PublishResponse, error) {
	return nil, nil
}

// Register registers service with a given GRPC server.
func Register(ctx context.Context, endpoint grpc.ServiceEndpoint) error {
	// Creates a new service instance.
	service := new(Service)

	// Registers GRPC service.
	api.RegisterRegistryServer(endpoint.GRPCServer, service)

	// Registers HTTP endpoint on GRPC gateway muxer.
	return api.RegisterRegistryHandler(ctx, endpoint.HTTPGWMuxer, endpoint.GRPCClientConn)
}
