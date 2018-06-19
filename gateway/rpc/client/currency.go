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

type CurrencyRPCCli struct {
	conn proto.CurrencyRPCService
}

func (s *CurrencyRPCCli) CallAdmin(name string) (rsp *proto.AdminResponse, err error) {
	rsp, err = s.conn.AdminCmd(context.TODO(), &proto.AdminRequest{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func NewCurrencyRPCCli() (u *CurrencyRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("currency.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_currency")
	greeter := proto.NewCurrencyRPCService(service_name, service.Client())
	u = &CurrencyRPCCli{
		conn: greeter,
	}
	return
}

// 调用 rpc 新增广告(买卖)
func (s *CurrencyRPCCli) CallAddAds(req *proto.AdsRequest) (int, error) {
	rsp, err := s.conn.AddAds(context.TODO(), req)
	return int(rsp.Code), err
}

// 调用 rpc 修改广告(买卖)
func (s *CurrencyRPCCli) CallUpdatedAds(req *proto.AdsRequest) (int, error) {
	rsp, err := s.conn.UpdatedAds(context.TODO(), req)
	return int(rsp.Code), err
}

// 调用 rpc 修改广告(买卖)状态
func (s *CurrencyRPCCli) CallUpdatedAdsStatus(req *proto.AdsStatusRequest) (int, error) {
	rsp, err := s.conn.UpdatedAdsStatus(context.TODO(), req)
	return int(rsp.Code), err
}

// 调用 rpc 法币交易列表 - (广告(买卖))
func (s *CurrencyRPCCli) CallAdsList(req *proto.AdsListRequest) (*proto.AdsListResponse, error) {
	rsp, err := s.conn.AdsList(context.TODO(), req)
	return rsp, err
}
