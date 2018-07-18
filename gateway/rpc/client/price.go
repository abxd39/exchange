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

type PriceRPCCli struct {
	conn proto.PriceRPCService
}
func (s *PriceRPCCli) CallCurrentPrice(p *proto.CurrentPriceRequest) (rsp *proto.CurrentPriceResponse, err error) {
	rsp, err = s.conn.CurrentPrice(context.TODO(), p)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}


func NewPriceRPCCli() (u *PriceRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("token.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_price")
	greeter := proto.NewPriceRPCService(service_name, service.Client())
	u = &PriceRPCCli{
		conn: greeter,
	}
	return
}