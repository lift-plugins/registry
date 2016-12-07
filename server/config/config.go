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

	IdentityAddress = os.Getenv("IDENTITY_ADDR")
	if IdentityAddress == "" {
		IdentityAddress = "https://localhost:9000"
	}

	// For development purposes, use the following command to regenerate cert:
	// openssl req -new -x509 -key cert-key.pem -out cert.pem -days 1920
	TLSCert = os.Getenv("TLS_CERT")
	if TLSCert == "" {
		log.Printf("[WARN] !!! UNSAFE: No TLS certificate configured, using default TLS certificate !!!")
		TLSCert = `
-----BEGIN CERTIFICATE-----
MIIDUzCCAtmgAwIBAgIJAKTf/aVGhWkYMAkGByqGSM49BAEwgZExCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhOZXcgWW9yazERMA8GA1UEBxMITmV3IFlvcmsxFzAVBgNV
BAoTDkhvb2tsaWZ0LCBJbmMuMRQwEgYDVQQLEwtFbmdpbmVlcmluZzEKMAgGA1UE
AxQBKjEhMB8GCSqGSIb3DQEJARYSY2FtaWxvQGhvb2tsaWZ0LmlvMCAXDTE2MTEx
NDE0NTgwNFoYDzIxMTUwNjA5MTQ1ODA0WjCBkTELMAkGA1UEBhMCVVMxETAPBgNV
BAgTCE5ldyBZb3JrMREwDwYDVQQHEwhOZXcgWW9yazEXMBUGA1UEChMOSG9va2xp
ZnQsIEluYy4xFDASBgNVBAsTC0VuZ2luZWVyaW5nMQowCAYDVQQDFAEqMSEwHwYJ
KoZIhvcNAQkBFhJjYW1pbG9AaG9va2xpZnQuaW8wdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAASH3bmfhqPNDE2YdeBG15Yl13GVWlex0QDCh85koZ3kbKMGdDBqgb5gqgwZ
F1rCCpjff+o3D3JaAMYosACOyHn8lnJOcpryqUkwCklxSQqleLJM4EGSitMm8119
tzYhaCajgfkwgfYwHQYDVR0OBBYEFMNqnVpZOU6jIqWaiHr7AnMXpBwWMIHGBgNV
HSMEgb4wgbuAFMNqnVpZOU6jIqWaiHr7AnMXpBwWoYGXpIGUMIGRMQswCQYDVQQG
EwJVUzERMA8GA1UECBMITmV3IFlvcmsxETAPBgNVBAcTCE5ldyBZb3JrMRcwFQYD
VQQKEw5Ib29rbGlmdCwgSW5jLjEUMBIGA1UECxMLRW5naW5lZXJpbmcxCjAIBgNV
BAMUASoxITAfBgkqhkiG9w0BCQEWEmNhbWlsb0Bob29rbGlmdC5pb4IJAKTf/aVG
hWkYMAwGA1UdEwQFMAMBAf8wCQYHKoZIzj0EAQNpADBmAjEAnvDrqcg7Sl2wK/bH
+98IMGMiYdT1FpSqCT3YyVQeCPELlxmnXbzNesY/R+l8oY9bAjEAhya4BL+ingli
o9FuJqdUS5o9Rgii55nFhNdzQvT/p/ANGHBCfQyUNtAjPp92KvXC
-----END CERTIFICATE-----
		`
	}

	// For development purposes, use the following to regenerate key:
	// openssl ecparam -genkey -name secp384r1 -out cert-key.pem
	TLSKey = os.Getenv("TLS_KEY")
	if TLSKey == "" {
		log.Printf("[WARN] !!! UNSAFE: No TLS private key configured, using default TLS private key !!!")
		TLSKey = `
-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDAD5WLfofxxT8EewIU/VYJ5hXRWAyjxwAhemboTVnnrCmqA5Icxz+oa
kVluFU7LiPWgBwYFK4EEACKhZANiAASH3bmfhqPNDE2YdeBG15Yl13GVWlex0QDC
h85koZ3kbKMGdDBqgb5gqgwZF1rCCpjff+o3D3JaAMYosACOyHn8lnJOcpryqUkw
CklxSQqleLJM4EGSitMm8119tzYhaCY=
-----END EC PRIVATE KEY-----
		`
	}

	ClientURI = os.Getenv("OAUTH2_CLIENT_URI")
	if ClientURI == "" {
		ClientURI = "https://lift.hooklift.io"
	}
}
