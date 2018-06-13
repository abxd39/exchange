package client

import (
	"context"
	cf "digicon/gateway/conf"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type PublciRPCCli struct {
	conn proto.PublicRPCService
}

func (s *PublciRPCCli) CallAdmin(name string) (rsp *proto.AdminResponse, err error) {
	rsp, err = s.conn.AdminCmd(context.TODO(), &proto.AdminRequest{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
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

func (s *PublciRPCCli) CallArticlesDesc(id int32) (rsp *proto.ArticlesDetailResponse, err error) {
	rsp, err = s.conn.ArticlesDetail(context.TODO(), &proto.ArticlesDetailRequest{
		Id: id,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *PublciRPCCli) CallArticlesList() (rsp *proto.ArticlesListResponse, err error) {
	rsp, err = s.conn.ArticlesList(context.TODO(), &proto.ArticlesListRequest{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
