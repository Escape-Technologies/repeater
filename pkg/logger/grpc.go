package logger

import (
	"log"

	grpc "github.com/Escape-Technologies/repeater/pkg/grpc"
)


func ConnectLogs(url, repeaterId string) (hasConnected bool) {
	stream, closer, err := grpc.LogStream(url, repeaterId)
	defer closer()
	if err != nil {
		log.Printf("Error creating stream: %v \n", err)
		return false
	}
	log.Println("Connected to server...")

	for {
		msg := <-logSink
		err := stream.Send(msg)
		if err != nil {
			log.Printf("Error receiving: %v \n", err)
			return true
		}
	}
}
