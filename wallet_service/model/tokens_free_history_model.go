package models

import (
	. "digicon/wallet_service/utils"
)

type Token_free_history struct {
	Id            int    `xorm:"not null pk autoincr INT(11)"`
	Token_id      int64  `xorm:"default '' comment('代币ID') BIGINT(20)"`
	Opt           int64  `xorm:"default '' comment('操作方向1加2减') BIGINT(20)"`
	Type          int64  `xorm:"default '' comment('流水类型') BIGINT(20)"`
	Num           int64  `xorm:"default 0 comment('数量') BIGINT(20)"`
	Created_time   int64  `xorm:"default '' comment('操作时间') BIGINT(20)"`
	Ukey          string `xorm:"default '' VARCHAR(255)"`
	Uid          int64 `xorm:"default '' BIGINT(20)"`
}

//写入数据
func (this *Token_free_history) InsertThis() (int,error) {
	affected,err := Engine_token.InsertOne(this)
	return int(affected),err
}