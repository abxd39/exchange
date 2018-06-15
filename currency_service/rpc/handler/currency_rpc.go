package handler

import (
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"log"
	"time"
)

type RPCServer struct{}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	log.Print("Received Say.Hello request")

	return nil
}

// 新增广告(买卖)
func (s *RPCServer) AddAds(ctx context.Context, req *proto.AdsRequest, rsp *proto.CurrencyResponse) error {

	// 数据过虑暂不做

	ads := new(model.Ads)
	ads.Uid = req.Uid
	ads.TypeId = req.TypeId
	ads.TokenId = req.TokenId
	ads.TokenName = req.TokenName
	ads.Price = req.Price
	ads.Num = req.Num
	ads.Premium = req.Premium
	ads.AcceptPrice = req.AcceptPrice
	ads.MinLimit = req.MinLimit
	ads.MaxLimit = req.MaxLimit
	ads.IsTwolevel = req.IsTwolevel
	ads.Pays = req.Pays
	ads.Remarks = req.Remarks
	ads.Reply = req.Reply
	ads.IsUsd = req.IsUsd
	ads.States = req.States
	ads.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	ads.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
	ads.IsDel = 0

	code := ads.Add()
	rsp.Code = int32(code)

	return nil
}

// 修改广告(买卖)
func (s *RPCServer) UpdatedAds(ctx context.Context, req *proto.AdsRequest, rsp *proto.CurrencyResponse) error {

	// 数据过虑暂不做

	ads := new(model.Ads)
	ads.Id = req.Id
	ads.Price = req.Price
	ads.Num = req.Num
	ads.Premium = req.Premium
	ads.AcceptPrice = req.AcceptPrice
	ads.MinLimit = req.MinLimit
	ads.MaxLimit = req.MaxLimit
	ads.IsTwolevel = req.IsTwolevel
	ads.Pays = req.Pays
	ads.Remarks = req.Remarks
	ads.Reply = req.Reply
	ads.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")

	code := ads.Update()
	rsp.Code = int32(code)

	return nil
}

// 修改广告(买卖)状态
func (s *RPCServer) UpdatedAdsStatus(ctx context.Context, req *proto.AdsStatusRequest, rsp *proto.CurrencyResponse) error {
	code := new(model.Ads).UpdatedAdsStatus(req.Id, req.StatusId)
	rsp.Code = int32(code)
	return nil
}
