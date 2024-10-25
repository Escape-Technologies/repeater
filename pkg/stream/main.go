package stream

import (
	"sync/atomic"
	"time"

	"github.com/Escape-Technologies/repeater/pkg/grpc"
	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/Escape-Technologies/repeater/pkg/roundtrip"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
)

func AlwaysConnectAndRun(url, repeaterId string, isConnected *atomic.Bool) {
	for {
		hasConnected := ConnectAndRun(url, repeaterId, isConnected)
		isConnected.Store(false)
		logger.Info("Disconnected...")
		if !hasConnected {
			logger.Info("Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}

func ConnectAndRun(url, repeaterId string, isConnected *atomic.Bool) (hasConnected bool) {
	// Set up proxy dialer
	proxyURL, _ := http.ProxyFromEnvironment(&http.Request{URL: &url.URL{Scheme: "https"}})
	var dialer func(context.Context, string) (net.Conn, error)  = nil
	if proxyURL != nil {
		proxyDialer := &net.Dialer{}
		dialer = func(ctx context.Context, addr string) (net.Conn, error) {
			return proxyDialer.DialContext(ctx, "tcp", proxyURL.Host)
		}
	}

	// Set up gRPC connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if dialer != nil {
		opts = append(opts, grpc.WithContextDialer(dialer))
	}

	// Create gRPC stream
	stream, closer, err := grpc.Stream(url, repeaterId)
	defer closer()
	if err != nil {
		for _, why := range extractWhyStreamCreateError(err) {
			logger.Error(why) //nolint:govet
		}
		return false
	}
	logger.Info("Repeater connected to server...")
	isConnected.Store(true)

	// Send healthcheck to the server
	go func() {
		retries := 0

		for {
			logger.Debug("Sending healthcheck...")
			err = stream.Send(&proto.Response{
				Data:        []byte(""),
				Correlation: 0,
			})
			if err != nil {
				logger.Error("Error sending healthcheck: %v", err)
				retries++
				if retries > 5 {
					logger.Warn("Too many retries, stopping healthchecks...")
					return
				}
			} else {
				retries = 0
			}
			logger.Debug("Healthcheck sent")
			time.Sleep(30 * time.Second)
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			for _, why := range extractWhyRecvError(err) {
				logger.Error(why) //nolint:govet
			}
			return true
		}
		logger.Info("Received incoming stream (%d)", req.Correlation)

		// Send request to server
		// Use a go func to avoid blocking the stream
		go func() {
			startTime := time.Now()
			res := roundtrip.HandleRequest(req)
			logger.Info("Processed stream in %v (%d)", time.Since(startTime), req.Correlation)

			err = stream.Send(res)
			if err != nil {
				logger.Error("Error processing stream (%d): %v", req.Correlation, err)
			}
		}()
	}
}
