// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: repeater/v1/main.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Repeater_Stream_FullMethodName = "/repeater.v1.Repeater/Stream"
)

// RepeaterClient is the client API for Repeater service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepeaterClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (Repeater_StreamClient, error)
}

type repeaterClient struct {
	cc grpc.ClientConnInterface
}

func NewRepeaterClient(cc grpc.ClientConnInterface) RepeaterClient {
	return &repeaterClient{cc}
}

func (c *repeaterClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Repeater_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Repeater_ServiceDesc.Streams[0], Repeater_Stream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &repeaterStreamClient{stream}
	return x, nil
}

type Repeater_StreamClient interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ClientStream
}

type repeaterStreamClient struct {
	grpc.ClientStream
}

func (x *repeaterStreamClient) Send(m *Response) error {
	return x.ClientStream.SendMsg(m)
}

func (x *repeaterStreamClient) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RepeaterServer is the server API for Repeater service.
// All implementations must embed UnimplementedRepeaterServer
// for forward compatibility
type RepeaterServer interface {
	Stream(Repeater_StreamServer) error
	mustEmbedUnimplementedRepeaterServer()
}

// UnimplementedRepeaterServer must be embedded to have forward compatible implementations.
type UnimplementedRepeaterServer struct {
}

func (UnimplementedRepeaterServer) Stream(Repeater_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedRepeaterServer) mustEmbedUnimplementedRepeaterServer() {}

// UnsafeRepeaterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepeaterServer will
// result in compilation errors.
type UnsafeRepeaterServer interface {
	mustEmbedUnimplementedRepeaterServer()
}

func RegisterRepeaterServer(s grpc.ServiceRegistrar, srv RepeaterServer) {
	s.RegisterService(&Repeater_ServiceDesc, srv)
}

func _Repeater_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RepeaterServer).Stream(&repeaterStreamServer{stream})
}

type Repeater_StreamServer interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ServerStream
}

type repeaterStreamServer struct {
	grpc.ServerStream
}

func (x *repeaterStreamServer) Send(m *Request) error {
	return x.ServerStream.SendMsg(m)
}

func (x *repeaterStreamServer) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Repeater_ServiceDesc is the grpc.ServiceDesc for Repeater service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Repeater_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "repeater.v1.Repeater",
	HandlerType: (*RepeaterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Repeater_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "repeater/v1/main.proto",
}