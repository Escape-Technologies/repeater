package logger

import (
	"fmt"
	"log"
	"time"

	grpc "github.com/Escape-Technologies/repeater/pkg/grpc"
	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	fmt.Printf("%v %v %v\n", timestamppb.Now().AsTime(), proto.LogLevel_INFO, "Log sender connected to server...")

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
