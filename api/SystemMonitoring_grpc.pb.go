// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// MonitoringClient is the client API for Monitoring service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MonitoringClient interface {
	GetInfo(ctx context.Context, in *Request, opts ...grpc.CallOption) (Monitoring_GetInfoClient, error)
}

type monitoringClient struct {
	cc grpc.ClientConnInterface
}

func NewMonitoringClient(cc grpc.ClientConnInterface) MonitoringClient {
	return &monitoringClient{cc}
}

func (c *monitoringClient) GetInfo(ctx context.Context, in *Request, opts ...grpc.CallOption) (Monitoring_GetInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &Monitoring_ServiceDesc.Streams[0], "/api.Monitoring/GetInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &monitoringGetInfoClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Monitoring_GetInfoClient interface {
	Recv() (*Result, error)
	grpc.ClientStream
}

type monitoringGetInfoClient struct {
	grpc.ClientStream
}

func (x *monitoringGetInfoClient) Recv() (*Result, error) {
	m := new(Result)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MonitoringServer is the server API for Monitoring service.
// All implementations must embed UnimplementedMonitoringServer
// for forward compatibility
type MonitoringServer interface {
	GetInfo(*Request, Monitoring_GetInfoServer) error
	mustEmbedUnimplementedMonitoringServer()
}

// UnimplementedMonitoringServer must be embedded to have forward compatible implementations.
type UnimplementedMonitoringServer struct {
}

func (UnimplementedMonitoringServer) GetInfo(*Request, Monitoring_GetInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedMonitoringServer) mustEmbedUnimplementedMonitoringServer() {}

// UnsafeMonitoringServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MonitoringServer will
// result in compilation errors.
type UnsafeMonitoringServer interface {
	mustEmbedUnimplementedMonitoringServer()
}

func RegisterMonitoringServer(s grpc.ServiceRegistrar, srv MonitoringServer) {
	s.RegisterService(&Monitoring_ServiceDesc, srv)
}

func _Monitoring_GetInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MonitoringServer).GetInfo(m, &monitoringGetInfoServer{stream})
}

type Monitoring_GetInfoServer interface {
	Send(*Result) error
	grpc.ServerStream
}

type monitoringGetInfoServer struct {
	grpc.ServerStream
}

func (x *monitoringGetInfoServer) Send(m *Result) error {
	return x.ServerStream.SendMsg(m)
}

// Monitoring_ServiceDesc is the grpc.ServiceDesc for Monitoring service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Monitoring_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Monitoring",
	HandlerType: (*MonitoringServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetInfo",
			Handler:       _Monitoring_GetInfo_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "SystemMonitoring.proto",
}
