// Code generated by protoc-gen-go.
// source: authenticator.proto
// DO NOT EDIT!

/*
Package specs is a generated protocol buffer package.

It is generated from these files:
	authenticator.proto
	common.proto

It has these top-level messages:
	User
	ACL
	Empty
*/
package specs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ACLAction int32

const (
	ACL_read  ACLAction = 0
	ACL_write ACLAction = 1
)

var ACLAction_name = map[int32]string{
	0: "read",
	1: "write",
}
var ACLAction_value = map[string]int32{
	"read":  0,
	"write": 1,
}

func (x ACLAction) String() string {
	return proto.EnumName(ACLAction_name, int32(x))
}
func (ACLAction) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type User struct {
	Error    ErrorCode `protobuf:"varint,1,opt,name=error,enum=specs.ErrorCode" json:"error,omitempty"`
	UserID   string    `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	Account  string    `protobuf:"bytes,3,opt,name=account" json:"account,omitempty"`
	Email    string    `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	Password string    `protobuf:"bytes,5,opt,name=password" json:"password,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *User) GetError() ErrorCode {
	if m != nil {
		return m.Error
	}
	return ErrorCode_OK
}

func (m *User) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *User) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type ACL struct {
	Error    ErrorCode   `protobuf:"varint,1,opt,name=error,enum=specs.ErrorCode" json:"error,omitempty"`
	UserID   string      `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	Resource string      `protobuf:"bytes,3,opt,name=resource" json:"resource,omitempty"`
	Actions  []ACLAction `protobuf:"varint,4,rep,packed,name=actions,enum=specs.ACLAction" json:"actions,omitempty"`
	Allow    bool        `protobuf:"varint,5,opt,name=allow" json:"allow,omitempty"`
}

func (m *ACL) Reset()                    { *m = ACL{} }
func (m *ACL) String() string            { return proto.CompactTextString(m) }
func (*ACL) ProtoMessage()               {}
func (*ACL) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ACL) GetError() ErrorCode {
	if m != nil {
		return m.Error
	}
	return ErrorCode_OK
}

func (m *ACL) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ACL) GetResource() string {
	if m != nil {
		return m.Resource
	}
	return ""
}

func (m *ACL) GetActions() []ACLAction {
	if m != nil {
		return m.Actions
	}
	return nil
}

func (m *ACL) GetAllow() bool {
	if m != nil {
		return m.Allow
	}
	return false
}

func init() {
	proto.RegisterType((*User)(nil), "specs.User")
	proto.RegisterType((*ACL)(nil), "specs.ACL")
	proto.RegisterEnum("specs.ACLAction", ACLAction_name, ACLAction_value)
}

func init() { proto.RegisterFile("authenticator.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 255 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x90, 0x5d, 0x4a, 0xc4, 0x30,
	0x14, 0x85, 0x8d, 0xfd, 0xb1, 0x73, 0x91, 0xa1, 0x46, 0x91, 0x50, 0x10, 0x4a, 0x1f, 0xa4, 0x20,
	0xf4, 0x61, 0x5c, 0xc1, 0x50, 0x7d, 0x10, 0xe6, 0xa9, 0xe0, 0x02, 0x62, 0x1a, 0xb0, 0xd0, 0xf6,
	0x96, 0x9b, 0x94, 0xee, 0xc3, 0x2d, 0xb9, 0x31, 0x69, 0xd2, 0x99, 0x0d, 0xf8, 0xf8, 0x9d, 0x73,
	0x08, 0x5f, 0x2e, 0xdc, 0xcb, 0xd9, 0x7e, 0xeb, 0xd1, 0x76, 0x4a, 0x5a, 0xa4, 0x6a, 0x22, 0xb4,
	0xc8, 0x23, 0x33, 0x69, 0x65, 0xb2, 0x5b, 0x85, 0xc3, 0x80, 0xa3, 0x0f, 0x8b, 0x1f, 0x06, 0xe1,
	0xa7, 0xd1, 0xc4, 0x9f, 0x21, 0xd2, 0x44, 0x48, 0x82, 0xe5, 0xac, 0xdc, 0x1f, 0xd2, 0xca, 0xad,
	0xab, 0xf7, 0x35, 0xab, 0xb1, 0xd5, 0x8d, 0xaf, 0xf9, 0x23, 0xc4, 0xb3, 0xd1, 0xf4, 0xf1, 0x26,
	0xae, 0x73, 0x56, 0xee, 0x9a, 0x8d, 0xb8, 0x80, 0x1b, 0xa9, 0x14, 0xce, 0xa3, 0x15, 0x81, 0x2b,
	0xce, 0xc8, 0x1f, 0x20, 0xd2, 0x83, 0xec, 0x7a, 0x11, 0xba, 0xdc, 0x03, 0xcf, 0x20, 0x99, 0xa4,
	0x31, 0x0b, 0x52, 0x2b, 0x22, 0x57, 0x5c, 0xb8, 0xf8, 0x65, 0x10, 0x1c, 0xeb, 0xd3, 0xbf, 0x9d,
	0x32, 0x48, 0x48, 0x1b, 0x9c, 0x49, 0xe9, 0x4d, 0xea, 0xc2, 0xfc, 0x65, 0xf5, 0xb5, 0x1d, 0x8e,
	0x46, 0x84, 0x79, 0x50, 0xee, 0x0f, 0x77, 0xdb, 0xeb, 0xc7, 0xfa, 0x54, 0xf9, 0xa6, 0x39, 0x2f,
	0xd6, 0x2f, 0xc8, 0xbe, 0xc7, 0xc5, 0x99, 0x26, 0x8d, 0x87, 0xe2, 0x09, 0x62, 0x3f, 0xe0, 0x09,
	0x84, 0xa4, 0x65, 0x9b, 0x5e, 0xf1, 0x1d, 0x44, 0x0b, 0x75, 0x56, 0xa7, 0xec, 0x2b, 0x76, 0x17,
	0x7e, 0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xa1, 0x15, 0x5e, 0x5a, 0x8d, 0x01, 0x00, 0x00,
}
