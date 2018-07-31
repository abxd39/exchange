package client

import (
	proto "digicon/proto/rpc"
	. "digicon/ws_service/conf"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro"
	"context"
	log "github.com/sirupsen/logrus"
)

type UserRPCCli struct {
	conn proto.UserRPCService
}


func NewUserRPCCli() (u *UserRPCCli) {
	consul_addr := Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := Cfg.MustValue("base", "service_client_user")
	greeter := proto.NewUserRPCService(service_name, service.Client())
	u = &UserRPCCli{
		conn: greeter,
	}
	return
}


func (s *UserRPCCli) CallTokenVerify(uid uint64, token []byte) (rsp *proto.TokenVerifyResponse, err error) {
	rsp, err = s.conn.TokenVerify(context.TODO(), &proto.TokenVerifyRequest{
		Uid:   uid,
		Token: token,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

