package internal

import (
	"log"
	"net/http"

	"github.com/Escape-Technologies/repeater/proto"
)

func protoErr(status int, corr int64) *proto.Response {
	res, err := responseToTransport(&http.Response{
		StatusCode: status,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}, corr)
	if err != nil {
		log.Fatalf("Error parsing %v response", status)
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
		log.Printf("Error parsing request : %v\n", err)
		return protoErr(499, protoReq.Correlation)
	}

	// work
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("ERROR sending request : %v\n", err)
		return protoErr(599, protoReq.Correlation)
	}

	protoRes, err := responseToTransport(httpRes, protoReq.Correlation)
	if err != nil {
		log.Printf("Error parsing response : %v\n", err)
		return protoErr(598, protoReq.Correlation)
	}

	return protoRes
}
