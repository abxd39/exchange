package handler

import (
	"context"
	"digicon/common/encryption"
	//log "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"digicon/currency_service/model"
	"digicon/proto/common"
	proto "digicon/proto/rpc"
	"encoding/json"
	"fmt"
	"time"

	"digicon/common/convert"
	"digicon/currency_service/rpc/client"
	"strconv"
)

// 获取订单列表
func (s *RPCServer) OrdersList(ctx context.Context, req *proto.OrdersListRequest, rsp *proto.OrdersListResponse) error {
	result := []model.OrderModel{}
	o := new(model.Order)
	rsp.Total, rsp.Page, rsp.PageNum, rsp.Err = o.List(req.Page, req.PageNum, req.AdType, req.States, req.Id, req.Uid, req.TokenId, req.StartTime, req.EndTime, &result)

	orders, err := json.Marshal(result)
	if err != nil {
		log.Errorln(err.Error())
		rsp.Orders = "[]"
		rsp.Message = err.Error()
		return err
	}
	rsp.Orders = string(orders)
	return nil
}

// 取消订单
func (s *RPCServer) CancelOrder(ctx context.Context, req *proto.CancelOrderRequest, rsp *proto.OrderResponse) error {
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	log.Println(req)
	code, msg := new(model.Order).Cancel(req.Id, req.CancelType, updateTimeStr, req.Uid)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 删除订单
func (s *RPCServer) DeleteOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	//fmt.Println(req.Id)
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Delete(req.Id, updateTimeStr)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 确认放行
func (s *RPCServer) ConfirmOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Confirm(req.Id, updateTimeStr, req.Uid)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 待放行
func (s *RPCServer) ReadyOrder(ctx context.Context, req *proto.OrderRequest, rsp *proto.OrderResponse) error {
	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Ready(req.Id, updateTimeStr)
	rsp.Code = code
	rsp.Message = msg
	return nil
}

// 添加订单
func (s *RPCServer) AddOrder(ctx context.Context, req *proto.AddOrderRequest, rsp *proto.OrderResponse) error {
	od := new(model.Order)
	if err := json.Unmarshal([]byte(req.Order), &od); err != nil {
		log.Println(err.Error())
		fmt.Println(err.Error())
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return nil
	}

	ads := new(model.Ads).Get(od.AdId)
	//var nowAds *model.Ads
	//nowAds = ads.Get(od.AdId)
	if ads == nil {
		rsp.Code = errdefine.ERRCODE_ADS_NOTEXIST
		return nil
	}

	fmt.Println("ads is two level: ", ads.IsTwolevel)
	if ads.IsTwolevel == 1 {
		authResp, err := client.InnerService.UserSevice.CallGetAuthInfo(uint64(req.Uid))
		if err != nil {
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
		var authInfo AuthInfo
		if err = json.Unmarshal([]byte(authResp.Data), &authInfo); err != nil {
			fmt.Println(err)
			rsp.Code = errdefine.ERRCODE_ADS_NEED_TWO_LEVEL
			return nil
		}
		fmt.Println("two level auth: ", req.Uid, authInfo.TwoLevelAuth, authInfo.TwoLevelAuth == 1)
		if authInfo.TwoLevelAuth != 1 {
			msg := "没有两次验证"
			fmt.Println(msg)
			log.Println(msg)
			rsp.Code = errdefine.ERRCODE_ADS_NEED_TWO_LEVEL
			//err := errors.New(msg)
			return nil
		}
	}

	od.AdType = ads.TypeId
	od.Price = int64(ads.Price)
	od.TokenId = uint64(ads.TokenId)

	if uint32(ads.TypeId) == 2 { //   广告状态为2(购买),那么当前用户肯定为出售
		od.SellId = ads.Uid
		od.BuyId = uint64(req.Uid)
	} else {
		od.BuyId = ads.Uid
		od.SellId = uint64(req.Uid)
	}

	od.PayId = ads.Pays

	//fmt.Println("od.selleid:", od.SellId, od.BuyId)
	if od.SellId == od.BuyId {
		msg := "无法下自己订单"
		//err := errors.New(msg)
		log.Errorln(msg)
		fmt.Println(msg)
		rsp.Code = errdefine.ERRCODE_TRADE_TO_SELF
		rsp.Message = msg
		return nil
	}

	var uids []uint64
	uids = append(uids, od.SellId)
	uids = append(uids, od.BuyId)

	nickNames, err := client.InnerService.UserSevice.CallGetNickName(uids) // rpc 获取用户信息
	for i := 0; i < 2; i++ {
		if err != nil {
			nickNames, err = client.InnerService.UserSevice.CallGetNickName(uids) // rpc 获取用户信息
		}
	}

	if err != nil {
		fmt.Println(err)
		log.Errorln(err.Error())
	} else {
		nickUsers := nickNames.User
		for i := 0; i < len(nickUsers); i++ {
			if nickUsers[i].Uid == od.SellId {
				od.SellName = nickUsers[i].NickName
			}
			if nickUsers[i].Uid == od.BuyId {
				od.BuyName = nickUsers[i].NickName
			}
		}
	}

	od.OrderId = encryption.CreateOrderId(uint64(req.Uid), int32(od.TokenId))
	od.States = 1

	now := time.Now()
	mm, _ := time.ParseDuration("15m") // 过期时间15分钟
	od.CreatedTime = now.Format("2006-01-02 15:04:05")
	od.UpdatedTime = now.Format("2006-01-02 15:04:05")
	od.ExpiryTime = now.Add(mm).Format("2006-01-02 15:04:05")

	od.NumTotalPrice = convert.Int64MulInt64By8Bit(od.Num, od.Price)
	od.FeePrice = convert.Int64MulInt64By8Bit(od.Fee, od.Price)

	//fmt.Println("od:", od)
	id, code := od.Add(req.Uid)
	rsp.Code = code
	rsp.Data = strconv.FormatUint(id, 10)
	return nil
}

// get Trade detail

func (s *RPCServer) TradeDetail(ctx context.Context, req *proto.TradeDetailRequest, rsp *proto.TradeDetailResponse) error {
	order := new(model.Order)
	aliPay := new(model.UserCurrencyAlipayPay)
	bankPay := new(model.UserCurrencyBankPay)
	paypalPay := new(model.UserCurrencyPaypalPay)
	wechatPay := new(model.UserCurrencyWechatPay)

	ctoken := new(model.CommonTokens)

	order.GetOrder(req.Id)
	sellid := order.SellId
	aliPay.GetByUid(sellid)
	bankPay.GetByUid(sellid)
	paypalPay.GetByUid(sellid)
	wechatPay.GetByUid(sellid)

	tkname := ctoken.Get(uint32(order.TokenId), "")
	tokenName := tkname.Mark

	type Data struct {
		SellId     uint64 `form:"sell_id"                json:"sell_id"`
		SellName   string `form:"sell_name"              json:"sell_name"`
		BuyId      uint64 `form:"buy_id"                 json:"buy_id"`
		BuyName    string `form:"buy_name"                json:"buy_name"`
		States     uint32 `form:"states"                 json:"states"`
		ExpiryTime string `xorm:"expiry_time"            json:"expiry_time" `
		TokenId    uint64 `form:"token_id"              json:"token_id"`
		TokenName  string `form:"token_name"            json:"token_name"`

		OrderId        string `form:"order_id"               json:"order_id"`
		PayPrice       int64  `form:"pay_price"              json:"pay_price"`
		Num            int64  `form:"num"                    json:"num"`
		Price          int64  `form:"price"                  json:"price"`
		AliPayName     string `form:"alipay_name"            json:"alipay_name"`
		Alipay         string `form:"alipay"                 json:"alipay"`
		AliReceiptCode string `form:"ali_receipt_code"       json:"ali_receipt_code"`

		BankpayName string `form:"bankpay_name"            json:"bankpay_name"`
		CardNum     string `form:"card_num"               json:"card_num"`
		BankName    string `form:"bank_name"              json:"bank_name"`
		BankInfo    string `form:"bank_info"              json:"bank_info"`

		WechatName        string `form:"wechat_name"            json:"wechat_name"`
		Wechat            string `form:"wechat"                 json:"wechat"`
		WechatReceiptCode string `form:"wechat_receipt_code"    json:"wechat_receipt_code"`
		PaypalNum         string `form:"paypal_num"             json:"paypal_num"`
	}
	var dt Data
	dt.SellId = order.SellId
	dt.SellName = order.SellName
	dt.BuyId = order.BuyId
	dt.BuyName = order.BuyName
	dt.States = order.States
	dt.ExpiryTime = order.ExpiryTime
	dt.TokenId = order.TokenId
	dt.TokenName = tokenName

	dt.OrderId = order.OrderId
	dt.Price = order.Price
	dt.Num = order.Num
	dt.PayPrice = convert.Int64MulInt64By8Bit(dt.Price, dt.Num)
	dt.AliPayName = aliPay.Name
	dt.Alipay = aliPay.Alipay
	dt.AliReceiptCode = aliPay.ReceiptCode
	dt.BankpayName = bankPay.Name
	dt.BankInfo = bankPay.BankInfo
	dt.CardNum = bankPay.CardNum
	dt.WechatName = wechatPay.Name
	dt.Wechat = wechatPay.Wechat
	dt.WechatReceiptCode = wechatPay.ReceiptCode
	dt.PaypalNum = paypalPay.Paypal

	resultdt, err := json.Marshal(dt)
	if err != nil {
		rsp.Data = ""
		rsp.Code = errdefine.ERRCODE_UNKNOWN
	} else {
		rsp.Data = string(resultdt)
		rsp.Code = errdefine.ERRCODE_SUCCESS
	}
	return nil
}

func (s *RPCServer) GetTradeHistory(ctx context.Context, req *proto.GetTradeHistoryRequest, rsp *proto.OtherResponse) error {
	od := new(model.Order)
	uOrderHistoryList, err := od.GetOrderHistory(req.StartTime, req.EndTime, req.Limit)
	if err != nil {
		log.Errorln(err.Error())
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return err
	}
	ohistlen := len(uOrderHistoryList)
	if ohistlen <= 0 {
		uOrderHistoryList = model.GenerateKline()
	}
	data, err := json.Marshal(uOrderHistoryList)
	rsp.Data = string(data)
	rsp.Code = errdefine.ERRCODE_SUCCESS
	return nil
}

/*
	获取用户资产明细
*/
func (s *RPCServer) GetAssetDetail(ctx context.Context, req *proto.GetAssetDetailRequest, rsp *proto.OtherResponse) error {
	uCurrencyHistory := new(model.UserCurrencyHistory)
	uAssetDeailList, total, page, pageNum, err := uCurrencyHistory.GetAssetDetail(int32(req.Uid), req.Page, req.PageNum)
	if err != nil {
		log.Errorln(err.Error())
		rsp.Code = errdefine.ERRCODE_UNKNOWN
		return err
	}
	var uids []uint64
	var tokenids []int
	for _, ua := range uAssetDeailList {
		uids = append(uids, uint64(ua.TradeUid))
		tokenids = append(tokenids, ua.TokenId)
	}
	nickNames, err := client.InnerService.UserSevice.CallGetNickName(uids) // rpc 获取用户信息

	fmt.Println("nickNames:", nickNames)

	userNameMap := map[uint64]string{}
	nickUsers := nickNames.User
	for i := 0; i < len(nickUsers); i++ {
		userNameMap[nickUsers[i].Uid] = nickUsers[i].NickName
	}

	tokenIdsMap := map[uint32]string{}
	comtoken := new(model.CommonTokens)
	tokenNames := comtoken.GetByTokenIds(tokenids)
	for _, tn := range tokenNames {
		tokenIdsMap[tn.Id] = tn.Mark
	}

	type NewUserCurrencyHisotry struct {
		Id          int     `json:"id"                  `
		Uid         int32   `json:"uid"               `
		TradeUid    int32   `json:"trade_uid"         `
		TokenId     int     `json:"token_id"            `
		TokenName   string  `json:"token_name"`
		Num         float64 `json:"num"                 `
		Operator    int     `json:"operator"            `
		CreatedTime string  `json:"created_time"        `
		TradeName   string  `json:"trade_name"         `
	}

	var NewUAssetDetaillList []NewUserCurrencyHisotry
	for _, ua := range uAssetDeailList {
		//fmt.Println(ua.CreatedTime)
		var tmp NewUserCurrencyHisotry
		tmp.TradeName = userNameMap[uint64(ua.TradeUid)]
		tmp.Uid = ua.Uid
		tmp.TradeUid = ua.TradeUid
		tmp.Num = convert.Int64ToFloat64By8Bit(ua.Num)
		tmp.CreatedTime = ua.CreatedTime
		tmp.TokenId = ua.TokenId
		tmp.TokenName = tokenIdsMap[uint32(ua.TokenId)]
		tmp.Operator = ua.Operator
		NewUAssetDetaillList = append(NewUAssetDetaillList, tmp)
	}
	type ResultData struct {
		NewList []NewUserCurrencyHisotry
		Total   int64  `json:"total"`
		Page    uint32 `json:"page"`
		PageNum uint32 `json:"page_num"`
	}
	resultdt := ResultData{NewList: NewUAssetDetaillList, Total: total, Page: page, PageNum: pageNum}
	data, err := json.Marshal(resultdt)
	rsp.Data = string(data)
	rsp.Code = errdefine.ERRCODE_SUCCESS
	return nil
}
