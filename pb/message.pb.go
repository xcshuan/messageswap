// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: message.proto

package messageswap_pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Message_MessageType int32

const (
	Message_PUT_VALUE      Message_MessageType = 0
	Message_GET_VALUE      Message_MessageType = 1
	Message_ADD_PROVIDER   Message_MessageType = 2
	Message_GET_PROVIDERS  Message_MessageType = 3
	Message_FIND_NODE      Message_MessageType = 4
	Message_PING           Message_MessageType = 5
	Message_GET_PREFIX     Message_MessageType = 6
	Message_APPEND_VALUE   Message_MessageType = 7
	Message_META_SYNC      Message_MessageType = 8
	Message_User_Init_Req  Message_MessageType = 9
	Message_User_Init_Res  Message_MessageType = 10
	Message_User_NewKP_Req Message_MessageType = 11
	Message_New_User_Notif Message_MessageType = 12
	Message_Block_Meta     Message_MessageType = 13
	Message_Delete_Block   Message_MessageType = 14
	Message_Challenge      Message_MessageType = 15
	Message_Proof          Message_MessageType = 16
	Message_Proof_Sync     Message_MessageType = 17
	Message_Repair         Message_MessageType = 18
	Message_Query_Info     Message_MessageType = 19
	Message_Repair_Res     Message_MessageType = 20
	Message_Storage_Sync   Message_MessageType = 21
	Message_MetaInfo       Message_MessageType = 22
	Message_MetaBroadcast  Message_MessageType = 23
)

var Message_MessageType_name = map[int32]string{
	0:  "PUT_VALUE",
	1:  "GET_VALUE",
	2:  "ADD_PROVIDER",
	3:  "GET_PROVIDERS",
	4:  "FIND_NODE",
	5:  "PING",
	6:  "GET_PREFIX",
	7:  "APPEND_VALUE",
	8:  "META_SYNC",
	9:  "User_Init_Req",
	10: "User_Init_Res",
	11: "User_NewKP_Req",
	12: "New_User_Notif",
	13: "Block_Meta",
	14: "Delete_Block",
	15: "Challenge",
	16: "Proof",
	17: "Proof_Sync",
	18: "Repair",
	19: "Query_Info",
	20: "Repair_Res",
	21: "Storage_Sync",
	22: "MetaInfo",
	23: "MetaBroadcast",
}
var Message_MessageType_value = map[string]int32{
	"PUT_VALUE":      0,
	"GET_VALUE":      1,
	"ADD_PROVIDER":   2,
	"GET_PROVIDERS":  3,
	"FIND_NODE":      4,
	"PING":           5,
	"GET_PREFIX":     6,
	"APPEND_VALUE":   7,
	"META_SYNC":      8,
	"User_Init_Req":  9,
	"User_Init_Res":  10,
	"User_NewKP_Req": 11,
	"New_User_Notif": 12,
	"Block_Meta":     13,
	"Delete_Block":   14,
	"Challenge":      15,
	"Proof":          16,
	"Proof_Sync":     17,
	"Repair":         18,
	"Query_Info":     19,
	"Repair_Res":     20,
	"Storage_Sync":   21,
	"MetaInfo":       22,
	"MetaBroadcast":  23,
}

func (x Message_MessageType) String() string {
	return proto.EnumName(Message_MessageType_name, int32(x))
}
func (Message_MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_message_59eb3dadee47de0f, []int{0, 0}
}

