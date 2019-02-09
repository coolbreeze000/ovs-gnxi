// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bgp/bgp.proto

package gnoi_bgp

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/openconfig/gnoi/types"
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

type ClearBGPNeighborRequest_Mode int32

const (
	ClearBGPNeighborRequest_SOFT   ClearBGPNeighborRequest_Mode = 0
	ClearBGPNeighborRequest_SOFTIN ClearBGPNeighborRequest_Mode = 1
	ClearBGPNeighborRequest_HARD   ClearBGPNeighborRequest_Mode = 2
)

var ClearBGPNeighborRequest_Mode_name = map[int32]string{
	0: "SOFT",
	1: "SOFTIN",
	2: "HARD",
}

var ClearBGPNeighborRequest_Mode_value = map[string]int32{
	"SOFT":   0,
	"SOFTIN": 1,
	"HARD":   2,
}

func (x ClearBGPNeighborRequest_Mode) String() string {
	return proto.EnumName(ClearBGPNeighborRequest_Mode_name, int32(x))
}

func (ClearBGPNeighborRequest_Mode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dd5905d96efb1c39, []int{0, 0}
}

type ClearBGPNeighborRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Routing instance containing the neighbor. Defaults to the global routing
	// table.
	RoutingInstance      string                       `protobuf:"bytes,2,opt,name=routing_instance,json=routingInstance,proto3" json:"routing_instance,omitempty"`
	Mode                 ClearBGPNeighborRequest_Mode `protobuf:"varint,3,opt,name=mode,proto3,enum=gnoi.bgp.ClearBGPNeighborRequest_Mode" json:"mode,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *ClearBGPNeighborRequest) Reset()         { *m = ClearBGPNeighborRequest{} }
func (m *ClearBGPNeighborRequest) String() string { return proto.CompactTextString(m) }
func (*ClearBGPNeighborRequest) ProtoMessage()    {}
func (*ClearBGPNeighborRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd5905d96efb1c39, []int{0}
}

func (m *ClearBGPNeighborRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearBGPNeighborRequest.Unmarshal(m, b)
}
func (m *ClearBGPNeighborRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearBGPNeighborRequest.Marshal(b, m, deterministic)
}
func (m *ClearBGPNeighborRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearBGPNeighborRequest.Merge(m, src)
}
func (m *ClearBGPNeighborRequest) XXX_Size() int {
	return xxx_messageInfo_ClearBGPNeighborRequest.Size(m)
}
func (m *ClearBGPNeighborRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearBGPNeighborRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClearBGPNeighborRequest proto.InternalMessageInfo

func (m *ClearBGPNeighborRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *ClearBGPNeighborRequest) GetRoutingInstance() string {
	if m != nil {
		return m.RoutingInstance
	}
	return ""
}

func (m *ClearBGPNeighborRequest) GetMode() ClearBGPNeighborRequest_Mode {
	if m != nil {
		return m.Mode
	}
	return ClearBGPNeighborRequest_SOFT
}

type ClearBGPNeighborResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearBGPNeighborResponse) Reset()         { *m = ClearBGPNeighborResponse{} }
func (m *ClearBGPNeighborResponse) String() string { return proto.CompactTextString(m) }
func (*ClearBGPNeighborResponse) ProtoMessage()    {}
func (*ClearBGPNeighborResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd5905d96efb1c39, []int{1}
}

func (m *ClearBGPNeighborResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearBGPNeighborResponse.Unmarshal(m, b)
}
func (m *ClearBGPNeighborResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearBGPNeighborResponse.Marshal(b, m, deterministic)
}
func (m *ClearBGPNeighborResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearBGPNeighborResponse.Merge(m, src)
}
func (m *ClearBGPNeighborResponse) XXX_Size() int {
	return xxx_messageInfo_ClearBGPNeighborResponse.Size(m)
}
func (m *ClearBGPNeighborResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearBGPNeighborResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ClearBGPNeighborResponse proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("gnoi.bgp.ClearBGPNeighborRequest_Mode", ClearBGPNeighborRequest_Mode_name, ClearBGPNeighborRequest_Mode_value)
	proto.RegisterType((*ClearBGPNeighborRequest)(nil), "gnoi.bgp.ClearBGPNeighborRequest")
	proto.RegisterType((*ClearBGPNeighborResponse)(nil), "gnoi.bgp.ClearBGPNeighborResponse")
}

func init() { proto.RegisterFile("bgp/bgp.proto", fileDescriptor_dd5905d96efb1c39) }

var fileDescriptor_dd5905d96efb1c39 = []byte{
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x50, 0xb1, 0x4e, 0xc3, 0x30,
	0x14, 0x6c, 0xda, 0x50, 0xc2, 0x93, 0x80, 0xc8, 0x0b, 0x51, 0xa6, 0x92, 0xa1, 0x2a, 0x12, 0x72,
	0x4a, 0xd9, 0x18, 0x90, 0x08, 0x88, 0xd2, 0x81, 0x52, 0x05, 0x36, 0x06, 0x14, 0x27, 0xc6, 0xb5,
	0x44, 0xfd, 0x4c, 0xec, 0x0c, 0xfc, 0x1e, 0x9f, 0xc0, 0x17, 0xa1, 0x24, 0xed, 0x02, 0x82, 0x2e,
	0x96, 0xef, 0xde, 0xdd, 0xbb, 0xd3, 0x83, 0x7d, 0x26, 0x74, 0xcc, 0x84, 0xa6, 0xba, 0x44, 0x8b,
	0xc4, 0x13, 0x0a, 0x25, 0x65, 0x42, 0x87, 0xa7, 0x42, 0xda, 0x65, 0xc5, 0x68, 0x8e, 0xab, 0x18,
	0x35, 0x57, 0x39, 0xaa, 0x57, 0x29, 0xe2, 0x7a, 0x1e, 0xdb, 0x0f, 0xcd, 0x4d, 0xfb, 0xb6, 0xbe,
	0xe8, 0xd3, 0x81, 0xa3, 0xeb, 0x37, 0x9e, 0x95, 0xc9, 0x74, 0x31, 0xe7, 0x52, 0x2c, 0x19, 0x96,
	0x29, 0x7f, 0xaf, 0xb8, 0xb1, 0x24, 0x80, 0xdd, 0xac, 0x28, 0x4a, 0x6e, 0x4c, 0xe0, 0x0c, 0x9c,
	0xd1, 0x5e, 0xba, 0x81, 0xe4, 0x04, 0xfc, 0x12, 0x2b, 0x2b, 0x95, 0x78, 0x91, 0xca, 0xd8, 0x4c,
	0xe5, 0x3c, 0xe8, 0x36, 0x92, 0xc3, 0x35, 0x3f, 0x5b, 0xd3, 0xe4, 0x02, 0xdc, 0x15, 0x16, 0x3c,
	0xe8, 0x0d, 0x9c, 0xd1, 0xc1, 0x64, 0x48, 0x37, 0x3d, 0xe9, 0x1f, 0xa9, 0xf4, 0x1e, 0x0b, 0x9e,
	0x36, 0x9e, 0x68, 0x08, 0x6e, 0x8d, 0x88, 0x07, 0xee, 0xe3, 0xc3, 0xed, 0x93, 0xdf, 0x21, 0x00,
	0xfd, 0xfa, 0x37, 0x9b, 0xfb, 0x4e, 0xcd, 0xde, 0x5d, 0xa5, 0x37, 0x7e, 0x37, 0x0a, 0x21, 0xf8,
	0xbd, 0xcd, 0x68, 0x54, 0x86, 0x4f, 0x18, 0xf4, 0x92, 0xe9, 0x82, 0x3c, 0x83, 0xff, 0x53, 0x42,
	0x8e, 0xb7, 0x96, 0x09, 0xa3, 0xff, 0x24, 0x6d, 0x42, 0xd4, 0x49, 0xbc, 0xaf, 0xcb, 0x9d, 0x31,
	0x3d, 0xa3, 0x63, 0xd6, 0x6f, 0xae, 0x7a, 0xfe, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x05, 0x3c, 0x32,
	0x2b, 0x9e, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BGPClient is the client API for BGP service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BGPClient interface {
	// ClearBGPNeighbor clears a BGP session.
	ClearBGPNeighbor(ctx context.Context, in *ClearBGPNeighborRequest, opts ...grpc.CallOption) (*ClearBGPNeighborResponse, error)
}

type bGPClient struct {
	cc *grpc.ClientConn
}

func NewBGPClient(cc *grpc.ClientConn) BGPClient {
	return &bGPClient{cc}
}

func (c *bGPClient) ClearBGPNeighbor(ctx context.Context, in *ClearBGPNeighborRequest, opts ...grpc.CallOption) (*ClearBGPNeighborResponse, error) {
	out := new(ClearBGPNeighborResponse)
	err := c.cc.Invoke(ctx, "/gnoi.bgp.BGP/ClearBGPNeighbor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BGPServer is the server API for BGP service.
type BGPServer interface {
	// ClearBGPNeighbor clears a BGP session.
	ClearBGPNeighbor(context.Context, *ClearBGPNeighborRequest) (*ClearBGPNeighborResponse, error)
}

func RegisterBGPServer(s *grpc.Server, srv BGPServer) {
	s.RegisterService(&_BGP_serviceDesc, srv)
}

func _BGP_ClearBGPNeighbor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearBGPNeighborRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BGPServer).ClearBGPNeighbor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gnoi.bgp.BGP/ClearBGPNeighbor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BGPServer).ClearBGPNeighbor(ctx, req.(*ClearBGPNeighborRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BGP_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gnoi.bgp.BGP",
	HandlerType: (*BGPServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClearBGPNeighbor",
			Handler:    _BGP_ClearBGPNeighbor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bgp/bgp.proto",
}
