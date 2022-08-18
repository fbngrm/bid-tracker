// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: auction/v1/auction.proto

package auctionv1

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	// Creates a new bid for an item.
	CreateBid(ctx context.Context, in *CreateBidRequest, opts ...grpc.CallOption) (*Bid, error)
	// Get the highest bid for an item.
	GetHighestBid(ctx context.Context, in *GetHighestBidRequest, opts ...grpc.CallOption) (*Bid, error)
	// Get all bids for an item.
	GetBids(ctx context.Context, in *GetBidsRequest, opts ...grpc.CallOption) (*Bids, error)
	// Get all items a user holds bids for.
	GetItemsForUserBids(ctx context.Context, in *GetItemsForUserBidsRequest, opts ...grpc.CallOption) (*Items, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) CreateBid(ctx context.Context, in *CreateBidRequest, opts ...grpc.CallOption) (*Bid, error) {
	out := new(Bid)
	err := c.cc.Invoke(ctx, "/auction.v1.Service/CreateBid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetHighestBid(ctx context.Context, in *GetHighestBidRequest, opts ...grpc.CallOption) (*Bid, error) {
	out := new(Bid)
	err := c.cc.Invoke(ctx, "/auction.v1.Service/GetHighestBid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetBids(ctx context.Context, in *GetBidsRequest, opts ...grpc.CallOption) (*Bids, error) {
	out := new(Bids)
	err := c.cc.Invoke(ctx, "/auction.v1.Service/GetBids", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetItemsForUserBids(ctx context.Context, in *GetItemsForUserBidsRequest, opts ...grpc.CallOption) (*Items, error) {
	out := new(Items)
	err := c.cc.Invoke(ctx, "/auction.v1.Service/GetItemsForUserBids", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations should embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	// Creates a new bid for an item.
	CreateBid(context.Context, *CreateBidRequest) (*Bid, error)
	// Get the highest bid for an item.
	GetHighestBid(context.Context, *GetHighestBidRequest) (*Bid, error)
	// Get all bids for an item.
	GetBids(context.Context, *GetBidsRequest) (*Bids, error)
	// Get all items a user holds bids for.
	GetItemsForUserBids(context.Context, *GetItemsForUserBidsRequest) (*Items, error)
}

// UnimplementedServiceServer should be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) CreateBid(context.Context, *CreateBidRequest) (*Bid, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBid not implemented")
}
func (UnimplementedServiceServer) GetHighestBid(context.Context, *GetHighestBidRequest) (*Bid, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHighestBid not implemented")
}
func (UnimplementedServiceServer) GetBids(context.Context, *GetBidsRequest) (*Bids, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBids not implemented")
}
func (UnimplementedServiceServer) GetItemsForUserBids(context.Context, *GetItemsForUserBidsRequest) (*Items, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemsForUserBids not implemented")
}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_CreateBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.v1.Service/CreateBid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateBid(ctx, req.(*CreateBidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetHighestBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHighestBidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetHighestBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.v1.Service/GetHighestBid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetHighestBid(ctx, req.(*GetHighestBidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetBids_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBidsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetBids(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.v1.Service/GetBids",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetBids(ctx, req.(*GetBidsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetItemsForUserBids_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemsForUserBidsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetItemsForUserBids(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.v1.Service/GetItemsForUserBids",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetItemsForUserBids(ctx, req.(*GetItemsForUserBidsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auction.v1.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBid",
			Handler:    _Service_CreateBid_Handler,
		},
		{
			MethodName: "GetHighestBid",
			Handler:    _Service_GetHighestBid_Handler,
		},
		{
			MethodName: "GetBids",
			Handler:    _Service_GetBids_Handler,
		},
		{
			MethodName: "GetItemsForUserBids",
			Handler:    _Service_GetItemsForUserBids_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auction/v1/auction.proto",
}