// `protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf
// --gogo_out=. *.proto`
type Message struct {
	Type                 Message_MessageType `protobuf:"varint,1,opt,name=type,proto3,enum=messageswap.pb.Message_MessageType" json:"type,omitempty"`
	Key                  Message_Key         `protobuf:"bytes,2,opt,name=key" json:"key"`
	Value                [][]byte            `protobuf:"bytes,3,rep,name=value" json:"value,omitempty"`
	Timestamp            string              `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_59eb3dadee47de0f, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetType() Message_MessageType {
	if m != nil {
		return m.Type
	}
	return Message_PUT_VALUE
}

func (m *Message) GetKey() Message_Key {
	if m != nil {
		return m.Key
	}
	return Message_Key{}
}

func (m *Message) GetValue() [][]byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Message) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type Message_Key struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Arguments            [][]byte `protobuf:"bytes,2,rep,name=arguments" json:"arguments,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message_Key) Reset()         { *m = Message_Key{} }
func (m *Message_Key) String() string { return proto.CompactTextString(m) }
func (*Message_Key) ProtoMessage()    {}
func (*Message_Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_59eb3dadee47de0f, []int{0, 0}
}
func (m *Message_Key) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Key.Unmarshal(m, b)
}
func (m *Message_Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Key.Marshal(b, m, deterministic)
}
func (dst *Message_Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Key.Merge(dst, src)
}
func (m *Message_Key) XXX_Size() int {
	return xxx_messageInfo_Message_Key.Size(m)
}
func (m *Message_Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Key proto.InternalMessageInfo

func (m *Message_Key) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Message_Key) GetArguments() [][]byte {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "messageswap.pb.Message")
	proto.RegisterType((*Message_Key)(nil), "messageswap.pb.Message.Key")
	proto.RegisterEnum("messageswap.pb.Message_MessageType", Message_MessageType_name, Message_MessageType_value)
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_message_59eb3dadee47de0f) }

