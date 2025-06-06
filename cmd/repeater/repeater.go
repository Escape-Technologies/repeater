package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"

	"github.com/Escape-Technologies/repeater/pkg/autoprovisioning"
	"github.com/Escape-Technologies/repeater/pkg/health"
	"github.com/Escape-Technologies/repeater/pkg/kube"
	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/Escape-Technologies/repeater/pkg/roundtrip"
	"github.com/Escape-Technologies/repeater/pkg/stream"

	_ "net/http/pprof"
)

// Injected by ldflags
var (
	version = "dev"
	commit  = "none"
)

var UUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var PORT = regexp.MustCompile(`^[1-9][0-9]{0,5}$`)

func setupHTTPClients() string {
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
		roundtrip.MTLSClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		}
	}

	url := os.Getenv("ESCAPE_REPEATER_URL")
	if url == "" {
		url = "repeater.escape.tech:443"
	} else {
		logger.Debug("Using custom repeater url: %s", url)
	}

	insecure := os.Getenv("ESCAPE_REPEATER_INSECURE")
	if insecure == "1" || insecure == "true" {
		if mTLScrt != "" && mTLSkey != "" {
			logger.Warn("Insecure SSL flag is enabled, so mTLS will not be used.")
			roundtrip.MTLSClient = nil
		}

		logger.Debug("Allowing insecure ssl connections")
		roundtrip.DefaultClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}

	disableRedirectsEnv := os.Getenv("ESCAPE_REPEATER_DISABLE_REDIRECTS")
	if disableRedirectsEnv == "1" || disableRedirectsEnv == "true" {
		roundtrip.DisableRedirects = true
	}

	return url
}

func setupProxyURL() string {
	proxyURL := os.Getenv("ESCAPE_REPEATER_PROXY_URL")
	if proxyURL != "" {
		logger.Debug("Using custom proxy url: %s", proxyURL)
	}
	return proxyURL
}

var isConnected = &atomic.Bool{}

func startHealth() {
	isConnected.Store(false)
	healthCheckPort := os.Getenv("HEALTH_CHECK_PORT")
	if healthCheckPort != "" {
		if !PORT.MatchString(healthCheckPort) {
			logger.Error("HEALTH_CHECK_PORT must be a valid port number, falling back to no health check")
		} else {
			health.HealthCheck(healthCheckPort, isConnected)
		}
	}
}

func getRepeaterId(ctx context.Context, ap *autoprovisioning.Autoprovisioner) string {
	repeaterId := os.Getenv("ESCAPE_REPEATER_ID")
	if ap != nil && repeaterId == "" {
		logger.Info("ESCAPE_REPEATER_ID is not set, using autoprovisioning")
		repeaterId, err := ap.GetId(ctx)
		if err != nil {
			logger.Error("Failed to get repeater id from autoprovisioning, %s", err.Error())
			os.Exit(1)
		}
		return repeaterId
	}
	if !UUID.MatchString(repeaterId) {
		logger.Error("ESCAPE_REPEATER_ID must be a UUID in lowercase")
		logger.Error("To get your repeater id, go to https://app.escape.tech/organization/network/")
		logger.Error("For more information, read the docs at https://docs.escape.tech/enterprise/repeater")
		os.Exit(1)
	}
	return repeaterId
}

func pprof() {
	logger.Info("Starting pprof on http://0.0.0.0:6060/debug/pprof/")
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		logger.Info("Failed to start pprof server, %s", err.Error())
	} else {
		logger.Info("Started pprof on http://0.0.0.0:6060/debug/pprof/")
	}
}

func getAutoprovisioner() *autoprovisioning.Autoprovisioner {
	ap, err := autoprovisioning.NewAutoprovisioner()
	if err != nil {
		logger.Info("Unable to setup autoprovisioning, using only local repeater id: %s", err.Error())
		return nil
	}
	return ap
}

func main() {
	ctx := context.Background()
	logger.Info("Running Escape repeater version %s, commit %s", version, commit)
	go pprof()
	go startHealth()

	proxyURL := setupProxyURL()
	ap := getAutoprovisioner()
	repeaterId := getRepeaterId(ctx, ap)
	url := setupHTTPClients()

	logger.Info("Starting repeater client...")

	go logger.AlwaysConnect(url, repeaterId, proxyURL)
	go kube.AlwaysConnectAndRun(ctx, ap, isConnected)
	stream.AlwaysConnectAndRun(url, repeaterId, isConnected, proxyURL)
}
