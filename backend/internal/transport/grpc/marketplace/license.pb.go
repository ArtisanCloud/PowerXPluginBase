package marketplace

import (
	"context"
	"reflect"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type IssueLicenseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId        string `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	ListingId       string `protobuf:"bytes,2,opt,name=listing_id,json=listingId,proto3" json:"listing_id,omitempty"`
	PlanId          string `protobuf:"bytes,3,opt,name=plan_id,json=planId,proto3" json:"plan_id,omitempty"`
	IssuedBy        string `protobuf:"bytes,4,opt,name=issued_by,json=issuedBy,proto3" json:"issued_by,omitempty"`
	Trial           bool   `protobuf:"varint,5,opt,name=trial,proto3" json:"trial,omitempty"`
	ExpiresAtUnix   int64  `protobuf:"varint,6,opt,name=expires_at_unix,json=expiresAtUnix,proto3" json:"expires_at_unix,omitempty"`
	PaymentIntentId string `protobuf:"bytes,7,opt,name=payment_intent_id,json=paymentIntentId,proto3" json:"payment_intent_id,omitempty"`
}

type IssueLicenseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	License *License `protobuf:"bytes,1,opt,name=license,proto3" json:"license,omitempty"`
}

type RenewLicenseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId      string `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	LicenseId     string `protobuf:"bytes,2,opt,name=license_id,json=licenseId,proto3" json:"license_id,omitempty"`
	PlanId        string `protobuf:"bytes,3,opt,name=plan_id,json=planId,proto3" json:"plan_id,omitempty"`
	RenewalToken  string `protobuf:"bytes,4,opt,name=renewal_token,json=renewalToken,proto3" json:"renewal_token,omitempty"`
	ExpiresAtUnix int64  `protobuf:"varint,5,opt,name=expires_at_unix,json=expiresAtUnix,proto3" json:"expires_at_unix,omitempty"`
}

type RenewLicenseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	License *License `protobuf:"bytes,1,opt,name=license,proto3" json:"license,omitempty"`
}

type VerifyLicenseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId  string `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	ListingId string `protobuf:"bytes,2,opt,name=listing_id,json=listingId,proto3" json:"listing_id,omitempty"`
}

type VerifyLicenseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Valid   bool     `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	Reason  string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	License *License `protobuf:"bytes,3,opt,name=license,proto3" json:"license,omitempty"`
}

type License struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                 string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TenantId           string  `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	ListingId          string  `protobuf:"bytes,3,opt,name=listing_id,json=listingId,proto3" json:"listing_id,omitempty"`
	PlanId             string  `protobuf:"bytes,4,opt,name=plan_id,json=planId,proto3" json:"plan_id,omitempty"`
	Status             string  `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	Token              string  `protobuf:"bytes,6,opt,name=token,proto3" json:"token,omitempty"`
	ExpiresAtUnix      int64   `protobuf:"varint,7,opt,name=expires_at_unix,json=expiresAtUnix,proto3" json:"expires_at_unix,omitempty"`
	OfflineUntilUnix   int64   `protobuf:"varint,8,opt,name=offline_until_unix,json=offlineUntilUnix,proto3" json:"offline_until_unix,omitempty"`
	RenewalToken       string  `protobuf:"bytes,9,opt,name=renewal_token,json=renewalToken,proto3" json:"renewal_token,omitempty"`
	SettlementCurrency string  `protobuf:"bytes,10,opt,name=settlement_currency,json=settlementCurrency,proto3" json:"settlement_currency,omitempty"`
	ExchangeRate       float64 `protobuf:"fixed64,11,opt,name=exchange_rate,json=exchangeRate,proto3" json:"exchange_rate,omitempty"`
}

