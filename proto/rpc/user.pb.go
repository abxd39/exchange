// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc/user.proto

package g2u

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

type CommonErrResponse struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err" json:"err,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommonErrResponse) Reset()         { *m = CommonErrResponse{} }
func (m *CommonErrResponse) String() string { return proto.CompactTextString(m) }
func (*CommonErrResponse) ProtoMessage()    {}
func (*CommonErrResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{0}
}
func (m *CommonErrResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonErrResponse.Unmarshal(m, b)
}
func (m *CommonErrResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonErrResponse.Marshal(b, m, deterministic)
}
func (dst *CommonErrResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonErrResponse.Merge(dst, src)
}
func (m *CommonErrResponse) XXX_Size() int {
	return xxx_messageInfo_CommonErrResponse.Size(m)
}
func (m *CommonErrResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonErrResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommonErrResponse proto.InternalMessageInfo

func (m *CommonErrResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

func (m *CommonErrResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type HelloRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{1}
}
func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloRequest.Unmarshal(m, b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
}
func (dst *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(dst, src)
}
func (m *HelloRequest) XXX_Size() int {
	return xxx_messageInfo_HelloRequest.Size(m)
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type HelloResponse struct {
	Greeting             string   `protobuf:"bytes,2,opt,name=greeting" json:"greeting,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloResponse) Reset()         { *m = HelloResponse{} }
func (m *HelloResponse) String() string { return proto.CompactTextString(m) }
func (*HelloResponse) ProtoMessage()    {}
func (*HelloResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{2}
}
func (m *HelloResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloResponse.Unmarshal(m, b)
}
func (m *HelloResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloResponse.Marshal(b, m, deterministic)
}
func (dst *HelloResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloResponse.Merge(dst, src)
}
func (m *HelloResponse) XXX_Size() int {
	return xxx_messageInfo_HelloResponse.Size(m)
}
func (m *HelloResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HelloResponse proto.InternalMessageInfo

func (m *HelloResponse) GetGreeting() string {
	if m != nil {
		return m.Greeting
	}
	return ""
}

type RegisterRequest struct {
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone,omitempty"`
	Pwd                  string   `protobuf:"bytes,2,opt,name=pwd" json:"pwd,omitempty"`
	InviteCode           string   `protobuf:"bytes,4,opt,name=invite_code,json=inviteCode" json:"invite_code,omitempty"`
	Country              int32    `protobuf:"varint,5,opt,name=country" json:"country,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{3}
}
func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (dst *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(dst, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RegisterRequest) GetPwd() string {
	if m != nil {
		return m.Pwd
	}
	return ""
}

func (m *RegisterRequest) GetInviteCode() string {
	if m != nil {
		return m.InviteCode
	}
	return ""
}

func (m *RegisterRequest) GetCountry() int32 {
	if m != nil {
		return m.Country
	}
	return 0
}

type RegisterResponse struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{4}
}
func (m *RegisterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResponse.Unmarshal(m, b)
}
func (m *RegisterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResponse.Marshal(b, m, deterministic)
}
func (dst *RegisterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResponse.Merge(dst, src)
}
func (m *RegisterResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterResponse.Size(m)
}
func (m *RegisterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

type LoginRequest struct {
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone,omitempty"`
	Pwd                  string   `protobuf:"bytes,2,opt,name=pwd" json:"pwd,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{5}
}
func (m *LoginRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginRequest.Unmarshal(m, b)
}
func (m *LoginRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginRequest.Marshal(b, m, deterministic)
}
func (dst *LoginRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginRequest.Merge(dst, src)
}
func (m *LoginRequest) XXX_Size() int {
	return xxx_messageInfo_LoginRequest.Size(m)
}
func (m *LoginRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoginRequest proto.InternalMessageInfo

func (m *LoginRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *LoginRequest) GetPwd() string {
	if m != nil {
		return m.Pwd
	}
	return ""
}

type LoginResponse struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{6}
}
func (m *LoginResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResponse.Unmarshal(m, b)
}
func (m *LoginResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResponse.Marshal(b, m, deterministic)
}
func (dst *LoginResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResponse.Merge(dst, src)
}
func (m *LoginResponse) XXX_Size() int {
	return xxx_messageInfo_LoginResponse.Size(m)
}
func (m *LoginResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResponse proto.InternalMessageInfo

func (m *LoginResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

type ForgetRequest struct {
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForgetRequest) Reset()         { *m = ForgetRequest{} }
func (m *ForgetRequest) String() string { return proto.CompactTextString(m) }
func (*ForgetRequest) ProtoMessage()    {}
func (*ForgetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{7}
}
func (m *ForgetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForgetRequest.Unmarshal(m, b)
}
func (m *ForgetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForgetRequest.Marshal(b, m, deterministic)
}
func (dst *ForgetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForgetRequest.Merge(dst, src)
}
func (m *ForgetRequest) XXX_Size() int {
	return xxx_messageInfo_ForgetRequest.Size(m)
}
func (m *ForgetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ForgetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ForgetRequest proto.InternalMessageInfo

func (m *ForgetRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type ForgetResponse struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err" json:"err,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone,omitempty"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForgetResponse) Reset()         { *m = ForgetResponse{} }
func (m *ForgetResponse) String() string { return proto.CompactTextString(m) }
func (*ForgetResponse) ProtoMessage()    {}
func (*ForgetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{8}
}
func (m *ForgetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForgetResponse.Unmarshal(m, b)
}
func (m *ForgetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForgetResponse.Marshal(b, m, deterministic)
}
func (dst *ForgetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForgetResponse.Merge(dst, src)
}
func (m *ForgetResponse) XXX_Size() int {
	return xxx_messageInfo_ForgetResponse.Size(m)
}
func (m *ForgetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ForgetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ForgetResponse proto.InternalMessageInfo

func (m *ForgetResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

func (m *ForgetResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ForgetResponse) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ForgetResponse) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type SecurityRequest struct {
	Phone                string   `protobuf:"bytes,1,opt,name=phone" json:"phone,omitempty"`
	PhoneAuthCode        string   `protobuf:"bytes,2,opt,name=phone_auth_code,json=phoneAuthCode" json:"phone_auth_code,omitempty"`
	EmailAuthCode        string   `protobuf:"bytes,3,opt,name=email_auth_code,json=emailAuthCode" json:"email_auth_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SecurityRequest) Reset()         { *m = SecurityRequest{} }
func (m *SecurityRequest) String() string { return proto.CompactTextString(m) }
func (*SecurityRequest) ProtoMessage()    {}
func (*SecurityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{9}
}
func (m *SecurityRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SecurityRequest.Unmarshal(m, b)
}
func (m *SecurityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SecurityRequest.Marshal(b, m, deterministic)
}
func (dst *SecurityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecurityRequest.Merge(dst, src)
}
func (m *SecurityRequest) XXX_Size() int {
	return xxx_messageInfo_SecurityRequest.Size(m)
}
func (m *SecurityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SecurityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SecurityRequest proto.InternalMessageInfo

func (m *SecurityRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *SecurityRequest) GetPhoneAuthCode() string {
	if m != nil {
		return m.PhoneAuthCode
	}
	return ""
}

func (m *SecurityRequest) GetEmailAuthCode() string {
	if m != nil {
		return m.EmailAuthCode
	}
	return ""
}

type SecurityResponse struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err" json:"err,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	SecurityKey          []byte   `protobuf:"bytes,3,opt,name=security_key,json=securityKey,proto3" json:"security_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SecurityResponse) Reset()         { *m = SecurityResponse{} }
func (m *SecurityResponse) String() string { return proto.CompactTextString(m) }
func (*SecurityResponse) ProtoMessage()    {}
func (*SecurityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{10}
}
func (m *SecurityResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SecurityResponse.Unmarshal(m, b)
}
func (m *SecurityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SecurityResponse.Marshal(b, m, deterministic)
}
func (dst *SecurityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecurityResponse.Merge(dst, src)
}
func (m *SecurityResponse) XXX_Size() int {
	return xxx_messageInfo_SecurityResponse.Size(m)
}
func (m *SecurityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SecurityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SecurityResponse proto.InternalMessageInfo

func (m *SecurityResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

func (m *SecurityResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *SecurityResponse) GetSecurityKey() []byte {
	if m != nil {
		return m.SecurityKey
	}
	return nil
}

type NoticeListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoticeListRequest) Reset()         { *m = NoticeListRequest{} }
func (m *NoticeListRequest) String() string { return proto.CompactTextString(m) }
func (*NoticeListRequest) ProtoMessage()    {}
func (*NoticeListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{11}
}
func (m *NoticeListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeListRequest.Unmarshal(m, b)
}
func (m *NoticeListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeListRequest.Marshal(b, m, deterministic)
}
func (dst *NoticeListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeListRequest.Merge(dst, src)
}
func (m *NoticeListRequest) XXX_Size() int {
	return xxx_messageInfo_NoticeListRequest.Size(m)
}
func (m *NoticeListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeListRequest proto.InternalMessageInfo

type NoticeListResponse struct {
	Notice               []*NoticeListResponse_Notice `protobuf:"bytes,1,rep,name=notice" json:"notice,omitempty"`
	Err                  int32                        `protobuf:"varint,2,opt,name=err" json:"err,omitempty"`
	Message              string                       `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *NoticeListResponse) Reset()         { *m = NoticeListResponse{} }
func (m *NoticeListResponse) String() string { return proto.CompactTextString(m) }
func (*NoticeListResponse) ProtoMessage()    {}
func (*NoticeListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{12}
}
func (m *NoticeListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeListResponse.Unmarshal(m, b)
}
func (m *NoticeListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeListResponse.Marshal(b, m, deterministic)
}
func (dst *NoticeListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeListResponse.Merge(dst, src)
}
func (m *NoticeListResponse) XXX_Size() int {
	return xxx_messageInfo_NoticeListResponse.Size(m)
}
func (m *NoticeListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeListResponse proto.InternalMessageInfo

func (m *NoticeListResponse) GetNotice() []*NoticeListResponse_Notice {
	if m != nil {
		return m.Notice
	}
	return nil
}

func (m *NoticeListResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

func (m *NoticeListResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type NoticeListResponse_Notice struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	Description          string   `protobuf:"bytes,3,opt,name=Description" json:"Description,omitempty"`
	CreateDateTime       string   `protobuf:"bytes,4,opt,name=createDateTime" json:"createDateTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoticeListResponse_Notice) Reset()         { *m = NoticeListResponse_Notice{} }
func (m *NoticeListResponse_Notice) String() string { return proto.CompactTextString(m) }
func (*NoticeListResponse_Notice) ProtoMessage()    {}
func (*NoticeListResponse_Notice) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{12, 0}
}
func (m *NoticeListResponse_Notice) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeListResponse_Notice.Unmarshal(m, b)
}
func (m *NoticeListResponse_Notice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeListResponse_Notice.Marshal(b, m, deterministic)
}
func (dst *NoticeListResponse_Notice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeListResponse_Notice.Merge(dst, src)
}
func (m *NoticeListResponse_Notice) XXX_Size() int {
	return xxx_messageInfo_NoticeListResponse_Notice.Size(m)
}
func (m *NoticeListResponse_Notice) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeListResponse_Notice.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeListResponse_Notice proto.InternalMessageInfo

func (m *NoticeListResponse_Notice) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NoticeListResponse_Notice) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *NoticeListResponse_Notice) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *NoticeListResponse_Notice) GetCreateDateTime() string {
	if m != nil {
		return m.CreateDateTime
	}
	return ""
}

type NoticeDetailRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoticeDetailRequest) Reset()         { *m = NoticeDetailRequest{} }
func (m *NoticeDetailRequest) String() string { return proto.CompactTextString(m) }
func (*NoticeDetailRequest) ProtoMessage()    {}
func (*NoticeDetailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{13}
}
func (m *NoticeDetailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeDetailRequest.Unmarshal(m, b)
}
func (m *NoticeDetailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeDetailRequest.Marshal(b, m, deterministic)
}
func (dst *NoticeDetailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeDetailRequest.Merge(dst, src)
}
func (m *NoticeDetailRequest) XXX_Size() int {
	return xxx_messageInfo_NoticeDetailRequest.Size(m)
}
func (m *NoticeDetailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeDetailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeDetailRequest proto.InternalMessageInfo

func (m *NoticeDetailRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type NoticeDetailResponse struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err" json:"err,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Id                   int32    `protobuf:"varint,3,opt,name=id" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,4,opt,name=title" json:"title,omitempty"`
	Description          string   `protobuf:"bytes,5,opt,name=description" json:"description,omitempty"`
	Content              string   `protobuf:"bytes,6,opt,name=content" json:"content,omitempty"`
	Covers               []byte   `protobuf:"bytes,7,opt,name=covers,proto3" json:"covers,omitempty"`
	ContentImages        []byte   `protobuf:"bytes,8,opt,name=content_images,json=contentImages,proto3" json:"content_images,omitempty"`
	Type                 int32    `protobuf:"varint,9,opt,name=type" json:"type,omitempty"`
	TypeName             string   `protobuf:"bytes,10,opt,name=type_name,json=typeName" json:"type_name,omitempty"`
	Author               string   `protobuf:"bytes,11,opt,name=author" json:"author,omitempty"`
	Weight               int32    `protobuf:"varint,12,opt,name=weight" json:"weight,omitempty"`
	Shares               int32    `protobuf:"varint,13,opt,name=shares" json:"shares,omitempty"`
	Hits                 int32    `protobuf:"varint,14,opt,name=hits" json:"hits,omitempty"`
	Comments             int32    `protobuf:"varint,15,opt,name=comments" json:"comments,omitempty"`
	DisplayMark          bool     `protobuf:"varint,16,opt,name=DisplayMark" json:"DisplayMark,omitempty"`
	CreateTime           string   `protobuf:"bytes,17,opt,name=create_time,json=createTime" json:"create_time,omitempty"`
	UpdateTime           string   `protobuf:"bytes,18,opt,name=update_time,json=updateTime" json:"update_time,omitempty"`
	AdminId              int32    `protobuf:"varint,19,opt,name=adminId" json:"adminId,omitempty"`
	AdminNickname        string   `protobuf:"bytes,20,opt,name=admin_nickname,json=adminNickname" json:"admin_nickname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoticeDetailResponse) Reset()         { *m = NoticeDetailResponse{} }
func (m *NoticeDetailResponse) String() string { return proto.CompactTextString(m) }
func (*NoticeDetailResponse) ProtoMessage()    {}
func (*NoticeDetailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_e94a6429fb0542dc, []int{14}
}
func (m *NoticeDetailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoticeDetailResponse.Unmarshal(m, b)
}
func (m *NoticeDetailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoticeDetailResponse.Marshal(b, m, deterministic)
}
func (dst *NoticeDetailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoticeDetailResponse.Merge(dst, src)
}
func (m *NoticeDetailResponse) XXX_Size() int {
	return xxx_messageInfo_NoticeDetailResponse.Size(m)
}
func (m *NoticeDetailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NoticeDetailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NoticeDetailResponse proto.InternalMessageInfo

func (m *NoticeDetailResponse) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

func (m *NoticeDetailResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *NoticeDetailResponse) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NoticeDetailResponse) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *NoticeDetailResponse) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *NoticeDetailResponse) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *NoticeDetailResponse) GetCovers() []byte {
	if m != nil {
		return m.Covers
	}
	return nil
}

func (m *NoticeDetailResponse) GetContentImages() []byte {
	if m != nil {
		return m.ContentImages
	}
	return nil
}

func (m *NoticeDetailResponse) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *NoticeDetailResponse) GetTypeName() string {
	if m != nil {
		return m.TypeName
	}
	return ""
}

func (m *NoticeDetailResponse) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *NoticeDetailResponse) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *NoticeDetailResponse) GetShares() int32 {
	if m != nil {
		return m.Shares
	}
	return 0
}

func (m *NoticeDetailResponse) GetHits() int32 {
	if m != nil {
		return m.Hits
	}
	return 0
}

func (m *NoticeDetailResponse) GetComments() int32 {
	if m != nil {
		return m.Comments
	}
	return 0
}

func (m *NoticeDetailResponse) GetDisplayMark() bool {
	if m != nil {
		return m.DisplayMark
	}
	return false
}

func (m *NoticeDetailResponse) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

func (m *NoticeDetailResponse) GetUpdateTime() string {
	if m != nil {
		return m.UpdateTime
	}
	return ""
}

func (m *NoticeDetailResponse) GetAdminId() int32 {
	if m != nil {
		return m.AdminId
	}
	return 0
}

func (m *NoticeDetailResponse) GetAdminNickname() string {
	if m != nil {
		return m.AdminNickname
	}
	return ""
}

func init() {
	proto.RegisterType((*CommonErrResponse)(nil), "g2u.CommonErrResponse")
	proto.RegisterType((*HelloRequest)(nil), "g2u.HelloRequest")
	proto.RegisterType((*HelloResponse)(nil), "g2u.HelloResponse")
	proto.RegisterType((*RegisterRequest)(nil), "g2u.RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "g2u.RegisterResponse")
	proto.RegisterType((*LoginRequest)(nil), "g2u.LoginRequest")
	proto.RegisterType((*LoginResponse)(nil), "g2u.LoginResponse")
	proto.RegisterType((*ForgetRequest)(nil), "g2u.ForgetRequest")
	proto.RegisterType((*ForgetResponse)(nil), "g2u.ForgetResponse")
	proto.RegisterType((*SecurityRequest)(nil), "g2u.SecurityRequest")
	proto.RegisterType((*SecurityResponse)(nil), "g2u.SecurityResponse")
	proto.RegisterType((*NoticeListRequest)(nil), "g2u.NoticeListRequest")
	proto.RegisterType((*NoticeListResponse)(nil), "g2u.NoticeListResponse")
	proto.RegisterType((*NoticeListResponse_Notice)(nil), "g2u.NoticeListResponse.Notice")
	proto.RegisterType((*NoticeDetailRequest)(nil), "g2u.NoticeDetailRequest")
	proto.RegisterType((*NoticeDetailResponse)(nil), "g2u.NoticeDetailResponse")
}

func init() { proto.RegisterFile("rpc/user.proto", fileDescriptor_user_e94a6429fb0542dc) }

var fileDescriptor_user_e94a6429fb0542dc = []byte{
	// 837 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x56, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0x8e, 0xe3, 0xd8, 0x89, 0x8f, 0xd7, 0x76, 0x32, 0x31, 0x65, 0x30, 0x12, 0xb8, 0x23, 0x12,
	0x45, 0x42, 0x0a, 0x95, 0x91, 0x22, 0x71, 0x81, 0x2a, 0xd4, 0x14, 0xa8, 0x28, 0x11, 0x5a, 0xe0,
	0xda, 0x1a, 0x76, 0x8f, 0xd6, 0xa3, 0x78, 0x77, 0x96, 0x99, 0xd9, 0x1a, 0xbf, 0x13, 0x37, 0xbc,
	0x1a, 0x4f, 0x50, 0xcd, 0xcf, 0xda, 0x6b, 0xa7, 0xa9, 0xd4, 0x5e, 0x65, 0xce, 0x37, 0xdf, 0xf9,
	0xf7, 0x7c, 0x1b, 0x18, 0xaa, 0x32, 0xf9, 0xa6, 0xd2, 0xa8, 0xae, 0x4b, 0x25, 0x8d, 0x24, 0xed,
	0x6c, 0x56, 0xb1, 0xe7, 0x70, 0xf6, 0x42, 0xe6, 0xb9, 0x2c, 0x5e, 0x2a, 0x15, 0xa3, 0x2e, 0x65,
	0xa1, 0x91, 0x9c, 0x42, 0x1b, 0x95, 0xa2, 0xad, 0x69, 0xeb, 0xaa, 0x13, 0xdb, 0x23, 0xa1, 0x70,
	0x9c, 0xa3, 0xd6, 0x3c, 0x43, 0x7a, 0x38, 0x6d, 0x5d, 0xf5, 0xe2, 0xda, 0x64, 0x0c, 0xa2, 0x9f,
	0x71, 0xb9, 0x94, 0x31, 0xfe, 0x5d, 0xa1, 0x36, 0x84, 0xc0, 0x51, 0xc1, 0x73, 0x74, 0xce, 0xbd,
	0xd8, 0x9d, 0xd9, 0xd7, 0x30, 0x08, 0x9c, 0x90, 0x60, 0x02, 0x27, 0x99, 0x42, 0x34, 0xa2, 0xc8,
	0x42, 0xbc, 0x8d, 0xcd, 0x0c, 0x8c, 0x62, 0xcc, 0x84, 0x36, 0xa8, 0xea, 0x98, 0x63, 0xe8, 0x94,
	0x0b, 0x59, 0xd4, 0x41, 0xbd, 0x61, 0xab, 0x2c, 0x57, 0x69, 0xf0, 0xb7, 0x47, 0xf2, 0x25, 0xf4,
	0x45, 0xf1, 0x46, 0x18, 0x9c, 0x27, 0x32, 0x45, 0x7a, 0xe4, 0x6e, 0xc0, 0x43, 0x2f, 0x64, 0x8a,
	0xb6, 0x8d, 0x44, 0x56, 0x85, 0x51, 0x6b, 0xda, 0x71, 0xcd, 0xd5, 0x26, 0xfb, 0x0a, 0x4e, 0xb7,
	0x59, 0x1f, 0x1b, 0x03, 0xbb, 0x81, 0xe8, 0xb5, 0xcc, 0x44, 0xf1, 0x81, 0x85, 0xb1, 0xa7, 0x30,
	0x08, 0x7e, 0x8f, 0x86, 0xbe, 0x80, 0xc1, 0x8f, 0x52, 0x65, 0x68, 0xde, 0x1b, 0x9b, 0x2d, 0x60,
	0x58, 0xd3, 0x3e, 0x7c, 0x59, 0xdb, 0x98, 0xed, 0x66, 0xbd, 0x63, 0xe8, 0x60, 0xce, 0xc5, 0x32,
	0x0c, 0xcc, 0x1b, 0x6c, 0x05, 0xa3, 0xdf, 0x31, 0xa9, 0x94, 0x30, 0xeb, 0xf7, 0xb7, 0x7b, 0x09,
	0x23, 0x77, 0x98, 0xf3, 0xca, 0x2c, 0xfc, 0xe4, 0x7d, 0xda, 0x81, 0x83, 0x7f, 0xa8, 0xcc, 0xc2,
	0x0d, 0xff, 0x12, 0x46, 0x2e, 0x72, 0x83, 0xe7, 0xcb, 0x18, 0x38, 0xb8, 0xe6, 0x31, 0x0e, 0xa7,
	0xdb, 0xc4, 0x1f, 0xd1, 0xe4, 0x53, 0x88, 0x74, 0xf0, 0x9f, 0xdf, 0xe3, 0xda, 0x25, 0x89, 0xe2,
	0x7e, 0x8d, 0xfd, 0x82, 0x6b, 0x76, 0x0e, 0x67, 0x77, 0xd2, 0x88, 0x04, 0x5f, 0x0b, 0x5d, 0x0f,
	0x9c, 0xfd, 0xdf, 0x02, 0xd2, 0x44, 0x43, 0xea, 0x1b, 0xe8, 0x16, 0x0e, 0xa5, 0xad, 0x69, 0xfb,
	0xaa, 0x3f, 0xfb, 0xe2, 0x3a, 0x9b, 0x55, 0xd7, 0x0f, 0x89, 0x01, 0x8a, 0x03, 0xbb, 0x2e, 0xf9,
	0xf0, 0x9d, 0x25, 0xb7, 0x77, 0x4a, 0x9e, 0xfc, 0x03, 0x5d, 0xef, 0x4d, 0x86, 0x70, 0x28, 0xd2,
	0xd0, 0xe7, 0xa1, 0x48, 0xed, 0xc8, 0x8d, 0x30, 0xcb, 0xba, 0x49, 0x6f, 0x90, 0x29, 0xf4, 0x6f,
	0x51, 0x27, 0x4a, 0x94, 0x46, 0xc8, 0x22, 0x44, 0x6b, 0x42, 0xe4, 0x12, 0x86, 0x89, 0x42, 0x6e,
	0xf0, 0x96, 0x1b, 0xfc, 0x43, 0xe4, 0xf5, 0x6b, 0xd8, 0x43, 0xd9, 0x05, 0x9c, 0xfb, 0xcc, 0xb7,
	0x68, 0xb8, 0x58, 0xd6, 0x9b, 0xde, 0x2b, 0x83, 0xfd, 0x77, 0x04, 0xe3, 0x5d, 0xde, 0x47, 0x2c,
	0xc6, 0x07, 0x6d, 0x3f, 0xec, 0xed, 0x68, 0xaf, 0xb7, 0xb4, 0xd1, 0x5b, 0xc7, 0xf7, 0xd6, 0x80,
	0xfc, 0x2b, 0x2e, 0x0c, 0x16, 0x86, 0x76, 0x7d, 0x86, 0x60, 0x92, 0x27, 0xd0, 0x4d, 0xe4, 0x1b,
	0x54, 0x9a, 0x1e, 0xbb, 0xa5, 0x07, 0x8b, 0x5c, 0xc0, 0x30, 0x50, 0xe6, 0x22, 0xe7, 0x19, 0x6a,
	0x7a, 0xe2, 0xee, 0x07, 0x01, 0x7d, 0xe5, 0x40, 0xab, 0x5d, 0x66, 0x5d, 0x22, 0xed, 0xb9, 0x12,
	0xdd, 0x99, 0x7c, 0x0e, 0x3d, 0xfb, 0x77, 0xee, 0x44, 0x0d, 0xbc, 0x56, 0x59, 0xe0, 0x8e, 0xe7,
	0x68, 0xf3, 0xd9, 0x1f, 0xb3, 0x54, 0xb4, 0xef, 0x6e, 0x82, 0x65, 0xf1, 0x15, 0x8a, 0x6c, 0x61,
	0x68, 0xe4, 0x42, 0x05, 0xcb, 0xe2, 0x7a, 0xc1, 0x15, 0x6a, 0x3a, 0xf0, 0xb8, 0xb7, 0x6c, 0xe2,
	0x85, 0x30, 0x9a, 0x0e, 0x7d, 0x62, 0x7b, 0xb6, 0x1a, 0x99, 0xc8, 0x3c, 0xc7, 0xc2, 0x68, 0x3a,
	0x72, 0xf8, 0xc6, 0x76, 0xfb, 0x17, 0xba, 0x5c, 0xf2, 0xf5, 0xaf, 0x5c, 0xdd, 0xd3, 0xd3, 0x69,
	0xeb, 0xea, 0x24, 0x6e, 0x42, 0x56, 0x0a, 0xfd, 0xa6, 0xe7, 0xc6, 0x2e, 0xff, 0xcc, 0x4b, 0xa1,
	0x87, 0xec, 0xe2, 0x2d, 0xa1, 0x2a, 0xd3, 0x0d, 0x81, 0x78, 0x82, 0x87, 0x1c, 0x81, 0xc2, 0x31,
	0x4f, 0x73, 0x51, 0xbc, 0x4a, 0xe9, 0xb9, 0xd7, 0xca, 0x60, 0xda, 0x69, 0xba, 0xe3, 0xbc, 0x10,
	0xc9, 0xbd, 0x9b, 0xcb, 0xd8, 0xbf, 0x63, 0x87, 0xde, 0x05, 0x70, 0xf6, 0x6f, 0x1b, 0xa2, 0x9f,
	0xb8, 0xc1, 0x15, 0x5f, 0xcf, 0xfe, 0xd4, 0xa8, 0xc8, 0x33, 0xe8, 0xb8, 0xcf, 0x00, 0x39, 0x73,
	0x4f, 0xa8, 0xf9, 0xd9, 0x98, 0x90, 0x26, 0xe4, 0x7f, 0x5b, 0xec, 0x80, 0x7c, 0x07, 0x27, 0xb5,
	0x2a, 0x93, 0xb1, 0x63, 0xec, 0x7d, 0x1a, 0x26, 0x9f, 0xec, 0xa1, 0x1b, 0xd7, 0x67, 0xd0, 0x71,
	0x92, 0x1b, 0x92, 0x35, 0x65, 0x3b, 0x24, 0xdb, 0x51, 0x64, 0x76, 0x40, 0x6e, 0xa0, 0xe7, 0xa5,
	0xf5, 0xb7, 0x55, 0x4a, 0x3c, 0x65, 0x47, 0x91, 0x27, 0xe7, 0x3b, 0xd8, 0xc6, 0xef, 0x7b, 0x88,
	0xac, 0x76, 0xd5, 0x9a, 0x15, 0x0a, 0xdd, 0xd3, 0xce, 0x50, 0xe8, 0xbe, 0xb0, 0xb1, 0x03, 0xf2,
	0x1c, 0x60, 0x2b, 0x26, 0xe4, 0xc9, 0x03, 0x75, 0xf1, 0xee, 0x9f, 0x3e, 0xa2, 0x3a, 0xec, 0x80,
	0xbc, 0x84, 0xa8, 0xf9, 0x34, 0x09, 0x6d, 0x50, 0x77, 0x5e, 0xf5, 0xe4, 0xb3, 0x77, 0xdc, 0xd4,
	0x61, 0xfe, 0xea, 0xba, 0xff, 0x0a, 0xbe, 0x7d, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x39, 0xec, 0xc4,
	0x8e, 0x27, 0x08, 0x00, 0x00,
}
