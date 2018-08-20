package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	log "github.com/sirupsen/logrus"
)

func Test() {
	AddTokenSess(&proto.AddTokenNumRequest{
		Uid:        1,
		TokenId:    1,
		Num:        2000000000,
		Opt:        proto.TOKEN_OPT_TYPE_ADD,
		Ukey:       []byte("test2"),
		Type:       proto.TOKEN_TYPE_OPERATOR_NONE,
		OptAddType: proto.TOKEN_OPT_TYPE_ADD_TYPE_BALANCE,
	})
}

//处理请求代币增加减少事务处理过程
func AddTokenSess(req *proto.AddTokenNumRequest) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.Errorln(err.Error())
		}
	}()

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

	//isAddFrozen := false
	if proto.TOKEN_OPT_TYPE_DEL == req.Opt {
		ret, err = u.SubMoney(session, req.Num, string(req.Ukey), req.Type)
		if err != nil {
			log.Errorln(err.Error())
			session.Rollback()
			return
		}
	} else if proto.TOKEN_OPT_TYPE_ADD == req.Opt {
		if proto.TOKEN_OPT_TYPE_ADD_TYPE_BALANCE == req.OptAddType { //加余额
			err = u.AddMoney(session, req.Num, string(req.Ukey), req.Type)
		} else { //加冻结余额
			//isAddFrozen = true
			err = u.AddFrozen(session, req.Num, string(req.Ukey), req.Type)
		}

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
	/*
		if isAddFrozen {
			err = new(Frozen).InsertRecord(session, &Frozen{
				Uid:     req.Uid,
				Ukey:    string(req.Ukey),
				Num:     req.Num,
				TokenId: int(req.TokenId),
				Type:    int(req.Type),
				Opt:     int(req.Opt),
			})
		} else {
			err = InsertRecord(session, &MoneyRecord{
				Uid:     req.Uid,
				TokenId: int(req.TokenId),
				Ukey:    string(req.Ukey),
				Opt:     int(req.Opt),
				Type:    int(req.Type),
				Balance: u.Balance,
				Num:     req.Num,
			})
		}

		if err != nil {
			log.Errorln(err.Error())
			session.Rollback()
			return
		}
	*/
	err = session.Commit()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return

}

//扣减金额并冻结
func SubTokenWithFronzen(req *proto.SubTokenWithFronzeRequest) (ret int32, err error) {
	u := &UserToken{}
	err = u.GetUserToken(req.Uid, int(req.TokenId))
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	var ok bool
	r := &MoneyRecord{}

	log.WithFields(log.Fields{
		"uid":      req.Uid,
		"token_id": req.TokenId,
		"opt":      req.Opt,
		"num":      req.Num,
		"ukey":     req.Ukey,
		"ukey_str": string(req.Ukey),
		"type":     req.Type,
	}).Errorf("inset  money record error %s", err.Error())
	ok, err = r.CheckExist(string(req.Ukey), req.Type)
	if err != nil {
		return
	}
	if ok {
		ret = ERR_TOKEN_REPEAT
		return
	}
	if u.Balance < req.Num {
		ret = ERR_TOKEN_LESS
		return
	}

	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()
	ret, err = u.SubMoneyWithFronzen(session, req.Num, string(req.Ukey), req.Type)
	if err != nil || ret != ERRCODE_SUCCESS {
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

func ConfirmSubFrozenToken(req *proto.ConfirmSubFrozenRequest) (err error) {
	u := &UserToken{}
	err = u.GetUserToken(req.Uid, int(req.TokenId))
	if err != nil {
		return
	}
	var ok bool
	r := &MoneyRecord{}
	ok, err = r.CheckExist(string(req.Ukey), req.Type)
	if err != nil {
		return
	}
	if ok {
		return
	}

	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()
	err = u.NotifyDelFronzen(session, req.Num, string(req.Ukey), req.Type)
	if err != nil {
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

func CancelFronzeToken(req *proto.CancelFronzeTokenRequest) (err error) {
	u := &UserToken{}
	err = u.GetUserToken(req.Uid, int(req.TokenId))
	if err != nil {
		return nil
	}
	session := DB.GetMysqlConn().NewSession()
	defer session.Close()
	err = session.Begin()
	err = u.ReturnFronzen(session, req.Num, string(req.Ukey), req.Type)
	if err != nil {
		session.Rollback()
		return
	}
	err = session.Commit()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	return nil
}
