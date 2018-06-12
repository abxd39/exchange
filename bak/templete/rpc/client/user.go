package client

import (
	"context"
	proto "digicon/proto/rpc"
	cf "digicon/wallet_service/conf"
	"github.com/golang/glog"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type UserRPCCli struct {
	conn proto.Gateway2WallerService
}

/*
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	fmt.Printf("[Log Wrapper] ctx: %v service: %s method: %s\n", md, req.Service(), req.Method())
	return l.Client.Call(ctx, req, rsp)
}

func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}
*/

func (s *UserRPCCli) CallGreet(name string) (rsp *proto.HelloResponse2, err error) {
	rsp, err = s.conn.Hello(context.TODO(), &proto.HelloRequest2{})
	if err != nil {
		glog.Errorln(err.Error())
		return
	}
	return
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
	greeter := proto.NewGateway2WallerService(service_name, service.Client())
	u = &UserRPCCli{
		conn: greeter,
	}
	return
	return
}
