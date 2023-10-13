package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/Escape-Technologies/repeater/pkg/roundtrip"
	"github.com/Escape-Technologies/repeater/pkg/stream"
)

// Injected by ldflags
var (
	version = "dev"
	commit  = "none"
)

var UUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func main() {
	logger.Info("Running Escape repeater version %s, commit %s\n", version, commit)

	repeaterId := os.Getenv("ESCAPE_REPEATER_ID")
	if !UUID.MatchString(repeaterId) {
		logger.Error("ESCAPE_REPEATER_ID must be a UUID in lowercase")
		logger.Error("To get your repeater id, go to https://app.escape.tech/organization/network/")
		logger.Error("For more information, read the docs at https://docs.escape.tech/enterprise/repeater")
		os.Exit(1)
	}

	mTLScrt := os.Getenv("ESCAPE_REPEATER_mTLS_CRT_FILE")
	mTLSkey := os.Getenv("ESCAPE_REPEATER_mTLS_KEY_FILE")
	if mTLScrt != "" && mTLSkey == "" {
		log.Println("ESCAPE_REPEATER_mTLS_KEY_FILE must be set if ESCAPE_REPEATER_mTLS_CRT_FILE is set")
		os.Exit(1)
	}
	if mTLScrt == "" && mTLSkey != "" {
		log.Println("ESCAPE_REPEATER_mTLS_CRT_FILE must be set if ESCAPE_REPEATER_mTLS_KEY_FILE is set")
		os.Exit(1)
	}
	if mTLScrt != "" && mTLSkey != "" {
		log.Printf("Using mTLS from files : %s, %s\n", mTLScrt, mTLSkey)

		cert, err := tls.LoadX509KeyPair(mTLScrt, mTLSkey)
		if err != nil {
			log.Fatalf("Error loading mTLS keypair: %v\n", err)
			os.Exit(1)
		}
		internal.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		}
	}

	start(repeaterId)
}

func getCon() *grpc.ClientConn {
	url := os.Getenv("ESCAPE_REPEATER_URL")
	if url == "" {
		url = "repeater.escape.tech:443"
	} else {
		logger.Debug("Using custom repeater url: %s\n", url)
	}

	insecure := os.Getenv("ESCAPE_REPEATER_INSECURE")
	if insecure == "1" || insecure == "true" {
		logger.Debug("Allowing insecure ssl connections")
		if roundtrip.Client.Transport == nil {
			roundtrip.Client.Transport = http.DefaultTransport
		}
		roundtrip.Client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	logger.Info("Starting repeater client...")

	go logger.AlwaysConnect(url, repeaterId)
	stream.AlwaysConnectAndRun(url, repeaterId)
}
