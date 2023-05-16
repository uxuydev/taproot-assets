// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mintrpc

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

// MintClient is the client API for Mint service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MintClient interface {
	// tapcli: `assets mint`
	// MintAsset will attempt to mint the set of assets (async by default to
	// ensure proper batching) specified in the request.
	MintAsset(ctx context.Context, in *MintAssetRequest, opts ...grpc.CallOption) (*MintAssetResponse, error)
	// tapcli: `assets mint finalize`
	// FinalizeBatch will attempt to finalize the current pending batch.
	FinalizeBatch(ctx context.Context, in *FinalizeBatchRequest, opts ...grpc.CallOption) (*FinalizeBatchResponse, error)
	// tapcli: `assets mint cancel`
	// CancelBatch will attempt to cancel the current pending batch.
	CancelBatch(ctx context.Context, in *CancelBatchRequest, opts ...grpc.CallOption) (*CancelBatchResponse, error)
	// tapcli: `assets mint batches`
	// ListBatches lists the set of batches submitted to the daemon, including
	// pending and cancelled batches.
	ListBatches(ctx context.Context, in *ListBatchRequest, opts ...grpc.CallOption) (*ListBatchResponse, error)
}

type mintClient struct {
	cc grpc.ClientConnInterface
}

func NewMintClient(cc grpc.ClientConnInterface) MintClient {
	return &mintClient{cc}
}

func (c *mintClient) MintAsset(ctx context.Context, in *MintAssetRequest, opts ...grpc.CallOption) (*MintAssetResponse, error) {
	out := new(MintAssetResponse)
	err := c.cc.Invoke(ctx, "/mintrpc.Mint/MintAsset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mintClient) FinalizeBatch(ctx context.Context, in *FinalizeBatchRequest, opts ...grpc.CallOption) (*FinalizeBatchResponse, error) {
	out := new(FinalizeBatchResponse)
	err := c.cc.Invoke(ctx, "/mintrpc.Mint/FinalizeBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mintClient) CancelBatch(ctx context.Context, in *CancelBatchRequest, opts ...grpc.CallOption) (*CancelBatchResponse, error) {
	out := new(CancelBatchResponse)
	err := c.cc.Invoke(ctx, "/mintrpc.Mint/CancelBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mintClient) ListBatches(ctx context.Context, in *ListBatchRequest, opts ...grpc.CallOption) (*ListBatchResponse, error) {
	out := new(ListBatchResponse)
	err := c.cc.Invoke(ctx, "/mintrpc.Mint/ListBatches", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MintServer is the server API for Mint service.
// All implementations must embed UnimplementedMintServer
// for forward compatibility
type MintServer interface {
	// tapcli: `assets mint`
	// MintAsset will attempt to mint the set of assets (async by default to
	// ensure proper batching) specified in the request.
	MintAsset(context.Context, *MintAssetRequest) (*MintAssetResponse, error)
	// tapcli: `assets mint finalize`
	// FinalizeBatch will attempt to finalize the current pending batch.
	FinalizeBatch(context.Context, *FinalizeBatchRequest) (*FinalizeBatchResponse, error)
	// tapcli: `assets mint cancel`
	// CancelBatch will attempt to cancel the current pending batch.
	CancelBatch(context.Context, *CancelBatchRequest) (*CancelBatchResponse, error)
	// tapcli: `assets mint batches`
	// ListBatches lists the set of batches submitted to the daemon, including
	// pending and cancelled batches.
	ListBatches(context.Context, *ListBatchRequest) (*ListBatchResponse, error)
	mustEmbedUnimplementedMintServer()
}

// UnimplementedMintServer must be embedded to have forward compatible implementations.
type UnimplementedMintServer struct {
}

func (UnimplementedMintServer) MintAsset(context.Context, *MintAssetRequest) (*MintAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MintAsset not implemented")
}
func (UnimplementedMintServer) FinalizeBatch(context.Context, *FinalizeBatchRequest) (*FinalizeBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinalizeBatch not implemented")
}
func (UnimplementedMintServer) CancelBatch(context.Context, *CancelBatchRequest) (*CancelBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelBatch not implemented")
}
func (UnimplementedMintServer) ListBatches(context.Context, *ListBatchRequest) (*ListBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBatches not implemented")
}
func (UnimplementedMintServer) mustEmbedUnimplementedMintServer() {}

// UnsafeMintServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MintServer will
// result in compilation errors.
type UnsafeMintServer interface {
	mustEmbedUnimplementedMintServer()
}

func RegisterMintServer(s grpc.ServiceRegistrar, srv MintServer) {
	s.RegisterService(&Mint_ServiceDesc, srv)
}

func _Mint_MintAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MintAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MintServer).MintAsset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mintrpc.Mint/MintAsset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MintServer).MintAsset(ctx, req.(*MintAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mint_FinalizeBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinalizeBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MintServer).FinalizeBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mintrpc.Mint/FinalizeBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MintServer).FinalizeBatch(ctx, req.(*FinalizeBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mint_CancelBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MintServer).CancelBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mintrpc.Mint/CancelBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MintServer).CancelBatch(ctx, req.(*CancelBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mint_ListBatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MintServer).ListBatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mintrpc.Mint/ListBatches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MintServer).ListBatches(ctx, req.(*ListBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Mint_ServiceDesc is the grpc.ServiceDesc for Mint service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mint_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mintrpc.Mint",
	HandlerType: (*MintServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MintAsset",
			Handler:    _Mint_MintAsset_Handler,
		},
		{
			MethodName: "FinalizeBatch",
			Handler:    _Mint_FinalizeBatch_Handler,
		},
		{
			MethodName: "CancelBatch",
			Handler:    _Mint_CancelBatch_Handler,
		},
		{
			MethodName: "ListBatches",
			Handler:    _Mint_ListBatches_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mintrpc/mint.proto",
}
