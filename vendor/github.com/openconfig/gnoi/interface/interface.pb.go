// Code generated by protoc-gen-go. DO NOT EDIT.
// source: interface/interface.proto

package gnoi_interface

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	types "github.com/openconfig/gnoi/types"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// SetLoopbackModeRequest requests the provide interface be have its loopback mode
// set to mode.  Mode may be vendor specific.  For example, on a transport
// device, available modes are "none", "mac", "phy", "phy_remote",
// "framer_facility", and "framer_terminal".
type SetLoopbackModeRequest struct {
	Interface            *types.Path `protobuf:"bytes,1,opt,name=interface,proto3" json:"interface,omitempty"`
	Mode                 string      `protobuf:"bytes,2,opt,name=mode,proto3" json:"mode,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *SetLoopbackModeRequest) Reset()         { *m = SetLoopbackModeRequest{} }
func (m *SetLoopbackModeRequest) String() string { return proto.CompactTextString(m) }
func (*SetLoopbackModeRequest) ProtoMessage()    {}
func (*SetLoopbackModeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_460c38aebb3cb2d6, []int{0}
}

func (m *SetLoopbackModeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetLoopbackModeRequest.Unmarshal(m, b)
}
func (m *SetLoopbackModeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetLoopbackModeRequest.Marshal(b, m, deterministic)
}
func (m *SetLoopbackModeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetLoopbackModeRequest.Merge(m, src)
}
func (m *SetLoopbackModeRequest) XXX_Size() int {
	return xxx_messageInfo_SetLoopbackModeRequest.Size(m)
}
func (m *SetLoopbackModeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetLoopbackModeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetLoopbackModeRequest proto.InternalMessageInfo

func (m *SetLoopbackModeRequest) GetInterface() *types.Path {
	if m != nil {
		return m.Interface
	}
	return nil
}

func (m *SetLoopbackModeRequest) GetMode() string {
	if m != nil {
		return m.Mode
	}
	return ""
}

type SetLoopbackModeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetLoopbackModeResponse) Reset()         { *m = SetLoopbackModeResponse{} }
func (m *SetLoopbackModeResponse) String() string { return proto.CompactTextString(m) }
func (*SetLoopbackModeResponse) ProtoMessage()    {}
func (*SetLoopbackModeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_460c38aebb3cb2d6, []int{1}
}

func (m *SetLoopbackModeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetLoopbackModeResponse.Unmarshal(m, b)
}
func (m *SetLoopbackModeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetLoopbackModeResponse.Marshal(b, m, deterministic)
}
func (m *SetLoopbackModeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetLoopbackModeResponse.Merge(m, src)
}
func (m *SetLoopbackModeResponse) XXX_Size() int {
	return xxx_messageInfo_SetLoopbackModeResponse.Size(m)
}
func (m *SetLoopbackModeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetLoopbackModeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetLoopbackModeResponse proto.InternalMessageInfo

type GetLoopbackModeRequest struct {
	Interface            *types.Path `protobuf:"bytes,1,opt,name=interface,proto3" json:"interface,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetLoopbackModeRequest) Reset()         { *m = GetLoopbackModeRequest{} }
func (m *GetLoopbackModeRequest) String() string { return proto.CompactTextString(m) }
func (*GetLoopbackModeRequest) ProtoMessage()    {}
func (*GetLoopbackModeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_460c38aebb3cb2d6, []int{2}
}

func (m *GetLoopbackModeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLoopbackModeRequest.Unmarshal(m, b)
}
func (m *GetLoopbackModeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLoopbackModeRequest.Marshal(b, m, deterministic)
}
func (m *GetLoopbackModeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLoopbackModeRequest.Merge(m, src)
}
func (m *GetLoopbackModeRequest) XXX_Size() int {
	return xxx_messageInfo_GetLoopbackModeRequest.Size(m)
}
func (m *GetLoopbackModeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLoopbackModeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetLoopbackModeRequest proto.InternalMessageInfo

func (m *GetLoopbackModeRequest) GetInterface() *types.Path {
	if m != nil {
		return m.Interface
	}
	return nil
}

type GetLoopbackModeResponse struct {
	Mode                 string   `protobuf:"bytes,1,opt,name=mode,proto3" json:"mode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetLoopbackModeResponse) Reset()         { *m = GetLoopbackModeResponse{} }
func (m *GetLoopbackModeResponse) String() string { return proto.CompactTextString(m) }
func (*GetLoopbackModeResponse) ProtoMessage()    {}
func (*GetLoopbackModeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_460c38aebb3cb2d6, []int{3}
}

func (m *GetLoopbackModeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLoopbackModeResponse.Unmarshal(m, b)
}
func (m *GetLoopbackModeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLoopbackModeResponse.Marshal(b, m, deterministic)
}
func (m *GetLoopbackModeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLoopbackModeResponse.Merge(m, src)
}
func (m *GetLoopbackModeResponse) XXX_Size() int {
	return xxx_messageInfo_GetLoopbackModeResponse.Size(m)
}
func (m *GetLoopbackModeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLoopbackModeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetLoopbackModeResponse proto.InternalMessageInfo

func (m *GetLoopbackModeResponse) GetMode() string {
	if m != nil {
		return m.Mode
	}
	return ""
}

type ClearInterfaceCountersRequest struct {
	Interface            []*types.Path `protobuf:"bytes,1,rep,name=interface,proto3" json:"interface,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ClearInterfaceCountersRequest) Reset()         { *m = ClearInterfaceCountersRequest{} }
func (m *ClearInterfaceCountersRequest) String() string { return proto.CompactTextString(m) }
func (*ClearInterfaceCountersRequest) ProtoMessage()    {}
func (*ClearInterfaceCountersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_460c38aebb3cb2d6, []int{4}
}

func (m *ClearInterfaceCountersRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearInterfaceCountersRequest.Unmarshal(m, b)
}
func (m *ClearInterfaceCountersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearInterfaceCountersRequest.Marshal(b, m, deterministic)
}
func (m *ClearInterfaceCountersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearInterfaceCountersRequest.Merge(m, src)
}
func (m *ClearInterfaceCountersRequest) XXX_Size() int {
	return xxx_messageInfo_ClearInterfaceCountersRequest.Size(m)
}
func (m *ClearInterfaceCountersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearInterfaceCountersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClearInterfaceCountersRequest proto.InternalMessageInfo

func (m *ClearInterfaceCountersRequest) GetInterface() []*types.Path {
	if m != nil {
		return m.Interface
	}
	return nil
}

type ClearInterfaceCountersResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearInterfaceCountersResponse) Reset()         { *m = ClearInterfaceCountersResponse{} }
func (m *ClearInterfaceCountersResponse) String() string { return proto.CompactTextString(m) }
func (*ClearInterfaceCountersResponse) ProtoMessage()    {}
func (*ClearInterfaceCountersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_460c38aebb3cb2d6, []int{5}
}

func (m *ClearInterfaceCountersResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearInterfaceCountersResponse.Unmarshal(m, b)
}
func (m *ClearInterfaceCountersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearInterfaceCountersResponse.Marshal(b, m, deterministic)
}
func (m *ClearInterfaceCountersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearInterfaceCountersResponse.Merge(m, src)
}
func (m *ClearInterfaceCountersResponse) XXX_Size() int {
	return xxx_messageInfo_ClearInterfaceCountersResponse.Size(m)
}
func (m *ClearInterfaceCountersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearInterfaceCountersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ClearInterfaceCountersResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SetLoopbackModeRequest)(nil), "gnoi.interface.SetLoopbackModeRequest")
	proto.RegisterType((*SetLoopbackModeResponse)(nil), "gnoi.interface.SetLoopbackModeResponse")
	proto.RegisterType((*GetLoopbackModeRequest)(nil), "gnoi.interface.GetLoopbackModeRequest")
	proto.RegisterType((*GetLoopbackModeResponse)(nil), "gnoi.interface.GetLoopbackModeResponse")
	proto.RegisterType((*ClearInterfaceCountersRequest)(nil), "gnoi.interface.ClearInterfaceCountersRequest")
	proto.RegisterType((*ClearInterfaceCountersResponse)(nil), "gnoi.interface.ClearInterfaceCountersResponse")
}

func init() { proto.RegisterFile("interface/interface.proto", fileDescriptor_460c38aebb3cb2d6) }

var fileDescriptor_460c38aebb3cb2d6 = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcc, 0xcc, 0x2b, 0x49,
	0x2d, 0x4a, 0x4b, 0x4c, 0x4e, 0xd5, 0x87, 0xb3, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xf8,
	0xd2, 0xf3, 0xf2, 0x33, 0xf5, 0xe0, 0xa2, 0x52, 0x3a, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a,
	0xc9, 0xf9, 0xb9, 0xfa, 0xf9, 0x05, 0xa9, 0x79, 0xc9, 0xf9, 0x79, 0x69, 0x99, 0xe9, 0xfa, 0x20,
	0x55, 0xfa, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x10, 0x12, 0xa2, 0x5b, 0x29, 0x86, 0x4b, 0x2c, 0x38,
	0xb5, 0xc4, 0x27, 0x3f, 0xbf, 0x20, 0x29, 0x31, 0x39, 0xdb, 0x37, 0x3f, 0x25, 0x35, 0x28, 0xb5,
	0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x48, 0x8f, 0x8b, 0x13, 0x6e, 0xa8, 0x04, 0xa3, 0x02, 0xa3, 0x06,
	0xb7, 0x91, 0x80, 0x1e, 0xd8, 0x2e, 0x88, 0xfe, 0x80, 0xc4, 0x92, 0x8c, 0x20, 0x84, 0x12, 0x21,
	0x21, 0x2e, 0x96, 0xdc, 0xfc, 0x94, 0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x30, 0x5b,
	0x49, 0x92, 0x4b, 0x1c, 0xc3, 0xf4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x25, 0x0f, 0x2e, 0x31,
	0x77, 0xaa, 0x58, 0xac, 0xa4, 0xcb, 0x25, 0xee, 0x8e, 0xdd, 0x12, 0xb8, 0x9b, 0x18, 0x91, 0xdc,
	0xe4, 0xcf, 0x25, 0xeb, 0x9c, 0x93, 0x9a, 0x58, 0xe4, 0x09, 0x33, 0xc0, 0x39, 0xbf, 0x14, 0xc4,
	0x2c, 0xc6, 0x61, 0x3f, 0x33, 0x21, 0xfb, 0x15, 0xb8, 0xe4, 0x70, 0x19, 0x08, 0x71, 0x86, 0xd1,
	0x25, 0x26, 0x2e, 0x4e, 0xb8, 0xac, 0x50, 0x0a, 0x17, 0x3f, 0x5a, 0xa0, 0x08, 0xa9, 0xe9, 0xa1,
	0x46, 0xa2, 0x1e, 0xf6, 0x38, 0x91, 0x52, 0x27, 0xa8, 0x0e, 0x1a, 0xba, 0x0c, 0x20, 0x5b, 0xdc,
	0x09, 0xd9, 0xe2, 0x4e, 0xa4, 0x2d, 0xee, 0x38, 0x6d, 0xa9, 0xe4, 0x12, 0xc3, 0xee, 0x77, 0x21,
	0x5d, 0x74, 0x43, 0xf0, 0x06, 0xba, 0x94, 0x1e, 0xb1, 0xca, 0x61, 0x56, 0x3b, 0x71, 0x5c, 0xb2,
	0x63, 0x35, 0xd0, 0x33, 0xd4, 0x33, 0x48, 0x62, 0x03, 0x27, 0x65, 0x63, 0x40, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x00, 0x88, 0x3f, 0x97, 0x25, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// InterfaceClient is the client API for Interface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InterfaceClient interface {
	// SetLoopbackMode is used to set the mode of loopback on a interface.
	SetLoopbackMode(ctx context.Context, in *SetLoopbackModeRequest, opts ...grpc.CallOption) (*SetLoopbackModeResponse, error)
	// GetLoopbackMode is used to get the mode of loopback on a interface.
	GetLoopbackMode(ctx context.Context, in *GetLoopbackModeRequest, opts ...grpc.CallOption) (*GetLoopbackModeResponse, error)
	// ClearInterfaceCounters will reset the counters for the provided interface.
	ClearInterfaceCounters(ctx context.Context, in *ClearInterfaceCountersRequest, opts ...grpc.CallOption) (*ClearInterfaceCountersResponse, error)
}

