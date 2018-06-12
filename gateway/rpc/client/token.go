package client

import (
	"context"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro"
	cf "digicon/gateway/conf"
	"github.com/micro/go-plugins/registry/consul"
)

type TokenRPCCli struct {
	conn proto.TokenRPCService
}

func (s *TokenRPCCli) CallAdmin(name string) (rsp *proto.AdminResponse, err error) {
	rsp, err = s.conn.AdminCmd(context.TODO(), &proto.AdminRequest{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}


func NewTokenRPCCli() (u *TokenRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("token.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_token")
	greeter := proto.NewTokenRPCService(service_name, service.Client())
	u = &TokenRPCCli{
		conn: greeter,
	}
	return
}
