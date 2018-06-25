package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
)

type UserEx struct {
	Uid           int    `xorm:"not null pk comment(' 用户ID') INT(11)"`
	NickName      string `xorm:"not null comment('昵称') VARCHAR(64)"`
	HeadSculpture string `xorm:"not null comment('头像图路径') VARCHAR(100)"`
	RegisterTime  int64  `xorm:"comment('注册时间') INT(20)"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId      int    `xorm:"comment('邀请者') INT(20)"`
	Invites       int    `xorm:"comment('邀请人数') INT(20)"`
}

func (s *UserEx) GetUserEx(uid int) (ret int32, err error) {
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

func (ex *UserEx) ModifyNickName(req *proto.UserModifyNickNameResquest, rsp *proto.UserModifyNickNameResponse) (ret int32, err error) {
	//
	engine := DB.GetMysqlConn()
	for _, id := range req.Uid {
		uex := &UserEx{}
		ok, verr := engine.Where("uid=?", id).Get(uex)
		if verr != nil {
			Log.Errorln(err.Error())
			ret = ERRCODE_UNKNOWN
			err = verr
			return
		}

		if ok {
			userEx := &proto.UserModifyNickNameResponse_UserNickName{
				Uid:           uex.Uid,
				NackName:      uex.NickName,
				HeadSculpture: uex.HeadSculpture,
			}

			rsp.User = append(rsp.User, userEx)
		}

	}
	return
}
