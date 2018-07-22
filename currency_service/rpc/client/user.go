package client

import (
	"context"
	cf "digicon/currency_service/conf"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type UserRPCCli struct {
	conn     proto.Gateway2WallerService
	userconn proto.UserRPCService
}

//func (s *UserRPCCli) CallGreet(name string) (rsp *proto.HelloResponse2, err error) {
//	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest2{})
//	if err != nil {
//		Log.Errorln(err.Error())
//		return
//	}
//	return
//}

func (s *UserRPCCli) CallGetNickName(uids []uint64) (rsp *proto.UserGetNickNameResponse, err error) {
	fmt.Println("uids:", uids)
	return s.userconn.GetNickName(context.TODO(), &proto.UserGetNickNameRequest{Uid: uids})
}

func (s *UserRPCCli) CallGetAuthInfo(uid uint64) (rsp *proto.GetAuthInfoResponse, err error) {
	return s.userconn.GetAuthInfo(context.TODO(), &proto.GetAuthInfoRequest{Uid: uid})
}

func NewUserRPCCli() (u *UserRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_name")
	user_client_name := cf.Cfg.MustValue("base", "service_client_user")

	fmt.Println("service_name,", service_name, " user_client_name: ", user_client_name)
	greeter := proto.NewGateway2WallerService(service_name, service.Client())
	userGreeter := proto.NewUserRPCService(user_client_name, service.Client())
	u = &UserRPCCli{
		conn:     greeter,
		userconn: userGreeter,
	}
	return
}
