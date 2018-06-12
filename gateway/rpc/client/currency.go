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

type CurrencyRPCCli struct {
	conn proto.CurrencyRPCService
}


func (s *CurrencyRPCCli) CallAdmin(name string) (rsp *proto.AdminResponse, err error) {
	rsp, err = s.conn.AdminCmd(context.TODO(), &proto.AdminRequest{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}


func NewCurrencyRPCCli() (u *CurrencyRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("currency.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_currency")
	greeter := proto.NewCurrencyRPCService(service_name, service.Client())
	u = &CurrencyRPCCli{
		conn: greeter,
	}
	return
}
