// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: rpc/user.proto

package g2u

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Gateway2User service

type Gateway2UserService interface {
	Hello(ctx context.Context, in *HelloRequest, opts ...client.CallOption) (*HelloResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	ForgetPwd(ctx context.Context, in *ForgetRequest, opts ...client.CallOption) (*ForgetResponse, error)
	AuthSecurity(ctx context.Context, in *SecurityRequest, opts ...client.CallOption) (*SecurityResponse, error)
	NoticeList(ctx context.Context, in *NoticeListRequest, opts ...client.CallOption) (*NoticeListResponse, error)
	NoticeDetail(ctx context.Context, in *NoticeDetailRequest, opts ...client.CallOption) (*NoticeDetailResponse, error)
}

type gateway2UserService struct {
	c    client.Client
	name string
}

func NewGateway2UserService(name string, c client.Client) Gateway2UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "g2u"
	}
	return &gateway2UserService{
		c:    c,
		name: name,
	}
}

func (c *gateway2UserService) Hello(ctx context.Context, in *HelloRequest, opts ...client.CallOption) (*HelloResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.Hello", in)
	out := new(HelloResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateway2UserService) Register(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*RegisterResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.Register", in)
	out := new(RegisterResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateway2UserService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateway2UserService) ForgetPwd(ctx context.Context, in *ForgetRequest, opts ...client.CallOption) (*ForgetResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.ForgetPwd", in)
	out := new(ForgetResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateway2UserService) AuthSecurity(ctx context.Context, in *SecurityRequest, opts ...client.CallOption) (*SecurityResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.AuthSecurity", in)
	out := new(SecurityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateway2UserService) NoticeList(ctx context.Context, in *NoticeListRequest, opts ...client.CallOption) (*NoticeListResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.NoticeList", in)
	out := new(NoticeListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateway2UserService) NoticeDetail(ctx context.Context, in *NoticeDetailRequest, opts ...client.CallOption) (*NoticeDetailResponse, error) {
	req := c.c.NewRequest(c.name, "Gateway2User.NoticeDetail", in)
	out := new(NoticeDetailResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Gateway2User service

type Gateway2UserHandler interface {
	Hello(context.Context, *HelloRequest, *HelloResponse) error
	Register(context.Context, *RegisterRequest, *RegisterResponse) error
	Login(context.Context, *LoginRequest, *LoginResponse) error
	ForgetPwd(context.Context, *ForgetRequest, *ForgetResponse) error
	AuthSecurity(context.Context, *SecurityRequest, *SecurityResponse) error
	NoticeList(context.Context, *NoticeListRequest, *NoticeListResponse) error
	NoticeDetail(context.Context, *NoticeDetailRequest, *NoticeDetailResponse) error
}

func RegisterGateway2UserHandler(s server.Server, hdlr Gateway2UserHandler, opts ...server.HandlerOption) {
	type gateway2User interface {
		Hello(ctx context.Context, in *HelloRequest, out *HelloResponse) error
		Register(ctx context.Context, in *RegisterRequest, out *RegisterResponse) error
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		ForgetPwd(ctx context.Context, in *ForgetRequest, out *ForgetResponse) error
		AuthSecurity(ctx context.Context, in *SecurityRequest, out *SecurityResponse) error
		NoticeList(ctx context.Context, in *NoticeListRequest, out *NoticeListResponse) error
		NoticeDetail(ctx context.Context, in *NoticeDetailRequest, out *NoticeDetailResponse) error
	}
	type Gateway2User struct {
		gateway2User
	}
	h := &gateway2UserHandler{hdlr}
	s.Handle(s.NewHandler(&Gateway2User{h}, opts...))
}

type gateway2UserHandler struct {
	Gateway2UserHandler
}

func (h *gateway2UserHandler) Hello(ctx context.Context, in *HelloRequest, out *HelloResponse) error {
	return h.Gateway2UserHandler.Hello(ctx, in, out)
}

func (h *gateway2UserHandler) Register(ctx context.Context, in *RegisterRequest, out *RegisterResponse) error {
	return h.Gateway2UserHandler.Register(ctx, in, out)
}

func (h *gateway2UserHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.Gateway2UserHandler.Login(ctx, in, out)
}

func (h *gateway2UserHandler) ForgetPwd(ctx context.Context, in *ForgetRequest, out *ForgetResponse) error {
	return h.Gateway2UserHandler.ForgetPwd(ctx, in, out)
}

func (h *gateway2UserHandler) AuthSecurity(ctx context.Context, in *SecurityRequest, out *SecurityResponse) error {
	return h.Gateway2UserHandler.AuthSecurity(ctx, in, out)
}

func (h *gateway2UserHandler) NoticeList(ctx context.Context, in *NoticeListRequest, out *NoticeListResponse) error {
	return h.Gateway2UserHandler.NoticeList(ctx, in, out)
}

func (h *gateway2UserHandler) NoticeDetail(ctx context.Context, in *NoticeDetailRequest, out *NoticeDetailResponse) error {
	return h.Gateway2UserHandler.NoticeDetail(ctx, in, out)
}
