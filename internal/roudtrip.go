package internal

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Escape-Technologies/repeater/proto"
)

func newErrorResponse(r *http.Request, status int, body string) *http.Response {
	resp := &http.Response{}
	resp.Request = r
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", "application/json")
	resp.StatusCode = status
	buf := bytes.NewBufferString(body)
	resp.ContentLength = int64(buf.Len())
	resp.Body = io.NopCloser(buf)
	return resp
}

func apiDownResponse(r *http.Request) *http.Response {
	return newErrorResponse(r, 599, `{"error": "API is down"}`)
}

func apiUnresponsiveResponse(r *http.Request) *http.Response {
	return newErrorResponse(r, 598, `{"error": "API is unresponsive"}`)
}

func RoundTrip(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ERROR sending request")
		return apiDownResponse(req)
	}
	return res
}

func HandleRequest(req *proto.Request, stream *proto.Repeater_StreamClient) {
	startTime := time.Now()
	httpReq, err := transportToRequest(req)
	if err != nil {
		log.Printf("ERROR[internal.TransportToRequest]: %v\n", err)
		return
	}

	// work
	httpRes := RoundTrip(httpReq)

	tRes, err := responseToTransport(httpRes, req.Correlation)
	if err != nil {
		log.Printf("ERROR[internal.ResponseToTransport]: %v\n", err)
		return
	}

	err = (*stream).Send(tRes)
	if err != nil {
		log.Printf("ERROR[stream.Send]: %v\n", err)
		return
	}
	log.Println("Ok in", time.Since(startTime))
}
