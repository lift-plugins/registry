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
	// ClientID is the OAuth2 client application identifier assigned by Hooklift Identity
	ClientID string
	// ClientSecret is required in order to be able to refresh access tokens.
	ClientSecret string
	// S3Bucket is the bucket where all published plugin packages are going to be stored.
	S3Bucket string
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

	// For development purposes, use the following command to regenerate cert:
	// openssl req -new -x509 -key cert-key.pem -out cert.pem -days 1920
	TLSCert = os.Getenv("TLS_CERT")
	if TLSCert == "" {
		log.Printf("[WARN] !!! UNSAFE: No TLS certificate configured, using default TLS certificate !!!")
		TLSCert = `
-----BEGIN CERTIFICATE-----
MIIDUTCCAtegAwIBAgIJAOWuQSsLeG+zMAkGByqGSM49BAEwgZExCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhOZXcgWW9yazERMA8GA1UEBxMITmV3IFlvcmsxFzAVBgNV
BAoTDkhvb2tsaWZ0LCBJbmMuMRQwEgYDVQQLEwtFbmdpbmVlcmluZzEKMAgGA1UE
AxQBKjEhMB8GCSqGSIb3DQEJARYSY2FtaWxvQGhvb2tsaWZ0LmlvMB4XDTE2MDgx
NTIwMjQyNVoXDTE2MTExMzIwMjQyNVowgZExCzAJBgNVBAYTAlVTMREwDwYDVQQI
EwhOZXcgWW9yazERMA8GA1UEBxMITmV3IFlvcmsxFzAVBgNVBAoTDkhvb2tsaWZ0
LCBJbmMuMRQwEgYDVQQLEwtFbmdpbmVlcmluZzEKMAgGA1UEAxQBKjEhMB8GCSqG
SIb3DQEJARYSY2FtaWxvQGhvb2tsaWZ0LmlvMHYwEAYHKoZIzj0CAQYFK4EEACID
YgAELp07EO1MczG950aucWp3qxo5FVT+9BZL5iDJiE31FkbqGZuFf7gOwB7kmeGW
1x+Ws7KGPaPgYKDaUHseuJkS+3+hguw4BY6eBcCbU1YDKS0bIgU6F5p2tiXbSBpC
K1GSo4H5MIH2MB0GA1UdDgQWBBRfGlwOi5nklQEonzq66YnoWp+yTzCBxgYDVR0j
BIG+MIG7gBRfGlwOi5nklQEonzq66YnoWp+yT6GBl6SBlDCBkTELMAkGA1UEBhMC
VVMxETAPBgNVBAgTCE5ldyBZb3JrMREwDwYDVQQHEwhOZXcgWW9yazEXMBUGA1UE
ChMOSG9va2xpZnQsIEluYy4xFDASBgNVBAsTC0VuZ2luZWVyaW5nMQowCAYDVQQD
FAEqMSEwHwYJKoZIhvcNAQkBFhJjYW1pbG9AaG9va2xpZnQuaW+CCQDlrkErC3hv
szAMBgNVHRMEBTADAQH/MAkGByqGSM49BAEDaQAwZgIxALaiHMepDgC+s/YOppjh
2Nj7ZVhRsyZXXirdBRv9WPJNr63ZVLc/ZknPtUCowr9IvgIxAL61ltwoDHcGRUj2
YwpZ+1QyNNCekHodFohHj/jKwcHebgPrGABvs86bStKpT4ThuQ==
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
MIGkAgEBBDDC0AdvtdeP7wL9qSCpDfBExh+1j5qKIaWTl0LfDD/Pkih28h5czN+A
MiOp70EtT7mgBwYFK4EEACKhZANiAAQunTsQ7UxzMb3nRq5xanerGjkVVP70Fkvm
IMmITfUWRuoZm4V/uA7AHuSZ4ZbXH5azsoY9o+BgoNpQex64mRL7f6GC7DgFjp4F
wJtTVgMpLRsiBToXmna2JdtIGkIrUZI=
-----END EC PRIVATE KEY-----
		`
	}
}
