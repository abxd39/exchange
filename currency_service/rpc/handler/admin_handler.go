package handler


import (
	"digicon/currency_service/model"
	"digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	//"github.com/gin-gonic/gin/json"
	"digicon/common/convert"
	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"
	"time"
)



/*
	给后台统计每个人的账户余额的折合
 */
func (s *RPCServer) GetUsersBalance(ctx context.Context, req *proto.GetUserBalanceUids, rsp *proto.UserBalancesResponse) error {
	uCurrency := new(model.UserCurrency)

	var UsersBalance []*proto.UserBalanceOne
	for _,uid := range  req.Uids{
		ucurrens, err  := uCurrency.GetByUid(uint64(uid))
		if err != nil {
			log.Errorln(err.Error())
			fmt.Println(err)
			continue
		}
		var udata proto.UserBalanceOne
		var balance int64
		var fronze int64

		for _, uc := range ucurrens {
			udata.Uid = int64(uc.Uid)
			price := model.GetCnyPrice(int32(uc.TokenId))

			numCny := convert.Int64MulInt64By8Bit(uc.Balance, price)
			feeCny := convert.Int64MulInt64By8Bit(uc.Freeze, price)
			balance += numCny
			fronze += feeCny
			fmt.Println("uid: ",uid, " tokenid:", uc.TokenId, " price:", price, " balance:", uc.Balance, "fee: ", feeCny)
		}

		udata.BalanceCnyInt = balance
		udata.BalanceCny =  fmt.Sprintf("%.3f",convert.Int64ToFloat64By8Bit(balance))
		udata.FrozenCnyInt = fronze
		udata.FrozenCny =   fmt.Sprintf("%.3f",convert.Int64ToFloat64By8Bit(fronze))
		udata.TotalCnyInt = balance + fronze
		udata.TotalCny =  fmt.Sprintf("%.3f", convert.Int64ToFloat64By8Bit(balance + fronze))
		log.Infoln(udata)
		UsersBalance = append(UsersBalance, &udata)
	}
	rsp.Data = UsersBalance
	rsp.Code = errdefine.ERRCODE_SUCCESS
	return nil
}


/*l
	后台放行
*/
func (s *RPCServer) AdminConfirm(ctx context.Context, req *proto.ConfirmOrderRequest, rsp *proto.OrderResponse) (err error){

	updateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	code, msg := new(model.Order).Confirm(req.Id, updateTimeStr, req.Uid)
	rsp.Code = code
	rsp.Message = msg

	return nil
}


 /*
	 后台取消
 */
 func (s *RPCServer) AdminCancel(ctx context.Context, req *proto.CancelOrderRequest, rsp *proto.OrderResponse) (err error){

	 code, err := model.CancelAction(req.Id, 3)    //3为系统取消
	 if err != nil {
	 	rsp.Code = code
	 	return nil
	 }else{
		 rsp.Code = errdefine.ERRCODE_SUCCESS
		 rsp.Message = ""
		 return nil
	 }

 }



