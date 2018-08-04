package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"fmt"
	log "github.com/sirupsen/logrus"
	."digicon/common/constant"
	"time"
)

type UserEx struct {
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

func (s *UserEx) TableName() string {
	return "user_ex"
}

func (s *UserEx) GetUserEx(uid uint64) (ret int32, err error) {
	ok, err := DB.GetMysqlConn().Where("uid=?", uid).Get(s)
	if err != nil {
		log.Errorln(err.Error())
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
	sql := "SELECT user.`uid`, account, user_ex.`nick_name`, user_ex.`head_sculpture` FROM g_common.`user` LEFT JOIN   user_ex ON user.`uid` = user_ex.`uid`"
	type UNickName struct {
		Uid        uint64     `json:"uid"`
		Account    string      `json:"account"`
		NickName   string      `json:"nick_name"`
		HeadSculpture  string   `json:"head_sculpture"`
	}
	var uex []UNickName
	err = engine.SQL(sql).In("uid", req.Uid).Find(&uex)
	//fmt.Println("uid:", req.Uid)
	//uex := make([]UserEx, 0)
	//err = engine.In("uid", req.Uid).Find(&uex)
	if err != nil {
		log.Errorln(err.Error())
		ret = ERRCODE_UNKNOWN
		return
	}
	for _, value := range uex {
		var nickname string
		if value.NickName == ""{
			nickname = value.Account
		}else{
			nickname = value.NickName
		}
		userEx := &proto.UserGetNickNameResponse_UserNickName{
			Uid:           uint64(value.Uid),
			NickName:      nickname,
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
		Uid: int64(req.Uid),
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
		Uid:           int64(req.Uid),
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

//申请一级实名认证
func (ex *UserEx) SetFirstVerify(req *proto.FirstVerifyRequest, rsp *proto.FirstVerifyResponse) (ret int32, err error) {
	engine := DB.GetMysqlConn()
	fmt.Println("---------------->258369")
	u := new(User)
	has, err := engine.Table("user").Where("uid=?", req.Uid).Get(u)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !has {
		ret = ERRCODE_ACCOUNT_NOTEXIST
		return
	}
	//电话验证
	result, err := AuthSms(u.Phone, SMS_REAL_NAME, req.PhoneCode)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if result != ERRCODE_SUCCESS {
		return result, nil
	}
	//google 认证
	if req.GoogleCode != 0 {
		_, err := u.AuthGoogleCode(u.GoogleVerifyId, req.GoogleCode)
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
	}
	has, err = engine.Where("uid=?", req.Uid).Get(ex)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if !has {
		return ERRCODE_ACCOUNT_NOTEXIST, nil
	}
	sess := engine.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return ERRCODE_UNKNOWN, err
	}
	//写数据库 如果 user 表上有 该用户，则 表user_ex 此表一定有该 用户
	if _, err = sess.Where("uid=?", req.Uid).Cols("real_name","identify_card","affirm_time","affirm_count").Update(&UserEx{
		RealName:     req.RealName,
		IdentifyCard: req.IdCode,
		AffirmTime:   time.Now().Unix(),
		AffirmCount:  ex.AffirmCount + 1,
	}); err != nil {
		sess.Rollback()

		return ERRCODE_UNKNOWN, err
	}
	u.SetTardeMark = u.SetTardeMark ^ APPLY_FOR_FIRST
	if _, err = sess.Table("user").Where("uid=?", req.Uid).Cols("set_tarde_mark").Update(&User{
		SetTardeMark: u.SetTardeMark,
	}); err != nil {
		sess.Rollback()
		return ERRCODE_UNKNOWN, err
	}
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	sess.Commit()

	u.ForceRefreshCache(u.Uid)
	return 0, nil
}

//获取认证次数 实名 和二级认证次数

func (ex *UserEx) GetVerifyCount(req *proto.VerifyCountRequest, rsp *proto.VerifyCountResponse) (ret int32, err error) {
	engine := DB.GetMysqlConn()
	_, err = engine.Where("uid=?", req.Uid).Get(ex)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	rsp.SecondCount, err = new(UserSecondaryCertification).GetVerifyCount(req.Uid)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	rsp.FirstCount = int32(ex.AffirmCount)
	return ERRCODE_SUCCESS, nil
}
