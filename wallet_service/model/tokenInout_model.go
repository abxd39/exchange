package models

import (
	"digicon/common/model"
	. "digicon/proto/common"
	"digicon/wallet_service/utils"
	"github.com/shopspring/decimal"
	"time"
	log "github.com/sirupsen/logrus"
	"digicon/common/encryption"
	"digicon/common/errors"
)

// 平台币账户的交易
// on neibu
type TokenInout struct {
	Id          int       `xorm:"not null pk INT(11)"`
	Uid         int       `xorm:"comment('用户id') INT(11)"`
	Opt         int       `xorm:"opt"`
	Txhash      string    `xorm:"not null comment('交易hash') VARCHAR(200)"`
	From        string    `xorm:"not null comment('打款方') VARCHAR(42)"`
	To          string    `xorm:"not null comment('收款方') VARCHAR(42)"`
	Amount      int64     `xorm:"not null comment('金额') BIGINT(20)"`
	Fee         int64     `xorm:"fee"`
	Value       string    `xorm:"comment('原始16进制转账数据') VARCHAR(32)"`
	Chainid     int       `xorm:"not null comment('链id') INT(11)"`
	Contract    string    `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
	Tokenid     int       `xorm:"not null comment('币种id') INT(11)"`
	States      int       `xorm:"comment('未处理0，已处理1') TINYINT(1)"`
	TokenName   string    `xorm:"comment('币种名称') VARCHAR(10)"`
	CreatedTime time.Time `xorm:"comment('创建时间') TIMESTAMP"`
	DoneTime    time.Time `xorm:"done_time"`
	Remarks     string    `xorm:"remarks"`
	AmountCny   int64     `xorm:"amount_cny"`
	FeeCny      int64     `xorm:"fee_cny"`
	Gas       int64    `xorm:"comment('gas数量') VARCHAR(30)"`
	Gas_price       int64    `xorm:"comment('gas价格,单位：wei') VARCHAR(30)"`
	Real_fee       int64    `xorm:"comment('实际消耗手续费') VARCHAR(30)"`
}

type User struct {
	Uid              uint64 `xorm:"not null pk autoincr comment('用户ID') BIGINT(11)"`
	Account          string `xorm:"comment('账号') unique VARCHAR(64)"`
	Pwd              string `xorm:"comment('密码') VARCHAR(255)"`
	Country          string `xorm:"comment('地区号') VARCHAR(32)"`
	Phone            string `xorm:"comment('手机') unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"comment('邮箱验证时间') INT(11)"`
	GoogleVerifyId   string `xorm:"comment('谷歌私钥') VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"comment('谷歌验证时间') INT(255)"`
	SmsTip           bool   `xorm:"default 0 comment('短信提醒') TINYINT(1)"`
	PayPwd           string `xorm:"comment('支付密码') VARCHAR(255)"`
	NeedPwd          bool   `xorm:"comment('免密设置1开启0关闭') TINYINT(1)"`
	NeedPwdTime      int    `xorm:"comment('免密周期') INT(11)"`
	Status           int    `xorm:"default 0 comment('用户状态，1正常，2冻结') INT(11)"`
	SecurityAuth     int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	SetTardeMark     int    `xorm:"comment('资金密码设置状态标识') INT(8)"`
}


type SumTokenOut struct {
	TotalFee   int64   `json:"total_fee"`
	Total      int64   `json:"total"`
}

type SumTokenIn struct {
	TotalPut   int64 `json:"total_put"`
}


func (this *TokenInout) Insert(txhash, from, to, total, contract string, chainid int, uid int, tokenid int, tokenname string, decim int,opt int,gas int64,gasprice int64,realfee int64) (int, error) {
	this.Id = 0
	this.Txhash = txhash
	this.From = from
	this.Opt = opt
	this.To = to
	this.Value = total
	amount,err := decimal.NewFromString(total)
	if err != nil {
		log.Info("error",err)
	}
	this.Amount = amount.IntPart()
	this.Contract = contract
	this.Chainid = chainid
	this.Tokenid = tokenid
	this.TokenName = tokenname
	this.Uid = uid
	this.Gas = gas
	this.Gas_price = gasprice
	this.Real_fee = realfee
	affected, err := utils.Engine_wallet.InsertOne(this)
	return int(affected), err
}

/// BTC insert
func (this *TokenInout) BtcInsert(txhash, from, to, tokenName string, amount int64, chainId, tokenId, states, uid int) (int, error) {
	this.Txhash = txhash
	this.From = from
	this.To = to
	this.Amount = amount
	this.Value = ""
	this.Chainid = chainId
	this.Contract = ""
	this.Tokenid = tokenId
	this.States = states
	this.TokenName = tokenName
	this.Uid = uid
	affected, err := utils.Engine_wallet.InsertOne(this)
	return int(affected), err

}

//更新比特币申请提币hash
func (this *TokenInout) UpdateApplyTiBi(applyid int, txhash string) (int, error) {
	affected, err := utils.Engine_wallet.Id(applyid).Update(TokenInout{Txhash: txhash}) //提币已经提交区块链，修改交易hash和正在提币状态
	return int(affected), err
}

//更新比特币申请提币hash
func (this *TokenInout) UpdateApplyTiBi2(applyid int,states int) (int, error) {
	affected, err := utils.Engine_wallet.Id(applyid).Update(TokenInout{States:states})
	return int(affected), err
}

