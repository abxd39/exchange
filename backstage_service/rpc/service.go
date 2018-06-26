package rpc

import (
	cf "digicon/backstage_service/conf"
	"digicon/backstage_service/rpc/handler"
	proto "digicon/proto/rpc"
	"fmt"
	"log"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func RPCServerInit() {
	service_name := cf.Inif.MustValue("base", "service_name")

	consul_addr := cf.Inif.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name(service_name),
		micro.Version("1.0.0"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(r),
	)
	service.Init()

	proto.RegisterBackstageRPCHandler(service.Server(), new(handler.RPCServer))

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
