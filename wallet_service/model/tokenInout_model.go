package models

import (
	"digicon/common/model"
	"digicon/wallet_service/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
	. "digicon/proto/common"
	"fmt"
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
	AmountCny         int64     `xorm:"amount_cny"`
	FeeCny         int64     `xorm:"fee_cny"`
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

func (this *TokenInout) Insert(txhash, from, to, value, contract string, chainid int, uid int, tokenid int, tokenname string, decim int) (int, error) {
	this.Id = 0
	this.Txhash = txhash
	this.From = from
	this.To = to
	this.Value = value
	temp, _ := new(big.Int).SetString(value[2:], 16)
	amount := decimal.NewFromBigInt(temp, int32(8-decim)).IntPart()
	this.Amount = amount
	this.Contract = contract
	this.Chainid = chainid
	this.Tokenid = tokenid
	this.TokenName = tokenname
	this.Uid = uid
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
func (this *TokenInout) TiBiApply(uid int,tokenid int,to string,amount string,fee string) (ret int,err error) {
	//查询form地址
	var walletToken = new(WalletToken)
	err = walletToken.GetByUid(uid)
	if err != nil {
		return ERRCODE_UNKNOWN,err
	}
	from := walletToken.Address

	this.From = from
	this.To = to
	temp, _ := new(big.Int).SetString(amount,10)
	amount1 := decimal.NewFromBigInt(temp, int32(8-8)).IntPart()
	this.Amount = amount1

	temp1, _ := new(big.Int).SetString(amount,10)
	fee1 := decimal.NewFromBigInt(temp1, int32(8-8)).IntPart()
	this.Fee = fee1

	this.Contract = walletToken.Contract
	this.Chainid = walletToken.Chainid
	this.Tokenid = tokenid
	this.Uid = uid
	affected, err := utils.Engine_wallet.InsertOne(this)
	fmt.Println("保存结果：",affected,err)
	return int(affected), err

}

//验证支付密码
func (this *TokenInout) AuthPayPwd(uid int32,password string) (ret int32,err error) {
	engine := utils.Engine_common
	var data = new(User)
	ok,err := engine.Where("uid=?",uid).Get(&data)
	fmt.Println("验证资金密码：",ok,err)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !ok {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	if data.PayPwd != password {
		return ERRCODE_UNKNOWN,nil
	}
	return ERRCODE_SUCCESS,nil
}
