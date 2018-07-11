package sync

import (
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	. "digicon/token_service/log"
	"digicon/token_service/model"
	"digicon/token_service/rpc"
)

var Symbols map[int]*proto.TokensData

func LoadSymbol(token_id int32) (err error) {

	//t := new(model.QuenesConfig).GetQuenes(ty)

	d := make([]model.QuenesConfig, 0)
	err = DB.GetMysqlConn().Where("switch=", 1).Find(&d)
	if err != nil {
		Log.Errorln(err.Error())
		return
	}

	ids := make([]int32, 0)
	for _, v := range d {
		ids = append(ids, int32(v.TokenId))
		ids = append(ids, int32(v.TokenTradeId))
	}

	h, err := rpc.InnerService.PublicSevice.CallGetTokensList(ids)
	if err != nil {
		Log.Errorln(err.Error())
		return err
	}

	//Symbols = make(map[int]*proto.TokensData)
	for _, v := range h.Tokens {
		Symbols[int(v.TokenId)] = v
	}
	return nil
}

func init() {
	Symbols = make(map[int]*proto.TokensData)
}
