package roundtrip

import (
	"net/http"

	"github.com/Escape-Technologies/repeater/pkg/logger"
	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
)

var Client = &http.Client{}

func protoErr(status int, corr int64) *proto.Response {
	res, err := responseToTransport(&http.Response{
		StatusCode: status,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}, corr)
	if err != nil {
		logger.Error("Error parsing %v response", status)
	}
	return res
}

// HandleRequest handles a request from the repeater server
//
// It MUST always return an HTTP response, even if there is an error
//
//   - Read request or `499 Unparsable request`
//   - Send request or `599 API down`
//   - Parse response or `598 API unresponsive`
func HandleRequest(protoReq *proto.Request) *proto.Response {
	httpReq, err := transportToRequest(protoReq)
	if err != nil {
		logger.Error("Error parsing request : %v", err)
		return protoErr(499, protoReq.Correlation)
	}
	if httpReq.Header.Get("X-Debug") == "true" {
		logger.Debug("Printing debug info for request %v", protoReq.Correlation)
		logger.Debug("Url : %v", protoReq.Url)
		dns(protoReq.Url)
		traceroute(protoReq.Url)
		tls(protoReq.Url)
	}

	logger.Debug("Sending request (%v)", protoReq.Correlation)
	httpRes, err := Client.Do(httpReq)
	if err != nil {
		logger.Error("ERROR sending request : %v", err)
		return protoErr(599, protoReq.Correlation)
	}
	logger.Debug("Received response code %d (%v)", httpRes.StatusCode, protoReq.Correlation)

	protoRes, err := responseToTransport(httpRes, protoReq.Correlation)
	if err != nil {
		logger.Error("Error parsing response : %v", err)
		return protoErr(598, protoReq.Correlation)
	}

	return protoRes
}
