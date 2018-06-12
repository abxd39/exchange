package rpc

import (
	proto "digicon/proto/rpc"
	cf "digicon/token_service/conf"
	"digicon/token_service/rpc/handler"
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

	proto.RegisterTokenRPCHandler(service.Server(), new(handler.RPCServer))

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
