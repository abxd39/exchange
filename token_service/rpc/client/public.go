package client

import (
	proto "digicon/proto/rpc"
	cf "digicon/token_service/conf"

	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type PublciRPCCli struct {
	conn proto.PublicRPCService
}

func (s *PublciRPCCli) CallGetTokensList(ids []int32) (*proto.TokensResponse, error) {

	return s.conn.GetTokensList(context.Background(), &proto.TokensRequest{
		Tokens: ids,
	})

}

func NewPublciRPCCli() (u *PublciRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("token.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_public")
	greeter := proto.NewPublicRPCService(service_name, service.Client())
	u = &PublciRPCCli{
		conn: greeter,
	}
	return
}
