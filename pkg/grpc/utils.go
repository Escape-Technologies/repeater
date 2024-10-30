package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var maxMsgSize = 1024 * 1024 * 1024

// Custom dialler for gRPC connection via proxy
func customDialer(proxyURL string) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, target string) (net.Conn, error) {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		// Establish HTTP2-compatible connection via proxy
		dialer := &net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}
		return dialer.DialContext(ctx, "tcp", proxy.Host)
	}
}

func GetCon(url, proxyURL string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	}

	if strings.HasPrefix(url, "localhost") {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("Error connecting: %v \n", err)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	if proxyURL != "" {
		opts = append(opts, grpc.WithContextDialer(customDialer(proxyURL)))
	}

	con, err := grpc.NewClient(url, opts...)
	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}
	return con
}
