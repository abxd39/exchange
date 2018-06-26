package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
)

type UserEx struct {
	Uid           uint64 `xorm:"not null pk comment(' 用户ID') BIGINT(11)"`
	NickName      string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	HeadSculpture string `xorm:"not null default '' comment('头像图片路径') VARCHAR(100)"`
	RegisterTime  int64  `xorm:"comment('注册时间') BIGINT(20)"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId      uint64 `xorm:"comment('邀请者id') BIGINT(11)"`
	Invites       int    `xorm:"default 0 comment('邀请人数') INT(11)"`
}

func (s *UserEx) GetUserEx(uid uint64) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().Where("uid=?", uid).Get(s)
	if err != nil {
		Log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	if ok {
		ret = ERRCODE_SUCCESS
		return
	}
	ret = ERRCODE_ACCOUNT_NOTEXIST
	return
}

func (ex *UserEx) GetNickName(req *proto.UserGetNickNameRequest, rsp *proto.UserGetNickNameResponse) (ret int32, err error) {
	//
	engine := DB.GetMysqlConn()

	uex := make([]UserEx, 0)
	err = engine.In("uid", req.Uid).Find(&uex)
	if err != nil {
		Log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}

	for _, value := range uex {
		userEx := &proto.UserGetNickNameResponse_UserNickName{
			Uid:           value.Uid,
			NickName:      value.NickName,
			HeadSculpture: value.HeadSculpture,
		}

		rsp.User = append(rsp.User, userEx)
	}

	ret = ERRCODE_SUCCESS
	return
}

func (ex *UserEx) SetNickName(req *proto.UserSetNickNameRequest, rsp *proto.UserSetNickNameResponse) (ret int32, err error) {
	//检查是否存在该uid
	engine := DB.GetMysqlConn()
	has, err := engine.Exist(&UserEx{
		Uid: req.Uid,
	})
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !has {
		ret = ERRCODE_ACCOUNT_NOTEXIST
		return
	}
	//验证token
	var result int64
	result, err = engine.ID(req.Uid).Update(&UserEx{
		Uid:           req.Uid,
		NickName:      req.NickName,
		HeadSculpture: req.HeadSculpture,
	})
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	if result != 0 {
		ret = ERRCODE_SUCCESS
	}
	return

}
