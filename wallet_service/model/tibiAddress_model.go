package models

import (
	"digicon/wallet_service/utils"
	proto "digicon/proto/rpc"
	"fmt"
)

type TibiAddress struct {
	Id      int    `xorm:"not null pk autoincr INT(11)"`
	Uid     int    `xorm:"not null comment('用户id') INT(11)"`
	TokenId int    `xorm:"not null comment('币种id') INT(11)"`
	Address string `xorm:"not null comment('地址') VARCHAR(60)"`
	Mark    string `xorm:"not null default '' comment('备注') VARCHAR(255)"`
}


func (this *TibiAddress)Save(uid int,tokenid int,address string,mark string)(int,error){
	this.Uid= uid
	this.TokenId=tokenid
	this.Address=address
	this.Mark=mark
	affected,err :=utils.Engine_wallet.Insert(this)
	return int(affected),err

}

func (this *TibiAddress)List(uid int,tokenid int)(lists []*proto.AddrlistPos,err error){
	this.Id= uid
	this.TokenId=tokenid

	//rets,err:=utils.Engine_wallet.Query("select * from tibi_address where uid=? and token_id=?",uid,tokenid)
	rets := make([]TibiAddress, 0)
	//lists:= make([]proto.AddrlistPos, 0)

	err =utils.Engine_wallet.Where("uid=? and token_id=?",uid,tokenid).Find(&rets)
	fmt.Println(rets)
	if err !=nil {
		return nil,err

	}
	for i:=0;i<len(rets);i++{
		temp := &proto.AddrlistPos{
			Id: int32(rets[i].Id),
			Uid:int32(rets[i].Uid),
			TokenId:int32(rets[i].TokenId),
			Address:rets[i].Address,
			Mark:rets[i].Mark,
		}
		lists = append(lists,temp)
	}
	return lists,err

}


func (this *TibiAddress)DeleteByid(id int,uid int)(int,error){
	utils.Engine_wallet.ShowSQL(true)
	affected,err:=utils.Engine_wallet.Where("id=? and uid=?",id,uid).Delete(this)
	return int(affected),err
}