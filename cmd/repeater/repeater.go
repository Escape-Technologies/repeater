package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Escape-Technologies/repeater/internal"
	"github.com/Escape-Technologies/repeater/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var UUID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var url = "0.0.0.0:8080"

func main() {
	// username := os.Getenv("ESCAPE_ORGANIZATION_ID")
	// if !UUID.MatchString(username) {
	// 	log.Printf("ESCAPE_ORGANIZATION_ID must be a UUID in lowercase")
	// 	log.Printf("To get your organization id, go to https://app.escape.tech/organization/settings/")
	// 	os.Exit(1)
	// }
	// password := os.Getenv("ESCAPE_REPEATER_ID")
	// if !UUID.MatchString(password) {
	// 	log.Printf("ESCAPE_REPEATER_ID must be a UUID in lowercase")
	// 	log.Printf("To get your API key, go to https://app.escape.tech/user/profile/")
	// 	os.Exit(1)
	// }

	start("username", "password")
}

func start(user string, pass string) {
	fmt.Println("Starting repeater client...")

	con, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}

	defer con.Close()
	client := proto.NewRepeaterClient(con)
	stream, err := client.Stream(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v \n", err)
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving: %v \n", err)
			continue
		}

		// Send request to server
		// Use a go func to avoid blocking the stream
		go func() {
			cReq, err := internal.TransportToRequest(req)
			if err != nil {
				return
			}

			// work
			cRes, err := http.DefaultClient.Do(cReq)
			if err != nil {
				return
			}

			tRes, err := internal.ResponseToTransport(cRes, req.Correlation)
			if err != nil {
				return
			}

			err = stream.Send(tRes)
			if err != nil {
				return
			}
		}()

	}
}
