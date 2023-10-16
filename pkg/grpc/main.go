package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
)


func LogStream(url, repeaterId string) (stream proto.Repeater_LogStreamClient, closer func(), err error) {
	con := GetCon(url)

	client := proto.NewRepeaterClient(con)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", repeaterId)
	stream, err = client.LogStream(ctx)

	return stream, func() { con.Close() }, err
}

func Stream(url, repeaterId string) (stream proto.Repeater_StreamClient, closer func(), err error) {
	con := GetCon(url)

	client := proto.NewRepeaterClient(con)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", repeaterId)
	stream, err = client.Stream(ctx)

	return stream, func() { con.Close() }, err
}