//更新提币完成状态
func (this *TokenInout) BteUpdateAppleDone(txhash string) (int, error) {
	affected, err := utils.Engine_wallet.Where("txhash = ?", txhash).Update(TokenInout{States: 2, DoneTime: time.Now()}) //提币已经完成，修改状态和完成时间
	return int(affected), err
}

// 列表
func (this *TokenInout) GetInOutList(pageIndex, pageSize int, filter map[string]interface{}) (*model.ModelList, []*TokenInout, error) {
	engine := utils.Engine_wallet
	query := engine.Desc("id")

	// 筛选
	query.Where("1=1")
	if v, ok := filter["uid"]; ok {
		query.And("uid=?", v)
	}
	if v, ok := filter["opt"]; ok {
		query.And("opt=?", v)
	}

	tempQuery := query.Clone()
	count, err := tempQuery.Count(this)
	if err != nil {
		return nil, nil, err
	}

	// 获取分页
	offset, modelList := model.Paging(pageIndex, pageSize, int(count))

	var list []*TokenInout
	err = query.Select("*").Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, err
	}
	modelList.Items = list

	return modelList, list, nil
}

//提币申请
func (this *TokenInout) TiBiApply(uid int,tokenid int,to string,amount string,fee string,amountCny int64,feeCny int64,from_address string) (ret int,err error) {
	//查询form地址
	var walletToken = new(WalletToken)
	err = walletToken.GetByUid(uid)
	if err != nil {
		return ERRCODE_UNKNOWN,err
	}

	//根据token_id获取token_name
	var tokenData = new(Tokens)
	_,err = tokenData.GetByid(tokenid)
	if err != nil {
		return
	}


	from := from_address //walletToken.Address

	this.From = from
	this.To = to

	t1,err := decimal.NewFromString(amount)
	if err != nil {
		return 0,errors.New("格式转换失败"+amount)
	}
	t1_c := decimal.NewFromFloat(float64(100000000))
	amount1 := t1.Mul(t1_c).IntPart()
	this.Amount = amount1

	t2,err := decimal.NewFromString(fee)
	if err != nil {
		return 0,errors.New("格式转换失败"+amount)
	}
	t2_c := decimal.NewFromFloat(float64(100000000))
	fee1 := t2.Mul(t2_c).IntPart()
	this.Fee = fee1


	this.Contract = walletToken.Contract
	this.Chainid = walletToken.Chainid
	this.Tokenid = tokenid
	this.Uid = uid
	this.TokenName = tokenData.Mark
	this.AmountCny = amountCny
	this.FeeCny = feeCny
	this.CreatedTime = time.Now()
	this.Contract = tokenData.Contract
	this.Opt = 2 //提币
	this.States = 1  //正在提币
	this.AmountCny = amountCny
	this.FeeCny = feeCny
	affected, err := utils.Engine_wallet.InsertOne(this)
	return int(affected), err

}

//验证支付密码
func (this *TokenInout) AuthPayPwd(uid int32,password string) (ret int32,err error) {
	engine := utils.Engine_common
	var data = new(User)
	engine.ShowSQL(true)
	ok,err := engine.Where("uid=?",uid).Get(data)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !ok {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	if data.PayPwd != encryption.GenMd5AndReverse(password) {
		return ERRCODE_UNKNOWN,nil
	}
	return ERRCODE_SUCCESS,nil
}

//取消提币
func (this *TokenInout) CancelTiBi(uid,id int) (int,error) {
	affected, err := utils.Engine_wallet.Where("uid = ? and id = ?",uid,id).Update(TokenInout{States:3})  //提币已取消
	return int(affected),err
}

//查询申请的提币单
func (this *TokenInout) GetApplyInOut(uid int,id int) (bool,error) {
	return utils.Engine_wallet.Where("uid = ? and id = ?",uid,id).Limit(1).Get(this)
}

func (this *TokenInout) TxhashExist(hash string, chainid int) (bool, error) {
	utils.Engine_wallet.ShowSQL(false)
	return utils.Engine_wallet.Where("txhash=?", hash).Get(this)

}

//根据交易hash，查询数据
func (this *TokenInout) GetByHash(txhash string) error {
	_, err := utils.Engine_wallet.Where("txhash =?", txhash).Get(this)
	if err != nil {
		return err
	}
	return nil
}


/*
	获取所有今天的转账
*/
func (this *TokenInout) GetInOutByTokenIdByTime(tkid uint32, startTime, endTime string) (tokensIntout []TokenInout, err error){
	err = utils.Engine_wallet.Table("token_inout").
		Where("tokenid=? AND created_time >= ? AND created_time <= ? AND states=2", tkid, startTime, endTime).
		Find(&tokensIntout)
	return
}

/*
	提币累计总金额
*/
func (this *TokenInout) GetOutSumByTokenId(tkid uint32, endTime string ) (outsum SumTokenOut, err error) {
	sql := "select sum(amount) as total, sum(fee) as total_fee from token_inout where tokenid=? and opt=2 and created_time < ?"
	_, err = utils.Engine_wallet.Table("token_inout").SQL(sql, tkid, endTime).Get(&outsum)
	return
}

/*
	充币累计总额
*/
func (this *TokenInout) GetInSumByTokenId(tkid uint32, endTime string)(insum SumTokenIn, err error) {
	sql := "select sum(amount) as total_put from token_inout where tokenid=? and opt=1 and created_time < ?"
	_, err = utils.Engine_wallet.Table("token_inout").SQL(sql, tkid,  endTime).Get(&insum)
	return
}

