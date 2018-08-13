package client

import (
	"context"
	cf "digicon/gateway/conf"
	proto "digicon/proto/rpc"
	log "github.com/sirupsen/logrus"

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
		log.Errorln(err.Error())
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

// 调用 rpc 获取广告(买卖)
func (s *CurrencyRPCCli) CallGetAds(req *proto.AdsGetRequest) (*proto.AdsModel, error) {
	rsp, err := s.conn.GetAds(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 新增广告(买卖)
func (s *CurrencyRPCCli) CallAddAds(req *proto.AdsModel) (int, error) {
	rsp, err := s.conn.AddAds(context.TODO(), req)
	return int(rsp.Code), err
}

// 调用 rpc 修改广告(买卖)
func (s *CurrencyRPCCli) CallUpdatedAds(req *proto.AdsModel) (int, error) {
	rsp, err := s.conn.UpdatedAds(context.TODO(), req)
	return int(rsp.Code), err
}

// 调用 rpc 修改广告(买卖)状态
func (s *CurrencyRPCCli) CallUpdatedAdsStatus(req *proto.AdsStatusRequest) (*proto.CurrencyResponse, error) {
	rsp, err := s.conn.UpdatedAdsStatus(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 法币交易列表 - (广告(买卖))
func (s *CurrencyRPCCli) CallAdsList(req *proto.AdsListRequest) (*proto.AdsListResponse, error) {
	rsp, err := s.conn.AdsList(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 个人法币交易列表 - (广告(买卖))
func (s *CurrencyRPCCli) CallAdsUserList(req *proto.AdsListRequest) (*proto.AdsListResponse, error) {
	rsp, err := s.conn.AdsUserList(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 获取货币类型
func (s *CurrencyRPCCli) CallGetCurrencyTokens(req *proto.CurrencyTokensRequest) (*proto.CurrencyTokens, error) {
	rsp, err := s.conn.GetCurrencyTokens(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 获取货币类型列表
func (s *CurrencyRPCCli) CallCurrencyTokensList(req *proto.CurrencyTokensRequest) (*proto.CurrencyTokensListResponse, error) {
	rsp, err := s.conn.CurrencyTokensList(context.TODO(), req)
	return rsp, err
}

//
func (s *CurrencyRPCCli) CallDisplayCurrencyTokens ( req *proto.CurrencyTokensRequest) ( *proto.CurrencyTokensListResponse, error){
	rsp, err := s.conn.DisplayCurrencyTokens(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 获取支付方式
func (s *CurrencyRPCCli) CallGetCurrencyPays(req *proto.CurrencyPaysRequest) (*proto.CurrencyPays, error) {
	rsp, err := s.conn.GetCurrencyPays(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 获取支付方式列表
func (s *CurrencyRPCCli) CallCurrencyPaysList(req *proto.CurrencyPaysRequest) (*proto.CurrencyPaysListResponse, error) {
	rsp, err := s.conn.CurrencyPaysList(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 新增订单聊天
func (s *CurrencyRPCCli) CallGetCurrencyChats(req *proto.CurrencyChats) (int, error) {
	rsp, err := s.conn.GetCurrencyChats(context.TODO(), req)
	return int(rsp.Code), err
}

// 调用 rpc 获取订单聊天列表
func (s *CurrencyRPCCli) CallCurrencyChatsList(req *proto.CurrencyChats) (*proto.CurrencyChatsListResponse, error) {
	rsp, err := s.conn.CurrencyChatsList(context.TODO(), req)
	return rsp, err
}

// 调用 rpc 获取用户虚拟货币资产
func (s *CurrencyRPCCli) CallGetUserCurrency(req *proto.UserCurrencyRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.GetUserCurrency(context.TODO(), req)
	return rsp, err
}
func (s *CurrencyRPCCli) CallGetUserCurrencyDetail(req *proto.UserCurrencyRequest) (*proto.UserCurrency, error) {
	rsp, err := s.conn.GetUserCurrencyDetail(context.TODO(), req)
	return rsp, err
}

// get 售价
func (s *CurrencyRPCCli) CallGetSellingPrice(req *proto.SellingPriceRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.GetSellingPrice(context.TODO(), req)
	return rsp, err
}

// get
func (s *CurrencyRPCCli) CallGetCurrencyBalance(req *proto.GetCurrencyBalanceRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.GetCurrencyBalance(context.TODO(), req)
	return rsp, err
}

/*
	获取用户好评率
*/
func (s *CurrencyRPCCli) CallGetUserRating(req *proto.GetUserRatingRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.GetUserRating(context.TODO(), req)
	return rsp, err
}

/*
	rpc添加用户余额
*/
func (s *CurrencyRPCCli) CallAddUserBalance(req *proto.AddUserBalanceRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.AddUserBalance(context.TODO(), req)
	return rsp, err
}

/*
	获取资金明细
*/
func (s *CurrencyRPCCli) CallGetAssetDetail(req *proto.GetAssetDetailRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.GetAssetDetail(context.TODO(), req)
	return rsp, err
}

/*
	法币划转到代币
*/
func (s *CurrencyRPCCli) CallTransferToToken(req *proto.TransferToTokenRequest) (*proto.OtherResponse, error) {
	rsp, err := s.conn.TransferToToken(context.TODO(), req)
	return rsp, err
}
