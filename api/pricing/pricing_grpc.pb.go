// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: pricing.proto

package pricing

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

// PricingServiceClient is the client API for PricingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PricingServiceClient interface {
	GetPrice(ctx context.Context, in *GetPriceRequest, opts ...grpc.CallOption) (*GetPriceResponse, error)
}

type pricingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPricingServiceClient(cc grpc.ClientConnInterface) PricingServiceClient {
	return &pricingServiceClient{cc}
}

func (c *pricingServiceClient) GetPrice(ctx context.Context, in *GetPriceRequest, opts ...grpc.CallOption) (*GetPriceResponse, error) {
	out := new(GetPriceResponse)
	err := c.cc.Invoke(ctx, "/pricing.PricingService/GetPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PricingServiceServer is the server API for PricingService service.
// All implementations must embed UnimplementedPricingServiceServer
// for forward compatibility
type PricingServiceServer interface {
	GetPrice(context.Context, *GetPriceRequest) (*GetPriceResponse, error)
	mustEmbedUnimplementedPricingServiceServer()
}

// UnimplementedPricingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPricingServiceServer struct {
}

func (UnimplementedPricingServiceServer) GetPrice(context.Context, *GetPriceRequest) (*GetPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrice not implemented")
}
func (UnimplementedPricingServiceServer) mustEmbedUnimplementedPricingServiceServer() {}

// UnsafePricingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PricingServiceServer will
// result in compilation errors.
type UnsafePricingServiceServer interface {
	mustEmbedUnimplementedPricingServiceServer()
}

func RegisterPricingServiceServer(s grpc.ServiceRegistrar, srv PricingServiceServer) {
	s.RegisterService(&PricingService_ServiceDesc, srv)
}

func _PricingService_GetPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PricingServiceServer).GetPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pricing.PricingService/GetPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PricingServiceServer).GetPrice(ctx, req.(*GetPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PricingService_ServiceDesc is the grpc.ServiceDesc for PricingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PricingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pricing.PricingService",
	HandlerType: (*PricingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPrice",
			Handler:    _PricingService_GetPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pricing.proto",
}
