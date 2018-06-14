package model

type User struct {
	Uid              int    `xorm:"not null pk autoincr INT(11)"`
	Pwd              string `xorm:"VARCHAR(255)"`
	Phone            string `xorm:"unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"INT(11)"`
	GoogleVerifyId   string `xorm:"VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"INT(255)"`
}