var fileDescriptor_message_59eb3dadee47de0f = []byte{
	// 485 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xc1, 0x6e, 0xda, 0x40,
	0x14, 0x45, 0x03, 0x36, 0x04, 0x3f, 0xc0, 0x9d, 0x4c, 0xd3, 0x16, 0xa5, 0x95, 0x8a, 0xd2, 0x0d,
	0x9b, 0x12, 0x29, 0x2c, 0xba, 0x86, 0xd8, 0x89, 0x2c, 0x8a, 0xe3, 0x1a, 0x88, 0xda, 0xd5, 0x68,
	0xa0, 0x0f, 0x07, 0x05, 0x3c, 0xae, 0x3d, 0x2e, 0xf2, 0x77, 0xf4, 0xa7, 0xfa, 0x15, 0xfd, 0x8c,
	0xae, 0xab, 0x19, 0x83, 0xd2, 0x56, 0xca, 0xca, 0x73, 0xef, 0xdc, 0x73, 0xdf, 0xb3, 0x34, 0xd0,
	0xde, 0x62, 0x96, 0xf1, 0x08, 0xfb, 0x49, 0x2a, 0xa4, 0xa0, 0xf6, 0x5e, 0x66, 0x3b, 0x9e, 0xf4,
	0x93, 0xc5, 0xd9, 0xfb, 0x68, 0x2d, 0xef, 0xf3, 0x45, 0x7f, 0x29, 0xb6, 0x17, 0x91, 0x88, 0xc4,
	0x85, 0x8e, 0x2d, 0xf2, 0x95, 0x56, 0x5a, 0xe8, 0x53, 0x89, 0x9f, 0xff, 0x36, 0xe1, 0x78, 0x52,
	0x36, 0xd0, 0x0f, 0x60, 0xca, 0x22, 0xc1, 0x4e, 0xa5, 0x5b, 0xe9, 0xd9, 0x97, 0xef, 0xfa, 0xff,
	0x36, 0xf7, 0xf7, 0xb1, 0xc3, 0x77, 0x56, 0x24, 0x18, 0x6a, 0x80, 0x0e, 0xc0, 0x78, 0xc0, 0xa2,
	0x53, 0xed, 0x56, 0x7a, 0xcd, 0xcb, 0xd7, 0x4f, 0x71, 0x63, 0x2c, 0x46, 0xe6, 0xcf, 0x5f, 0x6f,
	0x8f, 0x42, 0x95, 0xa6, 0xa7, 0x50, 0xfb, 0xce, 0x37, 0x39, 0x76, 0x8c, 0xae, 0xd1, 0x6b, 0x85,
	0xa5, 0xa0, 0x6f, 0xc0, 0x92, 0xeb, 0x2d, 0x66, 0x92, 0x6f, 0x93, 0x8e, 0xd9, 0xad, 0xf4, 0xac,
	0xf0, 0xd1, 0x38, 0x1b, 0x80, 0x31, 0xc6, 0x82, 0xda, 0x50, 0xf5, 0x1c, 0xbd, 0xa6, 0x15, 0x56,
	0x3d, 0x47, 0x41, 0x3c, 0x8d, 0xf2, 0x2d, 0xc6, 0x32, 0xeb, 0x54, 0x75, 0xdd, 0xa3, 0x71, 0xfe,
	0xc3, 0x80, 0xe6, 0x5f, 0x3b, 0xd3, 0x36, 0x58, 0xc1, 0x7c, 0xc6, 0xee, 0x86, 0x1f, 0xe7, 0x2e,
	0x39, 0x52, 0xf2, 0xc6, 0x3d, 0xc8, 0x0a, 0x25, 0xd0, 0x1a, 0x3a, 0x0e, 0x0b, 0xc2, 0xdb, 0x3b,
	0xcf, 0x71, 0x43, 0x52, 0xa5, 0x27, 0xd0, 0x56, 0x81, 0x83, 0x33, 0x25, 0x86, 0x62, 0xae, 0x3d,
	0xdf, 0x61, 0xfe, 0xad, 0xe3, 0x12, 0x93, 0x36, 0xc0, 0x0c, 0x3c, 0xff, 0x86, 0xd4, 0xa8, 0x0d,
	0x50, 0x66, 0xdd, 0x6b, 0xef, 0x33, 0xa9, 0xeb, 0xb6, 0x20, 0x70, 0x7d, 0x67, 0xdf, 0x7f, 0xac,
	0xd0, 0x89, 0x3b, 0x1b, 0xb2, 0xe9, 0x17, 0xff, 0x8a, 0x34, 0x54, 0xf9, 0x3c, 0xc3, 0x94, 0x79,
	0xf1, 0x5a, 0xb2, 0x10, 0xbf, 0x11, 0xeb, 0x7f, 0x2b, 0x23, 0x40, 0x29, 0xd8, 0xda, 0xf2, 0x71,
	0x37, 0x0e, 0x74, 0xac, 0xa9, 0x3c, 0x1f, 0x77, 0xac, 0xf4, 0x85, 0x5c, 0xaf, 0x48, 0x4b, 0x8d,
	0x1f, 0x6d, 0xc4, 0xf2, 0x81, 0x4d, 0x50, 0x72, 0xd2, 0x56, 0xe3, 0x1d, 0xdc, 0xa0, 0x44, 0xa6,
	0x6d, 0x62, 0xab, 0xf1, 0x57, 0xf7, 0x7c, 0xb3, 0xc1, 0x38, 0x42, 0xf2, 0x8c, 0x5a, 0x50, 0x0b,
	0x52, 0x21, 0x56, 0x84, 0x28, 0x56, 0x1f, 0xd9, 0xb4, 0x88, 0x97, 0xe4, 0x84, 0x02, 0xd4, 0x43,
	0x4c, 0xf8, 0x3a, 0x25, 0x54, 0xdd, 0x7d, 0xca, 0x31, 0x2d, 0x98, 0x17, 0xaf, 0x04, 0x79, 0xae,
	0x74, 0x79, 0xa7, 0xf7, 0x3b, 0x55, 0x73, 0xa6, 0x52, 0xa4, 0x3c, 0xc2, 0x92, 0x7e, 0x41, 0x5b,
	0xd0, 0x50, 0x3b, 0xe8, 0xfc, 0x4b, 0xf5, 0x4b, 0x4a, 0x8d, 0x52, 0xc1, 0xbf, 0x2e, 0x79, 0x26,
	0xc9, 0xab, 0x45, 0x5d, 0xbf, 0xbf, 0xc1, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x5e, 0xbf,
	0x74, 0xcf, 0x02, 0x00, 0x00,
}
