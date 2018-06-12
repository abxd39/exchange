package client

import (
	"context"
	proto "digicon/proto/rpc"
	cf "digicon/token_service/conf"
	. "digicon/token_service/log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type UserRPCCli struct {
	conn proto.Gateway2WallerService
}

func (s *UserRPCCli) CallGreet(name string) (rsp *proto.HelloResponse2, err error) {
	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest2{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func NewUserRPCCli() (u *UserRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_name")
	greeter := proto.NewGateway2WallerService(service_name, service.Client())
	u = &UserRPCCli{
		conn: greeter,
	}
	return
}
