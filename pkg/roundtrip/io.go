package roundtrip

import (
	"bufio"
	"net/http"
	"net/url"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
)

func responseToTransport(r *http.Response, correlation int64) (*proto.Response, error) { // In the other program
	du := newDump()
	err := r.Write(du)
	if err != nil {
		return nil, err
	}
	return &proto.Response{
		Data:        du.data,
		Correlation: correlation,
	}, nil
}

func transportToRequest(r *proto.Request) (*http.Request, error) { // In the other program
	lo := newLoad(r.Data)
	req, err := http.ReadRequest(bufio.NewReader(lo))
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(r.Url)
	if err != nil {
		return nil, err
	}
	req.URL = url
	req.RequestURI = ""

	return req, nil
}
