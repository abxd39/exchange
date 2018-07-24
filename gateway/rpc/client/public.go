package client

import (
	"context"
	cf "digicon/gateway/conf"
	log "github.com/sirupsen/logrus"
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
		log.Errorln(err.Error())
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

func (s *PublciRPCCli) CallArticle(id int32) (*proto.ArticleResponse, error) {
	return s.conn.Article(context.TODO(), &proto.ArticleRequest{
		Id: id,
	})

}

func (s *PublciRPCCli) CallArticleList(ty, page, page_num int32) (*proto.ArticleListResponse, error) {
	return s.conn.ArticleList(context.TODO(), &proto.ArticleListRequest{
		ArticleType: ty,
		Page:        page,
		PageNum:     page_num,
	})
}

func (s *PublciRPCCli) CallAddFriendlyLink(req *proto.AddFriendlyLinkRequest) (rsp *proto.AddFriendlyLinkResponse, err error) {
	return s.conn.AddFriendlyLink(context.TODO(), req)
}

func (s *PublciRPCCli) CallGetFriendlyLink(req *proto.FriendlyLinkRequest) (rsp *proto.FriendlyLinkResponse, err error) {
	return s.conn.GetFriendlyLink(context.TODO(), req)
}

func (s *PublciRPCCli) CallGetBannerList(req *proto.BannerRequest) (rsp *proto.BannerResponse, err error) {
	return s.conn.GetBannerList(context.TODO(), req)
}

func (s *PublciRPCCli) CallGetArticleTypeList(req *proto.ArticleTypeRequest) (rsp *proto.ArticleTypeListResponse, err error) {
	return s.conn.ArticleTypeList(context.TODO(), req)
}
