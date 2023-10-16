package main

import (
	"os"
	"regexp"
	"time"

	"github.com/Escape-Technologies/repeater/pkg/grpc"
	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/Escape-Technologies/repeater/pkg/roundtrip"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
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

	url := os.Getenv("ESCAPE_REPEATER_URL")
	if url == "" {
		url = "repeater.escape.tech:443"
	} else {
		logger.Debug("Using custom repeater url: %s\n", url)
	}

	logger.Info("Starting repeater client...")

	go func() {
		for {
			if !logger.ConnectLogs(url, repeaterId) {
				time.Sleep(time.Second)
			}
		}
	}()

	for {
		hasConnected := connectAndRun(url, repeaterId)
		logger.Info("Disconnected...")
		if !hasConnected {
			logger.Info("Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}

func connectAndRun(url, repeaterId string) (hasConnected bool) {
	stream, closer, err := grpc.Stream(url, repeaterId)
	defer closer()
	if err != nil {
		logger.Error("Error creating stream: %v", err)
		return false
	}
	logger.Info("Connected to server...")

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
			logger.Error("Error receiving: %v", err)
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
