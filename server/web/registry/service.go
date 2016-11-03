package registry

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	context "golang.org/x/net/context"

	api "github.com/hooklift/apis/go/lift-registry"
	"github.com/hooklift/lift-registry/server/domain/plugin"
	"github.com/hooklift/lift-registry/server/pkg/grpc"
)

// Service implements Lift Registry service.
type Service struct{}

// Search finds Lift plugins in the registry.
func (s *Service) Search(ctx context.Context, r *api.SearchRequest) (*api.SearchResponse, error) {
	matches, err := plugin.Search(r.Query, int(r.PageNumber), int(r.ResultPerPage))
	if err != nil {
		return nil, err
	}

	res := new(api.SearchResponse)

	// Annoying conversion from domain object to api object.
	for _, m := range matches {
		manifest := new(api.PluginManifest)
		manifest.Author = (*api.Author)(&m.Author)
		manifest.Description = m.Description
		manifest.Homepage = m.Homepage
		manifest.Name = m.Name
		manifest.Version = (*api.Version)(&m.Version)

		publishedAt, err := ptypes.TimestampProto(m.PublishedAt)
		if err != nil {
			glog.Errorf("invalid timestamp received by search index for package %q: %+v", err, m.Name)
			continue
		}
		manifest.PublishedAt = publishedAt

		manifest.Packages = make([]*api.Package, len(m.Packages))
		for _, p := range m.Packages {
			pkg := new(api.Package)
			pkg.Algorithm = api.Algorithm(api.Algorithm_value[string(p.Algorithm)])
			pkg.Checksum = p.Checksum
			pkg.Arch = api.Arch(api.Arch_value[string(p.Arch)])
			pkg.Os = api.OS(api.OS_value[string(p.OS)])
			pkg.Url = p.URL

			manifest.Packages = append(manifest.Packages, pkg)
		}
		res.Plugins = append(res.Plugins, manifest)
	}

	return res, nil
}

// Publish indexes plugin metadata.
func (s *Service) Publish(ctx context.Context, r *api.PublishRequest) (*api.PublishResponse, error) {
	// TODO(c4milo): get token from context
	// TODO(c4milo): Verify signature
	// TODO(c4milo): Verify expiration
	// TODO(c4milo): validate that token has access to lift registry and scope for publishing plugins
	// TODO(c4milo): Return unauthorized error if not
	manifest := new(plugin.Manifest)
	p := r.GetPlugin()

	manifest.Name = p.Name
	manifest.Author = plugin.Author(*p.Author)
	manifest.Description = p.Description
	manifest.Homepage = p.Homepage
	manifest.License = p.License
	manifest.Version = plugin.Version(*p.Version)

	manifest.Packages = make([]*plugin.Package, len(p.Packages))
	for _, pp := range p.GetPackages() {
		pkg := new(plugin.Package)
		pkg.Algorithm = plugin.Algorithm(pp.Algorithm)
		pkg.Arch = plugin.Arch(pp.Arch)
		pkg.Checksum = pp.Checksum
		pkg.OS = plugin.OS(pp.Os)
		pkg.URL = pp.Url

		manifest.Packages = append(manifest.Packages, pkg)
	}

	if err := plugin.Publish(manifest); err != nil {
		return nil, err
	}

	res := new(api.PublishResponse)
	return res, nil
}

// Unpublish ...
func (s *Service) Unpublish(ctx context.Context, r *api.UnpublishRequest) (*api.UnpublishResponse, error) {
	res := new(api.UnpublishResponse)
	if err := plugin.Unpublish(r.Id); err != nil {
		return nil, err
	}
	return res, nil
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
