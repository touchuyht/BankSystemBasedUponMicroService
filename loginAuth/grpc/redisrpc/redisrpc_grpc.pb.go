// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package redisrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// IRedisRPCClient is the client API for IRedisRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IRedisRPCClient interface {
	UpdateToken(ctx context.Context, in *User, opts ...grpc.CallOption) (*Result, error)
	InsertUsers(ctx context.Context, opts ...grpc.CallOption) (IRedisRPC_InsertUsersClient, error)
}

type iRedisRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewIRedisRPCClient(cc grpc.ClientConnInterface) IRedisRPCClient {
	return &iRedisRPCClient{cc}
}

func (c *iRedisRPCClient) UpdateToken(ctx context.Context, in *User, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/redisrpc.IRedisRPC/UpdateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iRedisRPCClient) InsertUsers(ctx context.Context, opts ...grpc.CallOption) (IRedisRPC_InsertUsersClient, error) {
	stream, err := c.cc.NewStream(ctx, &_IRedisRPC_serviceDesc.Streams[0], "/redisrpc.IRedisRPC/InsertUsers", opts...)
	if err != nil {
		return nil, err
	}
	x := &iRedisRPCInsertUsersClient{stream}
	return x, nil
}

type IRedisRPC_InsertUsersClient interface {
	Send(*User) error
	CloseAndRecv() (*Result, error)
	grpc.ClientStream
}

type iRedisRPCInsertUsersClient struct {
	grpc.ClientStream
}

func (x *iRedisRPCInsertUsersClient) Send(m *User) error {
	return x.ClientStream.SendMsg(m)
}

func (x *iRedisRPCInsertUsersClient) CloseAndRecv() (*Result, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Result)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// IRedisRPCServer is the server API for IRedisRPC service.
// All implementations must embed UnimplementedIRedisRPCServer
// for forward compatibility
type IRedisRPCServer interface {
	UpdateToken(context.Context, *User) (*Result, error)
	InsertUsers(IRedisRPC_InsertUsersServer) error
	mustEmbedUnimplementedIRedisRPCServer()
}

// UnimplementedIRedisRPCServer must be embedded to have forward compatible implementations.
type UnimplementedIRedisRPCServer struct {
}

func (UnimplementedIRedisRPCServer) UpdateToken(context.Context, *User) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateToken not implemented")
}
func (UnimplementedIRedisRPCServer) InsertUsers(IRedisRPC_InsertUsersServer) error {
	return status.Errorf(codes.Unimplemented, "method InsertUsers not implemented")
}
func (UnimplementedIRedisRPCServer) mustEmbedUnimplementedIRedisRPCServer() {}

// UnsafeIRedisRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IRedisRPCServer will
// result in compilation errors.
type UnsafeIRedisRPCServer interface {
	mustEmbedUnimplementedIRedisRPCServer()
}

func RegisterIRedisRPCServer(s grpc.ServiceRegistrar, srv IRedisRPCServer) {
	s.RegisterService(&_IRedisRPC_serviceDesc, srv)
}

func _IRedisRPC_UpdateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IRedisRPCServer).UpdateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redisrpc.IRedisRPC/UpdateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IRedisRPCServer).UpdateToken(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _IRedisRPC_InsertUsers_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IRedisRPCServer).InsertUsers(&iRedisRPCInsertUsersServer{stream})
}

type IRedisRPC_InsertUsersServer interface {
	SendAndClose(*Result) error
	Recv() (*User, error)
	grpc.ServerStream
}

type iRedisRPCInsertUsersServer struct {
	grpc.ServerStream
}

func (x *iRedisRPCInsertUsersServer) SendAndClose(m *Result) error {
	return x.ServerStream.SendMsg(m)
}

func (x *iRedisRPCInsertUsersServer) Recv() (*User, error) {
	m := new(User)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _IRedisRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "redisrpc.IRedisRPC",
	HandlerType: (*IRedisRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateToken",
			Handler:    _IRedisRPC_UpdateToken_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "InsertUsers",
			Handler:       _IRedisRPC_InsertUsers_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "grpc/redisrpc/redisrpc.proto",
}
