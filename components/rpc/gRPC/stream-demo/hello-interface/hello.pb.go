// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hello.proto

package hello_interface

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type String struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_hello_89b12a75f1fdbbb0, []int{0}
}
func (m *String) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_String.Unmarshal(m, b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_String.Marshal(b, m, deterministic)
}
func (dst *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(dst, src)
}
func (m *String) XXX_Size() int {
	return xxx_messageInfo_String.Size(m)
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*String)(nil), "HelloService.String")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HelloServiceClient is the client API for HelloService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HelloServiceClient interface {
	Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Channel(ctx context.Context, opts ...grpc.CallOption) (HelloService_ChannelClient, error)
}

type helloServiceClient struct {
	cc *grpc.ClientConn
}

func NewHelloServiceClient(cc *grpc.ClientConn) HelloServiceClient {
	return &helloServiceClient{cc}
}

func (c *helloServiceClient) Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/HelloService.HelloService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *helloServiceClient) Channel(ctx context.Context, opts ...grpc.CallOption) (HelloService_ChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &_HelloService_serviceDesc.Streams[0], "/HelloService.HelloService/Channel", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloServiceChannelClient{stream}
	return x, nil
}

type HelloService_ChannelClient interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ClientStream
}

type helloServiceChannelClient struct {
	grpc.ClientStream
}

func (x *helloServiceChannelClient) Send(m *String) error {
	return x.ClientStream.SendMsg(m)
}

func (x *helloServiceChannelClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HelloServiceServer is the server API for HelloService service.
type HelloServiceServer interface {
	Hello(context.Context, *String) (*String, error)
	Channel(HelloService_ChannelServer) error
}

func RegisterHelloServiceServer(s *grpc.Server, srv HelloServiceServer) {
	s.RegisterService(&_HelloService_serviceDesc, srv)
}

func _HelloService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HelloService.HelloService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServiceServer).Hello(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _HelloService_Channel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServiceServer).Channel(&helloServiceChannelServer{stream})
}

type HelloService_ChannelServer interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ServerStream
}

type helloServiceChannelServer struct {
	grpc.ServerStream
}

func (x *helloServiceChannelServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func (x *helloServiceChannelServer) Recv() (*String, error) {
	m := new(String)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _HelloService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "HelloService.HelloService",
	HandlerType: (*HelloServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _HelloService_Hello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Channel",
			Handler:       _HelloService_Channel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello.proto",
}

func init() { proto.RegisterFile("hello.proto", fileDescriptor_hello_89b12a75f1fdbbb0) }

var fileDescriptor_hello_89b12a75f1fdbbb0 = []byte{
	// 121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xf1, 0x00, 0x71, 0x82, 0x53, 0x8b, 0xca,
	0x32, 0x93, 0x53, 0x95, 0xe4, 0xb8, 0xd8, 0x82, 0x4b, 0x8a, 0x32, 0xf3, 0xd2, 0x85, 0x44, 0xb8,
	0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0xa3,
	0x3a, 0x2e, 0x14, 0xf5, 0x42, 0xc6, 0x5c, 0xac, 0x60, 0xbe, 0x90, 0x88, 0x1e, 0xb2, 0xb8, 0x1e,
	0xc4, 0x10, 0x29, 0xac, 0xa2, 0x42, 0x96, 0x5c, 0xec, 0xce, 0x19, 0x89, 0x79, 0x79, 0xa9, 0x39,
	0xa4, 0x68, 0xd3, 0x60, 0x34, 0x60, 0x4c, 0x62, 0x03, 0x3b, 0xda, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0x9a, 0x51, 0x94, 0xb1, 0xc3, 0x00, 0x00, 0x00,
}
