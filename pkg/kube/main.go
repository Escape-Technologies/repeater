package kube

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Escape-Technologies/repeater/pkg/logger"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/proxy"
)

const (
	defaultPort         = 8001
	defaultStaticPrefix = "/static/"
	defaultAPIPrefix    = "/"
	defaultAddress      = "127.0.0.1"
)

func inferConfig() (*rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig != "" {
		logger.Debug("Using kubeconfig : %s", kubeconfig)
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		logger.Debug("Using in cluster config")
		return rest.InClusterConfig()
	}
}

func connectAndRun(cfg *rest.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger.Debug("Connected to k8s API")

	srv, err := proxy.NewServer(
		"",
		"/",
		defaultStaticPrefix,
		nil,
		cfg,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("error creating proxy server: %w", err)
	}

	lis, err := srv.Listen(defaultAddress, defaultPort)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}

	go func() {
		<-ctx.Done()
		lis.Close()
	}()

	err = srv.ServeOnListener(lis)
	if err != nil {
		return fmt.Errorf("error serving: %w", err)
	}
	return nil
}

func AlwaysConnectAndRun() {
	logger.Debug("Checking if the k8s API is available...")

	cfg, err := inferConfig()
	if err != nil {
		logger.Debug("Not connected to k8s API")
		return
	}

	logger.Info("Exposing API on http://%s:%d", defaultAddress, defaultPort)

	for {
		err := connectAndRun(cfg)
		if err != nil {
			logger.Error("Error connecting to k8s API: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}
