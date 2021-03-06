package client

import (
	"context"
	cf "digicon/gateway/conf"
	//log "github.com/sirupsen/logrus"
	proto "digicon/proto/rpc"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type KlineRPCCli struct {
	conn proto.KlineRPCService
}

//func (s *KlineRPCCli) CallAdmin(name string) (rsp *proto.AdminResponse, err error) {
//	rsp, err = s.conn.AdminCmd(context.TODO(), &proto.AdminRequest{})
//	if err != nil {
//		log.Errorln(err.Error())
//		return
//	}
//	return
//}

func NewKlineRPCCli() (u *KlineRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("currency.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_currency")
	greeter := proto.NewKlineRPCService(service_name, service.Client())
	u = &KlineRPCCli{
		conn: greeter,
	}
	return
}

// 调用 hello
func (s *KlineRPCCli) CallHline(req *proto.KineRequest) (int, error) {
	rsp, err := s.conn.Hline(context.TODO(), req)
	return int(rsp.Code), err
}
