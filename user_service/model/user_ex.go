package model

type UserEx struct {
	Uid          int    `xorm:"not null pk comment(' 用户ID') INT(11)"`
	PayPwd       string `xorm:"comment('支付密码') VARCHAR(255)"`
	RegisterTime int64  `xorm:"comment('注册时间') INT(20)"`
	InviteCode   string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName     string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard string `xorm:"comment('身份证号') VARCHAR(64)"`
}
