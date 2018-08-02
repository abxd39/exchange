package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"time"
	"fmt"
)

type UserSecondaryCertification struct {
	Id                    int    `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	Uid                   int    `xorm:"not null comment('用户uid') INT(64)"`
	VerifyCount           int    `xorm:"not null comment('认证次数') TINYINT(4)"`
	VerifyTime            int  `xorm:"not null comment('认证时间戳') INT(11)"`
	VideoRecordingDigital string `xorm:"not null comment('视频录制的数字') VARCHAR(100)"`
	PositivePath          string `xorm:"not null comment('身份证正面图片路径') VARCHAR(100)"`
	ReverseSidePath       string `xorm:"not null comment('身份证反面图片路径') VARCHAR(100)"`
	InHandPicturePath     string `xorm:"not null comment('身份证反面图片路径') VARCHAR(100)"`
	VideoPath             string `xorm:"not null comment('视频路径') VARCHAR(100)"`
}

func (us *UserSecondaryCertification)GetVerifyCount(uid uint64)(int32,error ) {
	engine := DB.GetMysqlConn()
	_,err:=engine.Where("uid=?",uid).Get(us)
	if err!=nil{
		return 0,err
	}
	return int32(us.VerifyCount),nil
}

//申请二级认证
func (us *UserSecondaryCertification) SetSecondVerify(req *proto.SecondRequest, rsp *proto.SecondResponse) (ret int32, err error) {
	engine := DB.GetMysqlConn()
	u := new(User)
	has, err := engine.Where("uid=?", req.Uid).Get(u)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	if!has{
		return ERRCODE_ACCOUNT_NOTEXIST,nil
	}

	has, err = engine.Where("uid=?", req.Uid).Get(us)
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	sess :=engine.NewSession()
	defer  sess.Close()
	// 启动事务
	if err = sess.Begin(); err != nil {
		return ERRCODE_UNKNOWN,err
	}

	fmt.Println("--------------------->UserSecondaryCertification")
	if !has {
		//写数据库
		if _,err =sess.InsertOne(&UserSecondaryCertification{
			Uid: int(req.Uid),
			VideoRecordingDigital: req.Number,
			PositivePath:          req.FrontPath,
			ReverseSidePath:       req.ReversePath,
			VideoPath:             req.VideoPath,
			InHandPicturePath:req.HeadPath,
			VerifyTime:int( time.Now().Unix()),
		});err!=nil{
			sess.Rollback()
			return ERRCODE_UNKNOWN,err
		}
		u.SetTardeMark= u.SetTardeMark^4
		if _,err=engine.Table("user").Where("uid=?",req.Uid).Update(&User{
			SetTardeMark:u.SetTardeMark,
		});err!=nil{
			sess.Rollback()
			return ERRCODE_UNKNOWN,err
		}
		if err != nil {
			return ERRCODE_UNKNOWN, err
		}
		sess.Commit()
		return 0,nil
	}

	if _,err =sess.Where("uid=?",req.Uid).Update(&UserSecondaryCertification{
		Uid: int(req.Uid),
		VideoRecordingDigital: req.Number,
		PositivePath:          req.FrontPath,
		ReverseSidePath:       req.ReversePath,
		VideoPath:             req.VideoPath,
		InHandPicturePath:req.HeadPath,
		VerifyTime:int( time.Now().Unix()),
	});err!=nil{
		sess.Rollback()
		return ERRCODE_UNKNOWN,err
	}
	u.SetTardeMark= u.SetTardeMark^4
	if _,err=engine.Table("user").Where("uid=?",req.Uid).Update(&User{
		SetTardeMark:u.SetTardeMark,
	});err!=nil{
		sess.Rollback()
		return ERRCODE_UNKNOWN,err
	}
	if err != nil {
		return ERRCODE_UNKNOWN, err
	}
	sess.Commit()
	u.ForceRefreshCache(u.Uid)
	return 0, nil
}
