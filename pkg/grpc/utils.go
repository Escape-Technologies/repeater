package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var maxMsgSize = 1024 * 1024 * 1024

func GetCon(url, proxyURL string) *grpc.ClientConn {
	var creds grpc.DialOption
	if strings.Split(url, ":")[0] == "localhost" {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("Error connecting: %v \n", err)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		creds = grpc.WithTransportCredentials(cred)
	}
	con, err := grpc.NewClient(
		url,
		creds,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	)
	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}
	return con
}
