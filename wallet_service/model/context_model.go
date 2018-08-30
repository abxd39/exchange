package models

import "digicon/wallet_service/utils"

type Context struct {
	Id      int    `xorm:"not null pk autoincr INT(11)"`
	Number  int    `xorm:"comment('区块高度') INT(11)"`
	Txs     string `xorm:"comment('包含的交易') TEXT"`
	Hash    string `xorm:"comment('区块hash') VARCHAR(48)"`
	Node    string `xorm:"comment('节点url') VARCHAR(255)"`
	Chainid int    `xorm:"comment('chainid') INT(11)"`
	Type    string `xorm:"comment('项目类型btc，eth') VARCHAR(10)"`
}

func (this *Context) MaxNumber(url string, chainid int) (int, error) {
	this.Id = 0
	ok, err := utils.Engine_wallet.Where("node=? and chainid=?", url, chainid).Get(this)
	if err != nil {
		return 0, nil
	}
	if ok {
		return this.Number, nil
	}
	ok, err = utils.Engine_wallet.Where("chainid=?", chainid).Get(this)
	if err != nil {
		return 0, nil
	}
	if ok {
		return this.Number, nil
	}
	return 0, nil
}
func (this *Context) Save(node string, chainid int, blocknumber int) (int, error) {
	this.Node = node
	this.Chainid = chainid
	this.Number = blocknumber
	utils.Engine_wallet.ShowSQL(false)
	affected, err := utils.Engine_wallet.Where("node=? and chainid=?", node, chainid).Update(this)
	if err != nil {
		return 0, err
	}
	if affected > 0 {
		return int(affected), nil
	}
	affected, err = utils.Engine_wallet.InsertOne(this)
	return int(affected), nil
}

//查询最大区块高度
func (this *Context) GetMaxNumber(token_type string) (int,error) {
	ok, err := utils.Engine_wallet.Where("type = ?",token_type).Get(this)
	if err != nil {
		return 0, nil
	}
	if ok {
		return this.Number, nil
	}
	return 0,nil
}

//保存最大高度
func (this *Context) SaveMaxNumber(node string,chainid int,token_type string,blocknumber int) (int,error) {
	this.Number = blocknumber
	this.Node = node
	this.Chainid = chainid
	this.Type = token_type
	affected, err := utils.Engine_wallet.Where("type = ?",token_type).Update(this)
	if err != nil {
		return 0, err
	}
	if affected > 0 {
		return int(affected), nil
	}
	affected, err = utils.Engine_wallet.InsertOne(this)
	return int(affected), nil
}

