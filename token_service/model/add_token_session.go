package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	log "github.com/sirupsen/logrus"
)

//处理请求代币增加减少事务处理过程
func AddTokenSess(req *proto.AddTokenNumRequest) (ret int32, err error) {
	u := &UserToken{}
	err = u.GetUserToken(req.Uid, int(req.TokenId))
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	var ok bool
	r := &MoneyRecord{}
	ok, err = r.CheckExist(string(req.Ukey), req.Type)
	if err != nil {
		return
	}
	if ok {
		ret = ERR_TOKEN_REPEAT
		return
	}

	if req.Opt == proto.TOKEN_OPT_TYPE_DEL {
		if u.Balance < req.Num {
			ret = ERR_TOKEN_LESS
			return
		}
	}

	//开始入账
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()

	if proto.TOKEN_OPT_TYPE_DEL == req.Opt {
		ret, err = u.SubMoney(session, req.Num)
		if err != nil {
			log.Errorln(err.Error())
			session.Rollback()
			return
		}
	} else if proto.TOKEN_OPT_TYPE_ADD == req.Opt {
		err = u.AddMoney(session, req.Num)
		if err != nil {
			log.Errorln(err.Error())
			session.Rollback()
			return
		}
	} else {
		ret = ERRCODE_PARAM
		return
	}

	if ret != ERRCODE_SUCCESS {
		session.Rollback()
		return
	}

	err = new(MoneyRecord).InsertRecord(session, &MoneyRecord{
		Uid:     req.Uid,
		TokenId: int(req.TokenId),
		Ukey:    string(req.Ukey),
		Opt:     int(req.Opt),
		Type:    int(req.Type),
		Balance: u.Balance,
		Num:     req.Num,
	})

	if err != nil {
		log.Errorln(err.Error())
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return

}
