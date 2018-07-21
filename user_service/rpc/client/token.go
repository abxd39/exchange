package client

import (
	proto "digicon/proto/rpc"
	cf "digicon/user_service/conf"

	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type TokenRPCCli struct {
	conn proto.TokenRPCService
}

func NewTokenRPCCli() (t *TokenRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("user.token.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_token")
	greeter := proto.NewTokenRPCService(service_name, service.Client())
	t = &TokenRPCCli{
		conn: greeter,
	}
	return
}

func (s *TokenRPCCli) CallAddTokenNum(uid uint64, tokenId int32, num int64, opt proto.TOKEN_OPT_TYPE, ukey []byte, oType int32) (*proto.CommonErrResponse, error) {
	return s.conn.AddTokenNum(context.Background(), &proto.AddTokenNumRequest{
		Uid:     uid,
		TokenId: tokenId,
		Num:     num,
		Opt:     opt,
		Ukey:    ukey,
		Type:    oType,
	})
}
