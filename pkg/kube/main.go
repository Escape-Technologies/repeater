package kube

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/Escape-Technologies/repeater/pkg/autoprovisioning"
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

	autoprovisioningRetryInterval = 5 * time.Second
	autoprovisioningRetryCount    = 5
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

func connectAndRun(ctx context.Context, cfg *rest.Config, ap *autoprovisioning.Autoprovisioner, isConnected *atomic.Bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
		for !isConnected.Load() || ctx.Err() != nil {
			time.Sleep(1 * time.Second)
		}
		if ctx.Err() != nil {
			lis.Close()
			return
		}
		err := provisionIntegrationWithRetry(ctx, ap, 0)
		if err != nil {
			logger.Error("Error provisioning integration: %v", err)
		}
		<-ctx.Done()
		lis.Close()
	}()

	logger.Debug("Connecting to k8s API")
	err = srv.ServeOnListener(lis)
	if err != nil {
		return fmt.Errorf("error serving: %w", err)
	}
	return nil
}

func provisionIntegrationWithRetry(ctx context.Context, ap *autoprovisioning.Autoprovisioner, count int) error {
	if ap == nil {
		return nil
	}
	if count > autoprovisioningRetryCount {
		return errors.New("failed to provision integration")
	}
	err := ap.CreateIntegration(ctx)
	if err == nil {
		return nil
	} else {
		logger.Debug("Error provisioning integration: %v", err)
	}
	time.Sleep(autoprovisioningRetryInterval)
	return provisionIntegrationWithRetry(ctx, ap, count+1)
}

func AlwaysConnectAndRun(ctx context.Context, ap *autoprovisioning.Autoprovisioner, isConnected *atomic.Bool) {
	logger.Debug("Checking if the k8s API is available...")

	cfg, err := inferConfig()
	if err != nil {
		logger.Debug("Not connected to k8s API: %s", err.Error())
		return
	}

	logger.Info("Exposing API on http://%s:%d", defaultAddress, defaultPort)

	for {
		err := connectAndRun(ctx, cfg, ap, isConnected)
		if err != nil {
			logger.Error("Error connecting to k8s API: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}
