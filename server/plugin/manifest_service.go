package plugin

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"

	api "github.com/hooklift/apis/go/lift"
	"github.com/hooklift/lift-registry/server/pkg/grpc"
	"github.com/hooklift/lift-registry/server/pkg/identity"
	"github.com/hooklift/uaa/service/tokens/jwt"
)

// Service implements Lift Registry service.
type Service struct{}

// Search finds Lift plugins in the registry.
func (s *Service) Search(ctx context.Context, r *api.SearchRequest) (*api.SearchResponse, error) {
	res := new(api.SearchResponse)

	matches, err := Search(r.Query, int(r.PageNumber), int(r.ResultPerPage))
	if err != nil {

		return res, err
	}

	// Annoying conversion from domain object to api object.
	for _, m := range matches {
		manifest := new(api.PluginManifest)
		manifest.Author = &api.Author{
			Name:  m.Author.Name,
			Email: m.Author.Email,
		}
		manifest.Description = m.Description
		manifest.Homepage = m.Homepage
		manifest.Name = m.Name
		manifest.Version = m.Version
		manifest.FilesUri = m.FilesURI
		manifest.License = m.License

		publishedAt, err := ptypes.TimestampProto(m.PublishedAt)
		if err != nil {
			glog.Errorf("invalid timestamp received by search index for package %q: %+v", err, m.Name)
			continue
		}
		manifest.PublishedAt = publishedAt

		manifest.Packages = make([]*api.Package, 0)
		for _, p := range m.Packages {
			pkg := new(api.Package)
			pkg.Name = p.Name
			pkg.Algorithm = string(p.Algorithm)
			pkg.Checksum = p.Checksum
			pkg.Arch = string(p.Arch)
			pkg.Os = string(p.OS)

			manifest.Packages = append(manifest.Packages, pkg)
		}
		res.Plugins = append(res.Plugins, manifest)
	}

	return res, nil
}

// Publish indexes plugin metadata.
func (s *Service) Publish(ctx context.Context, r *api.PublishRequest) (*api.PublishResponse, error) {
	token, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	_, okg := token.Scopes["global"]
	_, okw := token.Scopes["write"]

	if !okg && !okw {
		return nil, errors.New("unauthorized")
	}

	manifest := new(Manifest)
	p := r.GetPlugin()

	manifest.Name = p.Name
	manifest.AccountID = token.Subject
	manifest.Author = Author(*p.Author)
	manifest.Description = p.Description
	manifest.Homepage = p.Homepage
	manifest.License = p.License
	manifest.Version = p.Version
	manifest.FilesURI = p.FilesUri

	manifest.Packages = make([]*Package, 0)
	for _, pp := range p.GetPackages() {
		pkg := new(Package)
		pkg.Name = pp.Name
		pkg.Algorithm = Algorithm(pp.Algorithm)
		pkg.Arch = Arch(pp.Arch)
		pkg.Checksum = pp.Checksum
		pkg.OS = OS(pp.Os)

		manifest.Packages = append(manifest.Packages, pkg)
	}

	res := new(api.PublishResponse)
	if err := Publish(manifest); err != nil {
		return nil, err
	}

	return res, nil
}

// Unpublish ...
func (s *Service) Unpublish(ctx context.Context, r *api.UnpublishRequest) (*api.UnpublishResponse, error) {
	token, ok := jwt.FromContext(ctx)
	if !ok && !token.Scopes.Has("global") && !token.Scopes.Has("write") {
		// Convert to grpc error
		return nil, errors.New("unauthorized")
	}

	res := new(api.UnpublishResponse)
	if err := Unpublish(r.Id, token.Subject); err != nil {
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
