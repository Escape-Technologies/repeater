package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Escape-Technologies/repeater/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
)

// Injected by ldflags
var (
	version = "dev"
	commit  = "none"
)

var UUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func main() {
	log.Printf("Running Escape repeater version %s, commit %s\n", version, commit)

	repeaterId := os.Getenv("ESCAPE_REPEATER_ID")
	if !UUID.MatchString(repeaterId) {
		log.Println("ESCAPE_REPEATER_ID must be a UUID in lowercase")
		log.Println("To get your repeater id, go to https://app.escape.tech/organization/network/")
		log.Println("For more information, read the docs at https://docs.escape.tech/enterprise/repeater")
		os.Exit(1)
	}

	start(repeaterId)
}

func getCon() *grpc.ClientConn {
	url := os.Getenv("ESCAPE_REPEATER_URL")
	if url == "" {
		url = "repeater.escape.tech:443"
	} else {
		log.Printf("Using custom repeater url: %s\n", url)
	}

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
	con, err := grpc.Dial(url, creds)
	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}
	return con
}

func start(repeaterId string) {
	log.Println("Starting repeater client...")

	con := getCon()
	defer con.Close()

	client := proto.NewRepeaterClient(con)

	for {
		hasConnected := connectAndRun(client, repeaterId)
		log.Println("Disconnected...")
		if !hasConnected {
			log.Println("Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}

func connectAndRun(client proto.RepeaterClient, repeaterId string) (hasConnected bool) {
	hasConnected = false
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", repeaterId)
	stream, err := client.Stream(ctx)
	if err != nil {
		log.Printf("Error creating stream: %v \n", err)
		return hasConnected
	}
	hasConnected = true
	log.Println("Connected to server...")

	// Send healthcheck to the server
	go func() {
		retries := 0

		for {
			log.Println("Sending healthcheck...")
			err = stream.Send(&proto.Response{
				Data:        []byte(""),
				Correlation: 0,
			})
			if err != nil {
				log.Printf("Error sending healthcheck: %v\n", err)
				retries++
				if retries > 5 {
					log.Println("Too many retries, stopping healthchecks...")
					return
				}
			} else {
				retries = 0
			}
			log.Println("Healthcheck sent")
			time.Sleep(30 * time.Second)
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving: %v \n", err)
			return hasConnected
		}
		log.Printf("Received incoming stream (%d)\n", req.Correlation)

		// Send request to server
		// Use a go func to avoid blocking the stream
		go func() {
			startTime := time.Now()
			res := internal.HandleRequest(req)
			log.Printf("Processed stream in %v (%d)\n", time.Since(startTime), req.Correlation)

			err = stream.Send(res)
			if err != nil {
				log.Printf("Error processing stream (%d): %v \n", req.Correlation, err)
			}
		}()
	}
}
