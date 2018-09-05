package models

import (
	proto "digicon/proto/rpc"
	"digicon/wallet_service/utils"
	"fmt"
)

type TibiAddress struct {
	Id      int    `xorm:"not null pk autoincr INT(11)"`
	Uid     int    `xorm:"not null comment('用户id') INT(11)"`
	TokenId int    `xorm:"not null comment('币种id') INT(11)"`
	Address string `xorm:"not null comment('地址') VARCHAR(60)"`
	Mark    string `xorm:"not null default '' comment('备注') VARCHAR(255)"`
}

func (this *TibiAddress) Save(uid int, tokenid int, address string, mark string) (int, error) {
	this.Uid = uid
	this.TokenId = tokenid
	this.Address = address
	this.Mark = mark
	affected, err := utils.Engine_wallet.Insert(this)
	return int(affected), err
}

func (this *TibiAddress) List(uid int) (lists []*proto.AddrlistPos, err error) {
	this.Id = uid
	rets := make([]TibiAddress, 0)

	err = utils.Engine_wallet.Where("uid=?", uid).Desc("created_time").Find(&rets)

	if err != nil {
		return nil, err
	}
	retsLen := len(rets)
	tokenIdsList := make([]int, 0, retsLen)
	for i := 0; i < len(rets); i++ {
		temp := &proto.AddrlistPos{
			Id:        int32(rets[i].Id),
			Uid:       int32(rets[i].Uid),
			TokenId:   int32(rets[i].TokenId),
			Address:   rets[i].Address,
			Mark:      rets[i].Mark,
			TokenName: "",
		}
		tokenIdsList = append(tokenIdsList, rets[i].TokenId)
		lists = append(lists, temp)
	}
	tks := []Tokens{}
	//fmt.Println("idsList:", tokenIdsList)
	err = utils.Engine_common.In("id", tokenIdsList).Find(&tks)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for i := 0; i < retsLen; i++ {
			for _, ti := range tks {
				if lists[i].TokenId == int32(ti.Id) {
					lists[i].TokenName = ti.Name
					break
				}
			}
		}
	}
	fmt.Println(lists)
	return lists, err

}

func (this *TibiAddress) DeleteByid(id int, uid int) (int, error) {
	utils.Engine_wallet.ShowSQL(true)
	affected, err := utils.Engine_wallet.Where("id=? and uid=?", id, uid).Delete(this)
	return int(affected), err
}

func (this *TibiAddress) GetByAddress(address string) (bool,error) {
	return utils.Engine_wallet.Where("address = ?",address).Get(this)
}

//判断提币地址是否存在
func (this *TibiAddress) TiBiAddressExists(address string) (bool,error) {
	return utils.Engine_wallet.Where("address = ?",address).Exist(this)
}
