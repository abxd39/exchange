package model

import (
	. "digicon/proto/common"
	. "digicon/user_service/dao"
	. "digicon/user_service/log"
)

type UserEx struct {
	Uid int `xorm:"not null pk comment(' 用户ID') INT(11)"`

	RegisterTime int64  `xorm:"comment('注册时间') INT(20)"`
	InviteCode   string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName     string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId     int    `xorm:"comment('邀请者') INT(20)"`
	Invites      int    `xorm:"comment('邀请人数') INT(20)"`

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
