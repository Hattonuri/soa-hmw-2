// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: internal/proto/base.proto

package base

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MafiaService_GetPlayersList_FullMethodName = "/base.MafiaService/GetPlayersList"
	MafiaService_Join_FullMethodName           = "/base.MafiaService/Join"
	MafiaService_ResponseEvent_FullMethodName  = "/base.MafiaService/ResponseEvent"
)

// MafiaServiceClient is the client API for MafiaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MafiaServiceClient interface {
	GetPlayersList(ctx context.Context, in *Room, opts ...grpc.CallOption) (*Players, error)
	Join(ctx context.Context, in *Player, opts ...grpc.CallOption) (MafiaService_JoinClient, error)
	ResponseEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*empty.Empty, error)
}

type mafiaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMafiaServiceClient(cc grpc.ClientConnInterface) MafiaServiceClient {
	return &mafiaServiceClient{cc}
}

func (c *mafiaServiceClient) GetPlayersList(ctx context.Context, in *Room, opts ...grpc.CallOption) (*Players, error) {
	out := new(Players)
	err := c.cc.Invoke(ctx, MafiaService_GetPlayersList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mafiaServiceClient) Join(ctx context.Context, in *Player, opts ...grpc.CallOption) (MafiaService_JoinClient, error) {
	stream, err := c.cc.NewStream(ctx, &MafiaService_ServiceDesc.Streams[0], MafiaService_Join_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &mafiaServiceJoinClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MafiaService_JoinClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type mafiaServiceJoinClient struct {
	grpc.ClientStream
}

func (x *mafiaServiceJoinClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mafiaServiceClient) ResponseEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, MafiaService_ResponseEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MafiaServiceServer is the server API for MafiaService service.
// All implementations must embed UnimplementedMafiaServiceServer
// for forward compatibility
type MafiaServiceServer interface {
	GetPlayersList(context.Context, *Room) (*Players, error)
	Join(*Player, MafiaService_JoinServer) error
	ResponseEvent(context.Context, *Event) (*empty.Empty, error)
	mustEmbedUnimplementedMafiaServiceServer()
}

// UnimplementedMafiaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMafiaServiceServer struct {
}

func (UnimplementedMafiaServiceServer) GetPlayersList(context.Context, *Room) (*Players, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlayersList not implemented")
}
func (UnimplementedMafiaServiceServer) Join(*Player, MafiaService_JoinServer) error {
	return status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedMafiaServiceServer) ResponseEvent(context.Context, *Event) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResponseEvent not implemented")
}
func (UnimplementedMafiaServiceServer) mustEmbedUnimplementedMafiaServiceServer() {}

// UnsafeMafiaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MafiaServiceServer will
// result in compilation errors.
type UnsafeMafiaServiceServer interface {
	mustEmbedUnimplementedMafiaServiceServer()
}

func RegisterMafiaServiceServer(s grpc.ServiceRegistrar, srv MafiaServiceServer) {
	s.RegisterService(&MafiaService_ServiceDesc, srv)
}

func _MafiaService_GetPlayersList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Room)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServiceServer).GetPlayersList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MafiaService_GetPlayersList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServiceServer).GetPlayersList(ctx, req.(*Room))
	}
	return interceptor(ctx, in, info, handler)
}

func _MafiaService_Join_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Player)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MafiaServiceServer).Join(m, &mafiaServiceJoinServer{stream})
}

type MafiaService_JoinServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type mafiaServiceJoinServer struct {
	grpc.ServerStream
}

func (x *mafiaServiceJoinServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func _MafiaService_ResponseEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MafiaServiceServer).ResponseEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MafiaService_ResponseEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MafiaServiceServer).ResponseEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

// MafiaService_ServiceDesc is the grpc.ServiceDesc for MafiaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MafiaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "base.MafiaService",
	HandlerType: (*MafiaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPlayersList",
			Handler:    _MafiaService_GetPlayersList_Handler,
		},
		{
			MethodName: "ResponseEvent",
			Handler:    _MafiaService_ResponseEvent_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Join",
			Handler:       _MafiaService_Join_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/proto/base.proto",
}