func (x *IssueLicenseRequest) Reset()         { *x = IssueLicenseRequest{} }
func (x *IssueLicenseRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*IssueLicenseRequest) ProtoMessage()    {}
func (x *IssueLicenseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *IssueLicenseResponse) Reset()         { *x = IssueLicenseResponse{} }
func (x *IssueLicenseResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*IssueLicenseResponse) ProtoMessage()    {}
func (x *IssueLicenseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *RenewLicenseRequest) Reset()         { *x = RenewLicenseRequest{} }
func (x *RenewLicenseRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*RenewLicenseRequest) ProtoMessage()    {}
func (x *RenewLicenseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *RenewLicenseResponse) Reset()         { *x = RenewLicenseResponse{} }
func (x *RenewLicenseResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*RenewLicenseResponse) ProtoMessage()    {}
func (x *RenewLicenseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *VerifyLicenseRequest) Reset()         { *x = VerifyLicenseRequest{} }
func (x *VerifyLicenseRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*VerifyLicenseRequest) ProtoMessage()    {}
func (x *VerifyLicenseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *VerifyLicenseResponse) Reset()         { *x = VerifyLicenseResponse{} }
func (x *VerifyLicenseResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*VerifyLicenseResponse) ProtoMessage()    {}
func (x *VerifyLicenseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *License) Reset()         { *x = License{} }
func (x *License) String() string { return protoimpl.X.MessageStringOf(x) }
func (*License) ProtoMessage()    {}
func (x *License) ProtoReflect() protoreflect.Message {
	mi := &file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_backend_internal_transport_grpc_marketplace_license_proto_goTypes = []interface{}{
	(*IssueLicenseRequest)(nil),
	(*IssueLicenseResponse)(nil),
	(*RenewLicenseRequest)(nil),
	(*RenewLicenseResponse)(nil),
	(*VerifyLicenseRequest)(nil),
	(*VerifyLicenseResponse)(nil),
	(*License)(nil),
}
var file_backend_internal_transport_grpc_marketplace_license_proto_depIdxs = []int32{
	6, // IssueLicenseResponse.license
	6, // RenewLicenseResponse.license
	6, // VerifyLicenseResponse.license
}

var File_backend_internal_transport_grpc_marketplace_license_proto protoreflect.FileDescriptor
var file_backend_internal_transport_grpc_marketplace_license_proto_once sync.Once

func file_backend_internal_transport_grpc_marketplace_license_proto_init() {
	file_backend_internal_transport_grpc_marketplace_license_proto_once.Do(func() {
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IssueLicenseRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IssueLicenseResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenewLicenseRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenewLicenseResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyLicenseRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyLicenseResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*License); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}

		type x struct{}
		out := protoimpl.TypeBuilder{
			File: protoimpl.DescBuilder{
				GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
				RawDescriptor: nil,
				NumEnums:      0,
				NumMessages:   7,
				NumExtensions: 0,
				NumServices:   1,
			},
			GoTypes:           file_backend_internal_transport_grpc_marketplace_license_proto_goTypes,
			DependencyIndexes: file_backend_internal_transport_grpc_marketplace_license_proto_depIdxs,
			MessageInfos:      file_backend_internal_transport_grpc_marketplace_license_proto_msgTypes,
		}.Build()

		File_backend_internal_transport_grpc_marketplace_license_proto = out.File
	})
}

func init() {
	file_backend_internal_transport_grpc_marketplace_license_proto_init()
}

// Reference imports to suppress errors if they are not otherwise used.
type _ = protoreflect.List

// LicenseServiceServer is the server API for LicenseService service.
type LicenseServiceServer interface {
	Issue(context.Context, *IssueLicenseRequest) (*IssueLicenseResponse, error)
	Renew(context.Context, *RenewLicenseRequest) (*RenewLicenseResponse, error)
	Verify(context.Context, *VerifyLicenseRequest) (*VerifyLicenseResponse, error)
	mustEmbedUnimplementedLicenseServiceServer()
}

type UnimplementedLicenseServiceServer struct{}

func (UnimplementedLicenseServiceServer) Issue(context.Context, *IssueLicenseRequest) (*IssueLicenseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Issue not implemented")
}
func (UnimplementedLicenseServiceServer) Renew(context.Context, *RenewLicenseRequest) (*RenewLicenseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Renew not implemented")
}
func (UnimplementedLicenseServiceServer) Verify(context.Context, *VerifyLicenseRequest) (*VerifyLicenseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}
func (UnimplementedLicenseServiceServer) mustEmbedUnimplementedLicenseServiceServer() {}

type UnsafeLicenseServiceServer interface {
	mustEmbedUnimplementedLicenseServiceServer()
}

func RegisterLicenseServiceServer(s grpc.ServiceRegistrar, srv LicenseServiceServer) {
	if srv != nil {
		s.RegisterService(&LicenseService_ServiceDesc, srv)
	}
}

func _LicenseService_Issue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueLicenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LicenseServiceServer).Issue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/powerx.marketplace.v1.LicenseService/Issue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LicenseServiceServer).Issue(ctx, req.(*IssueLicenseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LicenseService_Renew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenewLicenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LicenseServiceServer).Renew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/powerx.marketplace.v1.LicenseService/Renew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LicenseServiceServer).Renew(ctx, req.(*RenewLicenseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LicenseService_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyLicenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LicenseServiceServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/powerx.marketplace.v1.LicenseService/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LicenseServiceServer).Verify(ctx, req.(*VerifyLicenseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var LicenseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "powerx.marketplace.v1.LicenseService",
	HandlerType: (*LicenseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Issue",
			Handler:    _LicenseService_Issue_Handler,
		},
		{
			MethodName: "Renew",
			Handler:    _LicenseService_Renew_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _LicenseService_Verify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "marketplace/license.proto",
}
