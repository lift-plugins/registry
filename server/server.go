package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/c4milo/handlers/logger"
	"github.com/golang/glog"
	_ "google.golang.org/grpc/grpclog/glogger"

	"github.com/hooklift/lift-registry/client"
	"github.com/hooklift/lift-registry/server/config"
	"github.com/hooklift/lift-registry/server/domain"
	"github.com/hooklift/lift-registry/server/pkg/grpc"
	"github.com/hooklift/lift-registry/server/pkg/identity"
	"github.com/hooklift/lift-registry/server/service/files"
	"github.com/hooklift/lift-registry/server/service/registry"
)

var (
	// Version string is injected when building the binary from Makefile.
	Version string

	// AppName is also injected during the build process from the Makefile.
	AppName string
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Set("logtostderr", "true")
}

// Initializes plugins database
func initBleve() {
	glog.Infof("Opening Bleve index at %q...", config.IndexFile)

	index, err := bleve.Open(config.IndexFile)
	if err == bleve.ErrorIndexPathDoesNotExist {
		glog.Info("Bleve index does not exist, creating it....")

		index, err = bleve.New(config.IndexFile, bleve.NewIndexMapping())
		if err != nil {
			glog.Fatalf("unable to create Bleve index: %+v", err)
		}
	}

	if err != nil {
		glog.Fatalf("unable to open Bleve index at %q", config.IndexFile)
	}

	domain.Init(index)
}

func main() {
	appName := AppName + "-" + Version
	flag.Parse()

	// Reads configurations values
	config.Read()

	// Initializes Bleve index
	initBleve()

	// Initializes metrics sink
	// sink, _ := metrics.NewStatsiteSink(config.StatsiteAddr)
	// metrics.NewGlobal(metrics.DefaultConfig(AppName), sink)

	// GRPC services
	services := []grpc.ServiceRegisterFn{
		registry.Register,
	}

	// These middlewares are invoked bottom up and order matters.
	handler := client.Handler(http.DefaultServeMux)       // Single Page Application  web UI
	handler = files.Handler(handler)                      // File management API to upload or download packages
	handler = identity.Handler(handler, config.ClientURI) // HTTP security filter, for non gRPC requests
	handler = grpc.Handler(handler, services)             // gRPC services, uses interceptors to verify authorization tokens.
	handler = logger.Handler(handler, logger.AppName(appName))

	tlsKeyPair, err := tls.X509KeyPair([]byte(config.TLSCert), []byte(config.TLSKey))
	if err != nil {
		panic(err)
	}

	address := ":" + config.Port
	srv := &http.Server{
		Addr:    address,
		Handler: handler,
		// This is only for GRPC Gateway HTTP server, since GRPC handles its own transport security.
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsKeyPair},
		},
		// ReadTimeout: 15 * time.Second,
		//WriteTimeout: 15 * time.Second,
	}

	glog.Infof("Starting server at %s", address)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		glog.Fatalf("ListenAndServeTLS: %v", err)
	}
}
