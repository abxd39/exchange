package models

import (
	. "digicon/wallet_service/utils"
	"errors"
	"strings"
)

type Tokens struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	Name      string `xorm:"default '' comment('货币名称') VARCHAR(64)"`
	Detail    string `xorm:"default '' comment('详情地址') VARCHAR(255)"`
	Signature string `xorm:"default '' comment('签名方式') VARCHAR(255)"`
	Chainid   int    `xorm:"default 0 comment('链id') INT(11)"`
	Github    string `xorm:"default '' comment('项目地址') VARCHAR(255)"`
	Web       string `xorm:"default '' VARCHAR(255)"`
	Mark      string `xorm:"comment('英文标识') CHAR(10)"`
	Logo      string `xorm:"default '' comment('货币logo') VARCHAR(255)"`
	Contract  string `xorm:"default '' comment('合约地址') VARCHAR(255)"`
	Node      string `xorm:"default '' comment('节点地址') VARCHAR(100)"`
	Status    int    `xorm:"default '' comment('状态') VARCHAR(64)"`
	Decimal   int    `xorm:"not null default 1 comment('精度 1个eos最小精度的10的18次方') INT(11)"`
	Out_token_fee   float64    `xorm:"not null default 1 comment('提币手续费') BIGINT(11)"`
	Gather_limit   int64    `xorm:"not null default 1 comment('归总数据最低限制') BIGINT(11)"`
}

func (this *Tokens) GetByid(id int) (bool, error) {
	//log.Info("GetByid param id:",id)
	Engine_common.ShowSQL(true)
	session := Engine_common.NewSession()
	boo,err := session.Id(id).Get(this)
	//lastSQL, lastSQLArgs := session.LastSQL()
	//log.Info("last sql:",lastSQL, lastSQLArgs)

	defer session.Close()

	return boo,err
}

func (this *Tokens) GetidByContract(contract string, chainid int) (int, error) {
	this.Id = 0
	ok, err := Engine_common.Where("contract=?", contract).Get(this)
	//fmt.Println(contract,chainid)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("token not exist")
	}
	return this.Id, nil
}

func (this *Tokens) GetDecimal(id int) (int, error) {
	this.Id = id
	ok, err := this.GetByid(id)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("connot find this token")
	}
	return this.Decimal, nil
}

/*
	根据名称获取
*/
func (this *Tokens) GetByName(name string) (bool, error) {
	exists, err := Engine_common.Where("mark = ?", strings.ToUpper(name)).Get(this)
	if err != nil {
		return exists, err
	}
	return exists, nil
}


/*
	列出所有币种
*/
func (this *Tokens) ListTokens() (tokens []Tokens, err error){
	err = Engine_common.Table("tokens").Find(&tokens)
	return
}

//查询所有提币手续费
func (this *Tokens) GetAllTokenFee() ([]Tokens,error) {
	tokens := make([]Tokens,0)
	err := Engine_common.Table("tokens").Find(&tokens)
	return tokens,err
}