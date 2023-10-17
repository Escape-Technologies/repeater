package debug

import (
	"time"

	"github.com/Escape-Technologies/repeater/pkg/grpc"
	"github.com/Escape-Technologies/repeater/pkg/logger"
	"github.com/Escape-Technologies/repeater/pkg/roundtrip"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
)

func AlwaysConnectAndRun(url, repeaterId string) {
	for {
		hasConnected := ConnectAndRun(url, repeaterId)
		logger.Info("Disconnected debug stream...")
		if !hasConnected {
			logger.Info("Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}

func ConnectAndRun(url, repeaterId string) (hasConnected bool) {
	stream, closer, err := grpc.Debug(url, repeaterId)
	defer closer()
	if err != nil {
		logger.Error("Error creating stream: %v", err)
		return false
	}
	logger.Info("Connected to debug stream...")

	for {
		req, err := stream.Recv()
		if err != nil {
			logger.Error("Error receiving: %v", err)
			return true
		}
		logger.Info("Received incoming debug stream (%d)", req.Correlation)

		// Send request to server
		// Use a go func to avoid blocking the stream
		go func() {
			startTime := time.Now()
			res := roundtrip.HandleRequest(req)

			err = stream.Send(&proto.DebugResponse{
				Response: res,
				Infos: &proto.DebugInfo{
					Dns:        dns(req.Url),
					Traceroute: traceroute(req.Url),
				},
				Correlation: req.Correlation,
			})
			logger.Info("Processed debug stream in %v (%d)", time.Since(startTime), req.Correlation)
			if err != nil {
				logger.Error("Error processing stream (%d): %v", req.Correlation, err)
			}
		}()
	}
}
