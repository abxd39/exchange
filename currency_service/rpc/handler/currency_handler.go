package handler

import (
	"digicon/currency_service/model"
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"log"
	"time"
	"digicon/proto/common"
	"github.com/gin-gonic/gin/json"
	"fmt"
	"digicon/common/convert"
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
			//UserVolume:  data[i].Success,
			TypeId:      data[i].TypeId,
			TokenId:     data[i].TokenId,
			TokenName:   data[i].TokenName,
			//Balance:     data[i].Balance,
			//Freeze:      data[i].Freeze,
		}
		listData[i] = adsLists
	}

	rsp.Page = req.Page
	rsp.PageNum = req.PageNum
	rsp.Total = uint64(total)
	rsp.Data = listData
	//fmt.Println("listData:", listData)
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
	//listData := []*proto.AdsLists{}
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
			Balance:     data[i].Balance,
			Freeze:      data[i].Freeze,
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

// 新增订单聊天
func (s *RPCServer) GetCurrencyChats(ctx context.Context, req *proto.CurrencyChats, rsp *proto.CurrencyResponse) error {

	chats := new(model.Chats)

	chats.OrderId = req.OrderId
	chats.IsOrderUser = req.IsOrderUser
	chats.Uid = req.Uid
	chats.Uname = req.Uname
	chats.Content = req.Content
	chats.States = 1
	chats.CreatedTime = time.Now().Format("2006-01-02 15:04:05")

	code := chats.Add()
	rsp.Code = int32(code)

	return nil
}

// 获取订单聊天列表
func (s *RPCServer) CurrencyChatsList(ctx context.Context, req *proto.CurrencyChats, rsp *proto.CurrencyChatsListResponse) error {
	data := new(model.Chats).List(req.OrderId)
	if data == nil {
		return nil
	}

	listLen := len(data)
	listData := make([]*proto.CurrencyChats, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.CurrencyChats{
			Id:          data[i].Id,
			OrderId:     data[i].OrderId,
			IsOrderUser: data[i].IsOrderUser,
			Uid:         data[i].Uid,
			Uname:       data[i].Uname,
			Content:     data[i].Content,
			CreatedTime: data[i].CreatedTime,
		}
		listData[i] = adsLists
	}

	rsp.Data = listData
	return nil
}

// 获取用户虚拟货币资产
func (s *RPCServer) GetUserCurrency(ctx context.Context, req *proto.UserCurrencyRequest, rsp *proto.UserCurrency) error {
	data := new(model.UserCurrency).Get(req.Id, req.Uid, req.TokenId)
	if data == nil {
		return nil
	}
	rsp.Id = data.Id
	rsp.Uid = data.Uid
	rsp.TokenId = data.TokenId
	rsp.TokenName = data.TokenName
	rsp.Freeze = data.Freeze
	rsp.Balance = data.Balance
	rsp.Address = data.Address
	rsp.Version = data.Version
	return nil
}

// 获取当前法币账户余额
func (s *RPCServer) GetCurrencyBalance(ctx context.Context, req *proto.GetCurrencyBalanceRequest, rsp *proto.OtherResponse) error {
	balance, err := new(model.UserCurrency).GetBalance(req.Uid, req.TokenId)
	if err != nil {
		rsp.Data = string("0.00")
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	}else{
		rsp.Data = convert.Int64ToStringBy8Bit(balance.Balance)
		rsp.Code = errdefine.ERRCODE_SUCCESS
		return nil
	}


}


// 获取get售价
func (s *RPCServer) GetSellingPrice(ctx context.Context, req *proto.SellingPriceRequest, rsp *proto.OtherResponse) error {
	//
	sellingPriceMap := map[uint32]float64{2:48999.00, 3: 3003.34, 1: 7.08}     // 1 ustd, 2 btc, 3 eth, 4, SDC(平台币)
	key := req.TokenId
	type SellingPrice struct {
		Price float64
	}
	if v, ok := sellingPriceMap[key]; ok {
		dt := SellingPrice{Price:v}
		data, _ := json.Marshal(dt)
		rsp.Data = string(data)
		rsp.Code = errdefine.ERRCODE_SUCCESS
	} else {
		fmt.Println("Key Not Found")
		dt := SellingPrice{Price:v}
		data, _ := json.Marshal(dt)
		rsp.Data = string(data)
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		rsp.Message = "not found!"
		//rsp.Message = "not found!
	}
	return nil
}


// get GetUserRating
// 获取用戶评级
func (s *RPCServer) GetUserRating(ctx context.Context, req *proto.GetUserRatingRequest, rsp *proto.OtherResponse) error {
	uCurrencyCount := new(model.UserCurrencyCount)
	data, err := uCurrencyCount.GetUserCount(req.Uid)
	if err != nil {
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return err
	}
	rData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	}
	rsp.Data = string(rData)
	rsp.Code = errdefine.ERRCODE_SUCCESS
	return nil
}

