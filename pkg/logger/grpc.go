package logger

import (
	"log"
	"time"

	grpc "github.com/Escape-Technologies/repeater/pkg/grpc"
)

func AlwaysConnect(url, repeaterId string) {
	for {
		ConnectLogs(url, repeaterId)
		time.Sleep(time.Second)
	}
}

func ConnectLogs(url, repeaterId string) (hasConnected bool) {
	stream, closer, err := grpc.LogStream(url, repeaterId)
	defer closer()
	if err != nil {
		log.Printf("Error creating stream: %v \n", err)
		return false
	}
	log.Println("Connected to server...")

	for {
		msg := queue.Next()
		if msg == nil {
			time.Sleep(time.Second)
			continue
		}
		err := stream.Send(msg)
		if err != nil {
			log.Printf("Error receiving: %v \n", err)
			return true
		}
	}
}
