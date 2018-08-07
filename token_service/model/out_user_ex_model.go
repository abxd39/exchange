package model

import (
	"digicon/common/errors"
	"digicon/token_service/dao"
)

type OutUserEx struct {
	Uid           int64  `xorm:"not null pk comment(' 用户ID') BIGINT(11)"`
	NickName      string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	HeadSculpture string `xorm:"not null default '' comment('头像图片路径') VARCHAR(100)"`
	RegisterTime  int64  `xorm:"comment('注册时间') BIGINT(20)"`
	AffirmTime    int64  `xorm:"comment('实名认证时间') BIGINT(20)"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId      int64  `xorm:"comment('邀请者id') BIGINT(11)"`
	Invites       int    `xorm:"default 0 comment('邀请人数') INT(11)"`
	AffirmCount   int    `xorm:"default 0 comment('实名认证的次数') TINYINT(4)"`
	ChannelName   string `xorm:"not null default '' comment('邀请的渠道名称') VARCHAR(100)"`
}

func (s *OutUserEx) TableName() string {
	return "g_common.user_ex"
}

func (s *OutUserEx) Get(uid int64) (*OutUserEx, error) {
	userEx := &OutUserEx{}
	has, err := dao.DB.GetCommonMysqlConn().Where("uid=?", uid).Get(userEx)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	if !has {
		return nil, errors.NewNormal("用户不存在")
	}

	return userEx, nil
}
