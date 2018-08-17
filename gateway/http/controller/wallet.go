package controller

import (
	"digicon/common/convert"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type WalletGroup struct {
}

func (this *WalletGroup) Router(router *gin.Engine) {

	r := router.Group("/wallet")
	r.POST("/create", this.Create)       // 创建钱包
	r.POST("/signtx", this.Signtx)       // 签名
	r.POST("/sendrawtx", this.SendRawTx) // 广播
	r.POST("/tibi", this.Tibi)           //
	r.GET("/getvalue", this.GetValue)    // 查询链上余额
	r.POST("/address/save", this.AddressSave)
	r.GET("/address/list", this.AddressList)
	r.POST("/address/delete", this.AddressDelete)

	//btc signtx
	r.POST("/signtx_btc", this.BtcSigntx) // btc 签名
	r.POST("/biti_btc", this.BtcTiBi)     // btc

	r.GET("in_list", this.InList)
	r.GET("out_list", this.OutList)

	r.POST("/tibi_apply", this.TibiApply) //

	r.POST("/tibi_cancel", this.TiBiCancel) //

	r.POST("/get_address", this.GetAddress) // 获取充值地址

	r.POST("/sync_block", this.SyncEthBlockTx)  //同步区块信息
}

///////////////////////// start btc ///////////////////////////

func (this *WalletGroup) BtcSigntx(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid     int32  `form:"uid"          json:"uid"        binding:"required"` // 用户 id
		TokenId int32  `form:"token_id"     json:"token_id"   binding:"required"` // 币种ID
		Address string `form:"address"      json:"address"    binding:"required"` // 要发送给的地址
		Amount  string `form:"amount"       json:"amount"     binding:"required"` //交易总量
		Applyid int32  `form:"apply_id"     json:"amount"     binding:"required"` //申请提币id
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallBtcSigntx(&proto.BtcSigntxRequest{
		Uid:     param.Uid,
		Tokenid: param.TokenId,
		Address: param.Address,
		Amount:  param.Amount,
		Applyid: param.Applyid,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		fmt.Println(rsp.Data)
		//ret.SetDataValue(""rsp.Data)
		ret.SetDataSection("txhash", rsp.Data)
		ret.SetErrCode(rsp.Code, GetErrorMessage(rsp.Code))
	}

}

func (this *WalletGroup) BtcTiBi(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid     int32  `form:"uid"      json:"uid"       binding:"required"`
		TokenId int32  `form:"token_id" json:"token_id"  binding:"required"`
		Address string `form:"address"       json:"address"        binding:"required"`
		Amount  string `form:"amount"   json:"amount"    binding:"required"`
		Applyid int32 `form:"apply_id"     json:"apply_id"     binding:"required"` //申请提币id
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallBtcTibi(&proto.BtcTibiRequest{
		Uid:     param.Uid,
		Tokenid: param.TokenId,
		Address: param.Address,
		Amount:  param.Amount,
		Applyid: param.Applyid,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN,rsp.Message)
		return
	}

	if rsp.Code != 0 {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
	}

	ret.SetErrCode(int32(rsp.Code), GetErrorMessage(int32(rsp.Code)))
	ret.SetDataSection("txhash", rsp.Data)
	return
}

///////////////////////// end btc ///////////////////////////

func (this *WalletGroup) Index(ctx *gin.Context) {

}
func (this *WalletGroup) Create(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()

	type Param struct {
		Uid     int  `form:"uid"      json:"uid"       binding:"required"`
		TokenId int  `form:"token_id" json:"token_id"  binding:"required"`
	}

	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallCreateWallet(param.Uid, param.TokenId)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	if rsp.Code != "0" {
		ret.SetErrCode(ERRCODE_UNKNOWN, rsp.Msg)
		return
	}
	ret.SetDataSection("type", rsp.Data.Type)
	ret.SetDataSection("addr", rsp.Data.Addr)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

func (this *WalletGroup) Signtx(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()

	type Param struct {
		Uid     int  `form:"uid"      json:"uid"       binding:"required"`
		TokenId int  `form:"token_id" json:"token_id"  binding:"required"`
		To string  `form:"to" json:"to"  binding:"required"`
		Gasprice int64  `form:"gasprice" json:"gasprice"  binding:"required"`
		Amount string  `form:"amount" json:"amount"  binding:"required"`
	}

	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallSigntx(param.Uid, param.TokenId, param.To, param.Gasprice, param.Amount)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	if rsp.Code != "0" {
		ret.SetErrCode(ERRCODE_UNKNOWN,rsp.Msg)
		return
	}
	fmt.Println(rsp)
	ret.SetDataSection("signtx", rsp.Data.Signtx)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))

}

func (this *WalletGroup) Update(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()

	rsp, err := rpc.InnerService.WalletSevice.Callhello("eth")
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("data", rsp)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
}

func (this *WalletGroup) SendRawTx(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		TokenId int32  `form:"token_id" binding:"required"`
		Signtx  string `form:"signtx" binding:"required"`
		Applyid int32  `form:"apply_id"     json:"apply_id"     binding:"required"` //申请提币id
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallSendRawTx(param.TokenId, param.Signtx, param.Applyid)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}

	if rsp.Code != "0" {
		ret.SetErrCode(ERRCODE_UNKNOWN,rsp.Msg)
		return
	}
	
	ret.SetDataSection("result", rsp.Data.Result)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

//申请提币
func (this *WalletGroup) Tibi(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid      int32  `form:"uid" binding:"required"`
		Token_id int32  `form:"token_id" binding:"required"`
		To       string `form:"to" binding:"required"`
		Amount   string `form:"amount" binding:"required"`
		Gasprice int32  `form:"gasprice" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallTibi(param.Uid, param.Token_id, param.To, param.Gasprice, param.Amount)
	if err != nil {
		//fmt.Println(rsp.Code, rsp.Msg)
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	//ctx.JSON(http.StatusOK, rsp)
	ret.SetDataSection("msg", rsp.Msg)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

//查询钱包链上额度
func (this *WalletGroup) GetValue(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid      int32 `form:"uid" binding:"required"`
		Token_id int32 `form:"token_id" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallGetValue(param.Uid, param.Token_id)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetDataSection("msg", rsp.Msg)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

func (this *WalletGroup) AddressSave(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid      int32  `form:"uid" binding:"required"`
		Token_id int32  `form:"token_id" binding:"required"`
		Address  string `form:"address" binding:"required"`
		Mark     string `form:"mark" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressSave(param.Uid, param.Token_id, param.Address, param.Mark)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetDataSection("msg", rsp.Msg)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}
func (this *WalletGroup) AddressList(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid int32 `form:"uid"      binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressList(param.Uid)

	if err != nil {
		fmt.Println("ERRCODE_UNKNOWN:", ERRCODE_UNKNOWN)
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	} else {
		fmt.Println("SUCCESS:", ERRCODE_SUCCESS)
		ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
		ret.SetDataSection("data", rsp.Data)
	}

	return
}

func (this *WalletGroup) AddressDelete(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid int32 `form:"uid" binding:"required"`
		Id  int32 `form:"id" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallAddressDelete(param.Uid, param.Id)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("data", rsp.Data)
	ret.SetDataSection("msg", rsp.Msg)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

// 提币列表
func (this *WalletGroup) InList(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()

	param := &struct {
		Uid     int32  `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Page    int32  `form:"page" binding:"required"`
		PageNum int32  `form:"page_num"`
	}{}

	if err := ctx.ShouldBind(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallInList(&proto.InListRequest{
		Uid:     param.Uid,
		Page:    param.Page,
		PageNum: param.PageNum,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code, rsp.Msg)

	// 重组data
	type item struct {
		Id          int32  `json:"id"`
		TokenId     int32  `json:"token_id"`
		TokenName   string `json:"token_name"`
		Amount      string `json:"amount"`
		Address     string `json:"address"`
		States      int32  `json:"states"`
		CreatedTime int64  `json:"create_time"`
	}
	type list struct {
		PageIndex int32   `json:"page_index"`
		PageSize  int32   `json:"page_size"`
		TotalPage int32   `json:"total_page"`
		Total     int32   `json:"total"`
		Items     []*item `json:"items"`
	}

	newItems := make([]*item, len(rsp.Data.Items))
	for k, v := range rsp.Data.Items {
		newItems[k] = &item{
			Id:          v.Id,
			TokenId:     v.TokenId,
			TokenName:   v.TokenName,
			Amount:      convert.Int64ToStringBy8Bit(v.Amount),
			Address:     v.Address,
			States:      v.States,
			CreatedTime: v.CreatedTime,
		}
	}

	newList := &list{
		PageIndex: rsp.Data.PageIndex,
		PageSize:  rsp.Data.PageSize,
		TotalPage: rsp.Data.TotalPage,
		Total:     rsp.Data.Total,
		Items:     newItems,
	}

	ret.SetDataSection("list", newList)
	return
}

// 充币列表
func (this *WalletGroup) OutList(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()

	param := &struct {
		Uid     int32  `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Page    int32  `form:"page" binding:"required"`
		PageNum int32  `form:"page_num"`
	}{}

	if err := ctx.ShouldBind(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.WalletSevice.CallOutList(&proto.OutListRequest{
		Uid:     param.Uid,
		Page:    param.Page,
		PageNum: param.PageNum,
	})
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code, rsp.Msg)

	// 重组data
	type item struct {
		Id          int32  `json:"id"`
		TokenId     int32  `json:"token_id"`
		TokenName   string `json:"token_name"`
		Amount      string `json:"amount"`
		Fee         string `json:"fee"`
		Address     string `json:"address"`
		Remarks     string `json:"remarks"`
		States      int32  `json:"states"`
		CreatedTime int64  `json:"create_time"`
	}
	type list struct {
		PageIndex int32   `json:"page_index"`
		PageSize  int32   `json:"page_size"`
		TotalPage int32   `json:"total_page"`
		Total     int32   `json:"total"`
		Items     []*item `json:"items"`
	}

	newItems := make([]*item, len(rsp.Data.Items))
	for k, v := range rsp.Data.Items {
		newItems[k] = &item{
			Id:          v.Id,
			TokenId:     v.TokenId,
			TokenName:   v.TokenName,
			Amount:      convert.Int64ToStringBy8Bit(v.Amount),
			Fee:         convert.Int64ToStringBy8Bit(v.Fee),
			Address:     v.Address,
			Remarks:     v.Remarks,
			States:      v.States,
			CreatedTime: v.CreatedTime,
		}
	}

	newList := &list{
		PageIndex: rsp.Data.PageIndex,
		PageSize:  rsp.Data.PageSize,
		TotalPage: rsp.Data.TotalPage,
		Total:     rsp.Data.Total,
		Items:     newItems,
	}

	ret.SetDataSection("list", newList)
	return
}

//申请提币2
func (this *WalletGroup) TibiApply(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid        int32  `form:"uid" binding:"required"`
		Token_id   int32  `form:"token_id" binding:"required"`
		To         string `form:"to" binding:"required"`
		Amount     string `form:"amount" binding:"required"`
		Gasprice   string `form:"gasprice" binding:"required"`
		RealAmount string `form:"real_amount" binding:"required"`
		SmsCode    string `form:"sms_code" binding:"required"`
		EmailCode  string `form:"email_code" binding:"required"`
		Password   string `form:"password" binding:"required"`
		Phone   string `form:"phone" binding:"required"`
		Email   string `form:"email" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	rsp, err := rpc.InnerService.WalletSevice.CallTibiApply(param.Uid, param.Token_id, param.To, param.Gasprice, param.Amount, param.RealAmount, param.SmsCode, param.EmailCode, param.Password,param.Phone,param.Email)
	if err != nil {
		//fmt.Println(rsp.Code, rsp.Msg)
		log.Error(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	if rsp.Code != 0 {
		log.Error(rsp.Msg)
		ret.SetErrCode(ERRCODE_UNKNOWN, rsp.Msg)
		return
	}
	//ctx.JSON(http.StatusOK, rsp)
	//ret.SetDataSection("msg", rsp.Msg)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

//撤销提币
func (this *WalletGroup) TiBiCancel(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Uid int32 `form:"uid" binding:"required"`
		Id  int32 `form:"id" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	_, err := rpc.InnerService.WalletSevice.CallCancelTiBi(param.Uid, param.Id)
	if err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

func (this *WalletGroup) GetAddress(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	userid, _ := strconv.Atoi(ctx.PostForm("uid"))
	tokenid, _ := strconv.Atoi(ctx.PostForm("token_id"))

	rsp, err := rpc.InnerService.WalletSevice.CallGetAddress(userid, tokenid)
	if err != nil {
		//ret.SetDataSection("msg", rsp.Msg)
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("type", rsp.Type)
	ret.SetDataSection("addr", rsp.Addr)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

//同步以太坊区块交易
func (this *WalletGroup) SyncEthBlockTx(ctx *gin.Context) {
	ret := NewPublciError()
	defer func() {
		ctx.JSON(http.StatusOK, ret.GetResult())
	}()
	type Param struct {
		Block int32 `form:"block" binding:"required"`
	}
	var param Param
	if err := ctx.ShouldBind(&param); err != nil {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}


	rsp, err := rpc.InnerService.WalletSevice.CallSyncBlockTx(param.Block)
	if err != nil || rsp.Code != 0 {
		log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(rsp.Code, rsp.Msg)
	return
}
