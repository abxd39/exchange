package rpc

import (
	cf "digicon/token_service/conf"
	"digicon/token_service/rpc/handler"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"time"
	proto "digicon/proto/rpc"
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

	svr:=service.Server()

	err:=micro.RegisterSubscriber("topic.go.micro.srv.price",svr,new(handler.Subscriber))
	if err!=nil {
		log.Fatalln(err)
	}
	proto.RegisterTokenRPCHandler(svr, new(handler.RPCServer))
	if err!=nil {
		log.Fatalln(err)
	}
	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
