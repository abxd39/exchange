package client

import (
	proto "digicon/proto/rpc"
	cf "digicon/token_service/conf"

	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	log "github.com/sirupsen/logrus"
)

type PriceRPCCli struct {
	conn proto.PriceRPCService
}

func (s *PriceRPCCli) CallLastPrice(symbol string) (*proto.LastPriceResponse, error) {

	return s.conn.LastPrice(context.Background(), &proto.LastPriceRequest{
		Symbol: symbol,
	})

}

func (p *PriceRPCCli) CallGetSymbolsRate(symbols []string) (rsp *proto.GetSymbolsRateResponse, err error) {
	return p.conn.GetSymbolsRate(context.TODO(), &proto.GetSymbolsRateRequest{Symbols: symbols})
}

func (p *PriceRPCCli) CallGetCnyPrices(token_ids []int32) (rsp *proto.CnyPriceResponse, err error) {
	rsp, err = p.conn.GetCnyPrices(context.TODO(), &proto.CnyPriceRequest{
		TokenTradeId: token_ids,
	})
	if err != nil {
		log.Error(err)
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
