package  client

import (
	"context"
	proto "digicon/proto/rpc"
)


//type PriceRPCCli struct {
//	priceconn proto.PriceRPCService
//}
//
//
//func NewPriceRPCCli() ( p *PriceRPCCli) {
//	consul_addr := cf.Cfg.MustValue("consul", "addr")
//	r := consul.NewRegistry(registry.Addrs(consul_addr))
//	service := micro.NewService(
//		micro.Name("greeter.client"),
//		micro.Registry(r),
//	)
//	service.Init()
//	price_client_name := cf.Cfg.MustValue("base", "service_price")
//	fmt.Println("service_name,", price_client_name)
//	priceGreeter := proto.NewPriceRPCService(price_client_name, service.Client())
//	p = &PriceRPCCli{
//		priceconn: priceGreeter,
//	}
//	return
//}



func (p *UserRPCCli) CallGetSymbolsRate(symbols []string) (rsp *proto.GetSymbolsRateResponse, err error){
	return  p.priceconn.GetSymbolsRate(context.TODO(), &proto.GetSymbolsRateRequest{Symbols:symbols})
}