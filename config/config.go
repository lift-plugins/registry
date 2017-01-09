package config

import (
	"log"
	"os"
)

var (
	// PrimaryDomain is the main domain name by which the application is accessed.
	PrimaryDomain string
	// Port is the TCP port on which the service will accept connections.
	Port string
	// TLSCert is the PEM encoded value of the TLS certificate
	TLSCert string
	// TLSKey is the PEM encoded value of the TLS private key used to generate the certificate
	TLSKey string
	// ClientURI is the OAuth2 client application URI as registered in Hooklift Identity Provider
	ClientURI string
	// ClientSecret is required in order to be able to refresh access tokens.
	ClientSecret string
	// S3Bucket is the bucket where all published plugin packages are going to be stored.
	S3Bucket string
	// IndexFile contains the path to the database file where we store everything that is published.
	IndexFile string
	// IdentityAddress is the address to the identity server used to verify user's tokens.
	IdentityAddress string
)

// Read loads the configuration values.
func Read() {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		log.Fatal("AWS_ACCESS_KEY_ID with permissions to send emails using Amazon SES is required")
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatal("AWS_SECRET_ACCESS_KEY with permissions to send emails using Amazon SES is required")
	}

	if os.Getenv("AWS_REGION") == "" {
		os.Setenv("AWS_REGION", "us-east-1")
	}

	if os.Getenv("S3_BUCKET") == "" {
		S3Bucket = "hooklift-lift-registry"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}
	Port = port

	PrimaryDomain = os.Getenv("PRIMARY_DOMAIN")
	if PrimaryDomain == "" {
		PrimaryDomain = "localhost:" + Port
	}

	IndexFile = os.Getenv("INDEX_FILE")
	if IndexFile == "" {
		IndexFile = "tmp/registry.bleve"
	}

	// For development purposes, use the following command to regenerate cert:
	// openssl req -new -x509 -key cert-key.pem -out cert.pem -days 1920
	TLSCert = os.Getenv("TLS_CERT")
	if TLSCert == "" {
		log.Fatal("TLS_CERT config variable must be set")
	}

	// For development purposes, use the following to regenerate key:
	// openssl ecparam -genkey -name secp384r1 -out cert-key.pem
	TLSKey = os.Getenv("TLS_KEY")
	if TLSKey == "" {
		log.Fatal("TLS_KEY config variable must be set")
	}

	ClientURI = os.Getenv("OAUTH2_CLIENT_URI")
	if ClientURI == "" {
		ClientURI = "https://lift.hooklift.io"
	}
}
