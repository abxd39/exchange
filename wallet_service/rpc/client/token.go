package client

import (
	"context"
	proto "digicon/proto/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	. "digicon/wallet_service/utils"
)

type TokenRPCCli struct {
	conn proto.TokenRPCService
}

func (s *TokenRPCCli) CallSubTokenWithFronze(p *proto.SubTokenWithFronzeRequest) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.SubTokenWithFronzen(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *TokenRPCCli) CallCancelSubTokenWithFronze(p *proto.CancelFronzeTokenRequest) (rsp *proto.CommonErrResponse, err error) {
	rsp, err = s.conn.CancelFronzeToken(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func NewTokenRPCCli() (u *TokenRPCCli) {
	consul_addr := Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("token.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := Cfg.MustValue("base", "service_client_token")
	greeter := proto.NewTokenRPCService(service_name, service.Client())
	u = &TokenRPCCli{
		conn: greeter,
	}
	return
}
