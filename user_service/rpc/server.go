package rpc

import (
	proto "digicon/proto/rpc"
	cf "digicon/user_service/conf"
	"digicon/user_service/rpc/handler"
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
		micro.Version("1.0.0"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(r),
	)
	service.Init()

	proto.RegisterGateway2UserHandler(service.Server(), new(handler.RPCServer))

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
