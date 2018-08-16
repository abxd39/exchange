package handler

import (
	"digicon/currency_service/model"
	"digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	//"github.com/gin-gonic/gin/json"
	"digicon/common/convert"
	"digicon/common/errors"
	"digicon/common/random"
	"digicon/currency_service/rpc/client"
	"digicon/currency_service/utils"
	"encoding/json"
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
	//fmt.Println(req.TokenId, req.TokenName)

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
	fmt.Println("update req:", req.IsTwolevel, req.Id, req.MinLimit, req.MaxLimit)

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
	//fmt.Println(req.StatusId)
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
			Premium:     data[i].Premium,
			States:      data[i].States,
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
			States:      data[i].States,
			Premium:     data[i].Premium,
			//Balance:     data[i].Balance,
			//Freeze:      data[i].Freeze,
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
	//fmt.Println(req.Id)
	//data := new(model.Tokens).Get(req.Id, req.Name)
	data := new(model.CommonTokens).Get(req.Id, req.Name)
	if data == nil {
		return nil
	}
	//fmt.Println("data:", data)

	rsp.Id = data.Id
	rsp.CnName = data.Name
	rsp.Name = data.Mark
	//rsp.CnName = data.CnName

	return nil
}

// 获取货币类型列表
func (s *RPCServer) CurrencyTokensList(ctx context.Context, req *proto.CurrencyTokensRequest, rsp *proto.CurrencyTokensListResponse) error {

	//data := new(model.Tokens).List()
	data := new(model.CommonTokens).List()
	//fmt.Println("data:", data)
	if data == nil {
		return nil
	}

	//fmt.Println("data:", data)
	listLen := len(data)
	listData := make([]*proto.CurrencyTokens, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.CurrencyTokens{
			Id:     data[i].Id,
			Name:   data[i].Mark,
			CnName: data[i].Name,
			//
			//Name:   data[i].Name,
			//CnName: data[i].CnName,
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
func (s *RPCServer) GetUserCurrencyDetail(ctx context.Context, req *proto.UserCurrencyRequest, rsp *proto.UserCurrency) error {
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
	rsp.Valuation = 0 // 汇率转化
	return nil
}

func (s *RPCServer) GetUserCurrency(ctx context.Context, req *proto.UserCurrencyRequest, rsp *proto.OtherResponse) error {
	//fmt.Println(req.TokenId, req.NoZero)
	data, err := new(model.UserCurrency).GetUserCurrency(req.Uid, req.NoZero)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_USER_BALANCE
		return err
	}

	tkconfig := new(model.TokenConfigTokenCNy)

	fmt.Println("data:", data)

	var symbols []string
	var nosymbol []int

	otherSymbolMap := make(map[uint32]string)
	for _, dt := range data {
		if dt.TokenName != "BTC" {
			if dt.TokenName != "" {
				symbol := fmt.Sprintf("BTC/%s", dt.TokenName)
				symbols = append(symbols, symbol)
			} else {
				nosymbol = append(nosymbol, int(dt.TokenId))
			}
		}
	}

	if len(nosymbol) > 0 {
		tokensModel := new(model.CommonTokens)
		otherTokens := tokensModel.GetByTokenIds(nosymbol)
		for _, otk := range otherTokens {
			fmt.Println(otk.Mark)
			symbol := fmt.Sprintf("BTC/%s", otk.Mark)
			symbols = append(symbols, symbol)
			otherSymbolMap[otk.Id] = fmt.Sprintf("%s", otk.Mark)
		}
	}

	fmt.Println("symbols:", symbols)

	type RespBalance struct {
		Id        uint64  `json:"id"`
		Uid       uint64  `json:"uid"`
		TokenId   uint32  `json:"token_id"`
		TokenName string  `json:"token_name"`
		Address   string  `json:"address"`
		Freeze    string `json:"freeze"`
		Balance   string `json:"balance"`
		Valuation string `json:"valuation"`
	}
	var RespUCurrencyList []RespBalance
	type RespData struct {
		UCurrencyList []RespBalance
		Sum           string `json:"sum"`
		SumCNY        string `json:"sum_cny"`
	}

	var sum int64
	var sumcny int64

	symbolData, err := client.InnerService.UserSevice.CallGetSymbolsRate(symbols)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		rsp.Data = "{}"
		//return err
	} else {
		commontk := new(model.CommonTokens)
		BtcToken := commontk.Get(0, "BTC")
		btcTokenId := BtcToken.Id
		err = tkconfig.GetPrice(btcTokenId)
		var btcConfigPrice int64
		if err != nil {
			log.Println("get btc price error:", err)
			btcConfigPrice = 0
		} else {
			btcConfigPrice = tkconfig.Price
		}
		for _, dt := range data {
			var tmp RespBalance
			var valuation string
			if dt.TokenName == "BTC" {
				if btcConfigPrice <= 0 {
					sumcny += 0
				} else {
					int64valuetion := convert.Int64MulInt64By8Bit(btcConfigPrice, dt.Balance)
					valuation = utils.Int64ToStringBy8Bit(int64valuetion)
					sumcny += int64valuetion
				}
				sum += dt.Balance
			} else {
				var symbol string
				if dt.TokenName != "" {
					symbol = fmt.Sprintf("BTC/%s", dt.TokenName)
				} else {
					symbol = fmt.Sprintf("BTC/%s", otherSymbolMap[dt.TokenId])
				}

				//fmt.Println("symbol:", symbolData)
				symPrice := symbolData.Data[symbol]
				fmt.Println("symPrice:", symPrice, " symbol: ", symbol)
				if symPrice != nil {
					int64price, _ := convert.StringToInt64By8Bit(symPrice.Price)
					if int64price > 0 {
						sum += convert.Int64DivInt64By8Bit(dt.Balance, int64price)
						int64cynPrice := convert.Int64DivInt64By8Bit(btcConfigPrice, int64price)
						if int64cynPrice > 0 {
							int64Valueation := convert.Int64MulInt64By8Bit(dt.Balance, int64cynPrice)
							valuation = utils.Int64ToStringBy8Bit(int64Valueation)
							sumcny += int64Valueation
						} else {
							sumcny += 0
						}
					} else {
						sum += 0
						sumcny += 0
					}
				} else {
					sum += 0
					sumcny += 0
				}
			}

			if dt.TokenName != "" {
				tmp.TokenName = dt.TokenName
			} else {
				tmp.TokenName = otherSymbolMap[dt.TokenId]
			}

			tmp.TokenId = dt.TokenId
			tmp.Id = dt.Id
			tmp.Uid = dt.Uid
			tmp.Address = dt.Address
			tmp.Freeze = convert.Int64ToStringBy8Bit(dt.Freeze)
			tmp.Balance = convert.Int64ToStringBy8Bit(dt.Balance)
			tmp.Valuation = valuation
			RespUCurrencyList = append(RespUCurrencyList, tmp)
		}
	}

	var respdata RespData
	respdata.UCurrencyList = RespUCurrencyList
	respdata.Sum = convert.Int64ToStringBy8Bit(sum)
	respdata.SumCNY = utils.Int64ToStringBy8Bit(sumcny)

	result, err := json.Marshal(respdata)
	if err != nil {
		fmt.Println(err)
		rsp.Data = "{}"
		rsp.Message = err.Error()
		return err
	}

	rsp.Data = string(result)
	return nil
}

// 获取当前法币账户余额
func (s *RPCServer) GetCurrencyBalance(ctx context.Context, req *proto.GetCurrencyBalanceRequest, rsp *proto.OtherResponse) error {
	balance, err := new(model.UserCurrency).GetBalance(req.Uid, req.TokenId)
	sumLimit, err := new(model.Ads).GetUserAdsLimit(req.Uid, req.TokenId)
	//fmt.Println("sumLimit: ", sumLimit)
	var resultBalance int64
	if balance.Balance > sumLimit {
		resultBalance = balance.Balance - sumLimit
	} else {
		resultBalance = 0
	}

	if err != nil {
		rsp.Data = string("0.00")
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	} else {
		//rsp.Data = convert.Int64ToStringBy8Bit(balance.Balance)
		rsp.Data = convert.Int64ToStringBy8Bit(resultBalance)
		rsp.Code = errdefine.ERRCODE_SUCCESS
		return nil
	}
}

// 获取get售价
func (s *RPCServer) GetSellingPrice(ctx context.Context, req *proto.SellingPriceRequest, rsp *proto.OtherResponse) error {
	//
	//sellingPriceMap := map[uint32]float64{2: 48999.00, 3: 3003.34, 1: 7.08} // 1 ustd, 2 btc, 3 eth, 4, SDC(平台币)
	var maxmaxprice float64
	var minminprice float64
	var price float64
	maxPrice, minPrice, _ := new(model.Ads).GetOnlineAdsMaxMinPrice(req.TokenId)
	fmt.Println("maxPrice:", maxPrice, "minPrice: ", minPrice)
	if maxPrice != 0 {
		maxmaxprice = convert.Int64ToFloat64By8Bit(maxPrice)
	}
	if minPrice != 0 {
		minminprice = convert.Int64ToFloat64By8Bit(minPrice)
	}
	if maxPrice == 0 || minPrice == 0 {
		tokenConfigCny := new(model.TokenConfigTokenCNy)
		err := tokenConfigCny.GetPrice(req.TokenId)
		if err != nil {
			log.Println(err.Error())
			price = 0.00
		} else {
			price = convert.Int64ToFloat64By8Bit(tokenConfigCny.Price)
		}
		if minPrice == 0 {
			minminprice = price
		}
		if maxPrice == 0 {
			maxmaxprice = price
		}
	}

	type SellingPrice struct {
		Cny    float64
		Mincny float64
	}

	dt := SellingPrice{Cny: maxmaxprice, Mincny: minminprice}
	data, _ := json.Marshal(dt)
	rsp.Data = string(data)
	rsp.Code = errdefine.ERRCODE_SUCCESS

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

	authResp, err := client.InnerService.UserSevice.CallGetAuthInfo(req.Uid)
	//fmt.Println("authResp:", authResp)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(authResp)
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return err
	}
	type AuthInfo struct {
		EmailAuth     int32  `json:"email_auth"`     //
		PhoneAuth     int32  `json:"phone_auth"`     //
		RealName      int32  `json:"real_name"`      //
		TwoLevelAuth  int32  `json:"two_level_auth"` //
		NickName      string `json:"nick_name"`
		HeadSculpture string `json:"head_scul"`
		CreatedTime   string `json:"created_time"`
	}
	type UserRateAndAuth struct {
		model.UserCurrencyCount
		AuthInfo
		CompleteRate float64 `json:"complete_rate"` //  完成率
		MonthRate    int64   `json:"month_rate"`    // 30日成单
		AverageTo    int64   `json:"average_to"`    // 120 分钟
	}
	var authInfo AuthInfo
	if err = json.Unmarshal([]byte(authResp.Data), &authInfo); err != nil {
		fmt.Println(err)
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return err
	}

	uOrder := new(model.Order)
	now := time.Now()

	startTime := now.Add(-30 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	endTime := now.Format("2006-01-02 15:04:05")

	//fmt.Println("startTime:", startTime, " endTime:",endTime)

	Orders, err := uOrder.GetOrderByTime(req.Uid, startTime, endTime)
	var monthrate int64
	orderLen := len(Orders)
	if err != nil {
		fmt.Println(err.Error())
		monthrate = 0
	} else {
		monthrate = int64(orderLen)
	}

	var allminute int64
	for _, od := range Orders {
		fmt.Println(model.GetHourDiffer(od.CreatedTime, od.ConfirmTime.String))
		allminute = allminute + model.GetHourDiffer(od.CreatedTime, od.ConfirmTime.String)

	}

	rateAndAuth := new(UserRateAndAuth)
	var averageto int64
	if orderLen <= 0 {
		averageto = 0
	} else {
		averageto = allminute / int64(orderLen)
	}
	rateAndAuth.AverageTo = int64(averageto)
	rateAndAuth.MonthRate = monthrate

	rateAndAuth.Uid = data.Uid
	rateAndAuth.Success = data.Success
	rateAndAuth.Failure = data.Failure
	rateAndAuth.Good = data.Good
	rateAndAuth.Cancel = data.Cancel
	rateAndAuth.Orders = data.Orders
	if data.Orders <= 0 {
		rateAndAuth.CompleteRate = 100.0
	} else {
		rateAndAuth.CompleteRate = float64((data.Success / data.Orders) * 100)
	}

	rateAndAuth.RealName = authInfo.RealName
	rateAndAuth.TwoLevelAuth = authInfo.TwoLevelAuth
	rateAndAuth.EmailAuth = authInfo.EmailAuth
	rateAndAuth.PhoneAuth = authInfo.PhoneAuth
	rateAndAuth.NickName = authInfo.NickName
	rateAndAuth.HeadSculpture = authInfo.HeadSculpture
	rateAndAuth.CreatedTime = authInfo.CreatedTime
	fmt.Println("rateAndAuth:", rateAndAuth)

	if rateAndAuth.HeadSculpture == "" {
		rateAndAuth.HeadSculpture = random.GetRandHead()
	}

	rData, err := json.Marshal(rateAndAuth)
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

/*
	AddUserBalance
*/
func (s *RPCServer) AddUserBalance(ctx context.Context, req *proto.AddUserBalanceRequest, rsp *proto.OtherResponse) error {
	uCurrency, err := new(model.UserCurrency).GetBalance(req.Uid, req.TokenId)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		rsp.Message = "get balance error!"
		return err
	}
	intAmount, err := convert.StringToInt64By8Bit(req.Amount)
	if err != nil {
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		rsp.Message = "amount strint convert to int64 err!"
		return err
	}
	err = uCurrency.SetBalance(req.Uid, req.TokenId, intAmount)
	if err != nil {
		fmt.Println(err.Error())
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		rsp.Message = "set balance error!"
		return err
	} else {
		rsp.Code = errdefine.ERRCODE_SUCCESS
		return nil
	}
}

/*
	获取最新交易价格
*/
func (s *RPCServer) GetRecentTransactionPrice(ctx context.Context, req *proto.GetRecentTransactionPriceRequest, rsp *proto.OtherResponse) error {
	//fmt.Println(req.PriceType)
	type TransactionPrice struct {
		MarketPrice float64 `json:"market_price"`
		LatestPrice float64 `json:"latest_price"`
	}
	tp := TransactionPrice{}

	ctk := new(model.CommonTokens)
	fctk := ctk.Get(0, "BTC")
	tokenId := fctk.Id
	tctcy := new(model.TokenConfigTokenCNy)
	tctcy.GetPrice(uint32(tokenId))
	price := tctcy.Price
	tp.MarketPrice = utils.Round2(convert.Int64ToFloat64By8Bit(price), 2)

	chistory := new(model.UserCurrencyHistory)

	err, price := chistory.GetLastPrice(tokenId)
	//fmt.Println("last price:", price)
	if err != nil {
		tp.LatestPrice = tp.MarketPrice
		//tp.LatestPrice = 0.00
	} else {
		tp.LatestPrice = utils.Round2(convert.Int64ToFloat64By8Bit(price), 2)
	}

	data, err := json.Marshal(tp)
	if err != nil {
		fmt.Println(err)
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return err
	} else {
		rsp.Data = string(data)
		rsp.Code = errdefine.ERRCODE_SUCCESS
		return nil
	}
}

// 获取需要广告显示的币种
func (s *RPCServer) DisplayCurrencyTokens(ctx context.Context, req *proto.CurrencyTokensRequest, rsp *proto.CurrencyTokensListResponse) error {
	data := new(model.CommonTokens).DisplayCurrencyTokens()
	//fmt.Println("data:", data)
	if data == nil {
		return nil
	}
	listLen := len(data)
	listData := make([]*proto.CurrencyTokens, listLen)
	for i := 0; i < listLen; i++ {
		adsLists := &proto.CurrencyTokens{
			Id:     data[i].Id,
			Name:   data[i].Mark,
			CnName: data[i].Name,
		}
		listData[i] = adsLists
	}
	rsp.Data = listData
	return nil
}



// 法币划转到代币
func (s *RPCServer) TransferToToken(ctx context.Context, req *proto.TransferToTokenRequest, rsp *proto.OtherResponse) error {
	userCurrencyModel := &model.UserCurrency{}
	err := userCurrencyModel.TransferToToken(req.Uid, int(req.TokenId), int64(req.Num))
	if err != nil {
		rsp.Code = int32(errors.GetErrStatus(err))
		rsp.Message = errors.GetErrMsg(err)
		return nil
	}
	// 划转成功后，检查广告是否需要下架
	go model.AutoDownlineUserAds(req.Uid, uint64(req.TokenId))

	return nil
}



