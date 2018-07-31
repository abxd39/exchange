package client

import (
	"context"
	cf "digicon/price_service/conf"
	proto "digicon/proto/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	log "github.com/sirupsen/logrus"
)

type TokenRPCCli struct {
	conn proto.TokenRPCService
}

func (s *TokenRPCCli) CallGetConfigQuene() (rsp *proto.ConfigQueneResponse, err error) {
	rsp, err = s.conn.GetConfigQuene(context.TODO(), &proto.NullRequest{})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func NewTokenRPCCli() (u *TokenRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_token")
	if service_name == "" {
		log.Fatalln("err config please check config")
	}
	greeter := proto.NewTokenRPCService(service_name, service.Client())
	u = &TokenRPCCli{
		conn: greeter,
	}
	return
}
