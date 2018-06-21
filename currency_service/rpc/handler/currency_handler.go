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

// 获取广告(买卖)
func (s *RPCServer) GetAds(ctx context.Context, req *proto.AdsGetRequest, rsp *proto.AdsModel) error {

	data := new(model.Ads).Get(req.Id)
	if data == nil {
		return nil
	}

	rsp.Id = data.Id
	rsp.Uid = data.Uid
	rsp.TypeId = data.TypeId
	rsp.TokenId = data.TokenId
	rsp.TokenName = data.TokenName
	rsp.Price = data.Price
	rsp.Num = data.Num
	rsp.Premium = data.Premium
	rsp.AcceptPrice = data.AcceptPrice
	rsp.MinLimit = data.MinLimit
	rsp.MaxLimit = data.MaxLimit
	rsp.IsTwolevel = data.IsTwolevel
	rsp.Pays = data.Pays
	rsp.Remarks = data.Remarks
	rsp.Reply = data.Reply
	rsp.IsUsd = data.IsUsd
	rsp.States = data.States
	rsp.CreatedTime = data.CreatedTime
	rsp.UpdatedTime = data.UpdatedTime

	return nil
}

// 新增广告(买卖)
func (s *RPCServer) AddAds(ctx context.Context, req *proto.AdsModel, rsp *proto.CurrencyResponse) error {

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
func (s *RPCServer) UpdatedAds(ctx context.Context, req *proto.AdsModel, rsp *proto.CurrencyResponse) error {

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

// 法币交易列表 (广告(买卖))
func (s *RPCServer) AdsList(ctx context.Context, req *proto.AdsListRequest, rsp *proto.AdsListResponse) error {
	data, total := new(model.Ads).AdsList(req.TypeId, req.TokenId, req.Page, req.PageNum)
	if data == nil || total <= 0 {
		return nil
	}

	listLen := len(data)
	listData := make([]*proto.AdsLists, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.AdsLists{
			Id:          data[i].Id,
			Uid:         data[i].Uid,
			Price:       data[i].Price,
			Num:         data[i].Num,
			MinLimit:    data[i].MinLimit,
			MaxLimit:    data[i].MaxLimit,
			Pays:        data[i].Pays,
			CreatedTime: data[i].CreatedTime,
			UpdatedTime: data[i].UpdatedTime,
			UserVolume:  data[i].Success,
			TypeId:      data[i].TypeId,
			TokenId:     data[i].TokenId,
			TokenName:   data[i].TokenName,
		}
		listData[i] = adsLists
	}

	rsp.Page = req.Page
	rsp.PageNum = req.PageNum
	rsp.Total = uint64(total)
	rsp.Data = listData

	return nil
}

// 个人法币交易列表 (广告(买卖))
func (s *RPCServer) AdsUserList(ctx context.Context, req *proto.AdsListRequest, rsp *proto.AdsListResponse) error {
	data, total := new(model.Ads).AdsUserList(req.Uid, req.TypeId, req.Page, req.PageNum)
	if data == nil || total <= 0 {
		return nil
	}

	listLen := len(data)
	listData := make([]*proto.AdsLists, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.AdsLists{
			Id:          data[i].Id,
			Uid:         data[i].Uid,
			Price:       data[i].Price,
			Num:         data[i].Num,
			MinLimit:    data[i].MinLimit,
			MaxLimit:    data[i].MaxLimit,
			Pays:        data[i].Pays,
			CreatedTime: data[i].CreatedTime,
			UpdatedTime: data[i].UpdatedTime,
			TypeId:      data[i].TypeId,
			TokenId:     data[i].TokenId,
			TokenName:   data[i].TokenName,
		}
		listData[i] = adsLists
	}

	rsp.Page = req.Page
	rsp.PageNum = req.PageNum
	rsp.Total = uint64(total)
	rsp.Data = listData

	return nil
}

// 获取货币类型
func (s *RPCServer) GetCurrencyTokens(ctx context.Context, req *proto.CurrencyTokensRequest, rsp *proto.CurrencyTokens) error {
	data := new(model.Tokens).Get(req.Id, req.Name)
	if data == nil {
		return nil
	}

	rsp.Id = data.Id
	rsp.Name = data.Name
	rsp.CnName = data.CnName

	return nil
}

// 获取货币类型列表
func (s *RPCServer) CurrencyTokensList(ctx context.Context, req *proto.CurrencyTokensRequest, rsp *proto.CurrencyTokensListResponse) error {
	data := new(model.Tokens).List()
	if data == nil {
		return nil
	}

	listLen := len(data)
	listData := make([]*proto.CurrencyTokens, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.CurrencyTokens{
			Id:     data[i].Id,
			Name:   data[i].Name,
			CnName: data[i].CnName,
		}
		listData[i] = adsLists
	}

	rsp.Data = listData
	return nil
}

// 获取支付方式
func (s *RPCServer) GetCurrencyPays(ctx context.Context, req *proto.CurrencyPaysRequest, rsp *proto.CurrencyPays) error {
	data := new(model.Pays).Get(req.Id, req.EnPay)
	if data == nil {
		return nil
	}

	rsp.Id = data.Id
	rsp.TypeId = data.TypeId
	rsp.ZhPay = data.ZhPay
	rsp.EnPay = data.EnPay
	rsp.States = data.States

	return nil
}

// 获取支付方式列表
func (s *RPCServer) CurrencyPaysList(ctx context.Context, req *proto.CurrencyPaysRequest, rsp *proto.CurrencyPaysListResponse) error {
	data := new(model.Pays).List()
	if data == nil {
		return nil
	}

	listLen := len(data)
	listData := make([]*proto.CurrencyPays, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.CurrencyPays{
			Id:     data[i].Id,
			ZhPay:  data[i].ZhPay,
			EnPay:  data[i].EnPay,
			States: data[i].States,
		}
		listData[i] = adsLists
	}

	rsp.Data = listData
	return nil
}
