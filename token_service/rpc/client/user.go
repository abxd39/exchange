package client

import (
	"context"
	proto "digicon/proto/rpc"
	cf "digicon/token_service/conf"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type UserRPCCli struct {
	conn proto.UserRPCService
}

func (s *UserRPCCli) CallGetUserFeeInfo(uid uint64) (rsp *proto.GetUserFeeInfoResponse, err error) {
	rsp, err = s.conn.GetUserFeeInfo(context.TODO(), &proto.InnerCommonRequest{
		Uid: uid,
	})
	return
}

func NewUserRPCCli() (u *UserRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("token.user.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_user")
	greeter := proto.NewUserRPCService(service_name, service.Client())
	u = &UserRPCCli{
		conn: greeter,
	}
	return
}