type interfaceClient struct {
	cc *grpc.ClientConn
}

func NewInterfaceClient(cc *grpc.ClientConn) InterfaceClient {
	return &interfaceClient{cc}
}

func (c *interfaceClient) SetLoopbackMode(ctx context.Context, in *SetLoopbackModeRequest, opts ...grpc.CallOption) (*SetLoopbackModeResponse, error) {
	out := new(SetLoopbackModeResponse)
	err := c.cc.Invoke(ctx, "/gnoi.interface.Interface/SetLoopbackMode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interfaceClient) GetLoopbackMode(ctx context.Context, in *GetLoopbackModeRequest, opts ...grpc.CallOption) (*GetLoopbackModeResponse, error) {
	out := new(GetLoopbackModeResponse)
	err := c.cc.Invoke(ctx, "/gnoi.interface.Interface/GetLoopbackMode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interfaceClient) ClearInterfaceCounters(ctx context.Context, in *ClearInterfaceCountersRequest, opts ...grpc.CallOption) (*ClearInterfaceCountersResponse, error) {
	out := new(ClearInterfaceCountersResponse)
	err := c.cc.Invoke(ctx, "/gnoi.interface.Interface/ClearInterfaceCounters", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InterfaceServer is the server API for Interface service.
type InterfaceServer interface {
	// SetLoopbackMode is used to set the mode of loopback on a interface.
	SetLoopbackMode(context.Context, *SetLoopbackModeRequest) (*SetLoopbackModeResponse, error)
	// GetLoopbackMode is used to get the mode of loopback on a interface.
	GetLoopbackMode(context.Context, *GetLoopbackModeRequest) (*GetLoopbackModeResponse, error)
	// ClearInterfaceCounters will reset the counters for the provided interface.
	ClearInterfaceCounters(context.Context, *ClearInterfaceCountersRequest) (*ClearInterfaceCountersResponse, error)
}

func RegisterInterfaceServer(s *grpc.Server, srv InterfaceServer) {
	s.RegisterService(&_Interface_serviceDesc, srv)
}

func _Interface_SetLoopbackMode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetLoopbackModeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServer).SetLoopbackMode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.interface.Interface/SetLoopbackMode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServer).SetLoopbackMode(ctx, req.(*SetLoopbackModeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Interface_GetLoopbackMode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLoopbackModeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServer).GetLoopbackMode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.interface.Interface/GetLoopbackMode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServer).GetLoopbackMode(ctx, req.(*GetLoopbackModeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Interface_ClearInterfaceCounters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearInterfaceCountersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServer).ClearInterfaceCounters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.interface.Interface/ClearInterfaceCounters",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServer).ClearInterfaceCounters(ctx, req.(*ClearInterfaceCountersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Interface_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gnoi.interface.Interface",
	HandlerType: (*InterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetLoopbackMode",
			Handler:    _Interface_SetLoopbackMode_Handler,
		},
		{
			MethodName: "GetLoopbackMode",
			Handler:    _Interface_GetLoopbackMode_Handler,
		},
		{
			MethodName: "ClearInterfaceCounters",
			Handler:    _Interface_ClearInterfaceCounters_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interface/interface.proto",
}