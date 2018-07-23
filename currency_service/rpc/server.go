package rpc

import (
	cf "digicon/currency_service/conf"
	"digicon/currency_service/rpc/handler"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"time"
)



func RPCServerInit() {
	service_name := cf.Cfg.MustValue("base", "service_name")

	addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(addr))
	service := micro.NewService(
		micro.Name(service_name),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(r),
	)
	service.Init()
	proto.RegisterCurrencyRPCHandler(service.Server(), new(handler.RPCServer))

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
