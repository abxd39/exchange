package models

import (
	"digicon/common/model"
	"digicon/wallet_service/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
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
