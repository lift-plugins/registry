package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/c4milo/handlers/logger"
	"github.com/golang/glog"
	_ "google.golang.org/grpc/grpclog/glogger"

	"github.com/hooklift/lift-registry/client"
	"github.com/hooklift/lift-registry/server/config"
	"github.com/hooklift/lift-registry/server/grpc"
	"github.com/hooklift/lift-registry/server/web/fileupload"
	"github.com/hooklift/lift-registry/server/web/registry"
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
func initDatabase() {
	// mapping := bleve.NewIndexMapping()
	// index, err := bleve.New("lift-plugins.bleve", mapping)
	// if err != nil {
	// 	panic(err)
	// }
}

func main() {
	appName := AppName + "-" + Version
	flag.Parse()

	// Reads configurations values
	config.Read()

	// Initializes database
	initDatabase()

	// Initializes metrics sink
	// sink, _ := metrics.NewStatsiteSink(config.StatsiteAddr)
	// metrics.NewGlobal(metrics.DefaultConfig(AppName), sink)

	// GRPC services
	services := []grpc.ServiceRegisterFn{
		registry.Register,
	}
	mux := http.DefaultServeMux

	// These middlewares are invoked bottom up and order matters.
	rack := client.Handler(mux)
	rack = fileupload.Handler(rack, new(fileupload.S3))
	rack = grpc.Handler(rack, services)
	rack = logger.Handler(rack, logger.AppName(appName))

	tlsKeyPair, err := tls.X509KeyPair([]byte(config.TLSCert), []byte(config.TLSKey))
	if err != nil {
		panic(err)
	}

	address := ":" + config.Port
	srv := &http.Server{
		Addr:    address,
		Handler: rack,
		// This is only for GRPC Gateway HTTP server, since GRPC handles its own transport security.
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsKeyPair},
		},
		//ReadTimeout:  5 * time.Second,
		//WriteTimeout: 15 * time.Second,
	}

	glog.Infof("Starting server at %s", address)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		glog.Fatalf("ListenAndServeTLS: %v", err)
	}
}
