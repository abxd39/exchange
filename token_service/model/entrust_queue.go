package model

import (
	. "digicon/user_service/dao"
)
type EntrustQuene struct{
	
}

type OrderDetail struct {
	Uid int32
} 

//限价委托入队列
func (s *EntrustQuene) JoinSellQuene(token_quene_id string,uid int32,)  {

	DB.GetRedisConn().ZAdd(token_quene_id,)
}



