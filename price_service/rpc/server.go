package rpc

import (
	cf "digicon/price_service/conf"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"time"
	"digicon/price_service/rpc/handler"
	"github.com/micro/go-micro/broker"
	"encoding/json"
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

	//pubsub := micro.NewPublisher("user.created", service.Client())
	pubsub := service.Server().Options().Broker
	if err := pubsub.Connect(); err != nil {
		log.Fatal(err)
	}
	_, err := pubsub.Subscribe("", func(p broker.Publication) error {

		return nil
	})
	if err!=nil {
		log.Fatalln(err)
	}
	//proto.RegisterPriceRPCHandler(service.Server(), new(handler.RPCServer))
	proto.RegisterPriceRPCHandler(service.Server(), handler.NewRPCServer(pubsub))

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

}
