package rpc

import (
	proto "digicon/proto/rpc"
	cf "digicon/wallet_service/conf"
	"digicon/wallet_service/rpc/handler"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"time"
)

func RPCServerInit() {
	service_name := cf.Cfg.MustValue("base", "service_name")
	r := consul.NewRegistry(registry.Addrs("47.106.136.96:8500"))
	service := micro.NewService(
		micro.Name(service_name),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(r),
	)
	service.Init()

	proto.RegisterGateway2WallerHandler(service.Server(), new(handler.Greeter))

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
