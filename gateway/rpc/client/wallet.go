package client

import (
	"context"
	cf "digicon/gateway/conf"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type WalletRPCCli struct {
	conn proto.Gateway2WallerService
}

func (s *WalletRPCCli) CallGreet(name string) (rsp *proto.HelloResponse2, err error) {
	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest2{Name: name})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func NewWalletRPCCli() (u *WalletRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("wallet.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_wallet")
	greeter := proto.NewGateway2WallerService(service_name, service.Client())
	u = &WalletRPCCli{
		conn: greeter,
	}
	return
}
