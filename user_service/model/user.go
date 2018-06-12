package model

type User struct {
	Uid              int
	Pwd              string
	Phone            string
	PhoneVerifyTime  int
	Email            string
	EmailVerifyTime  int
	GoogleVerifyTime int
}

type UserEx struct {
	Uid          int
	RegisterTime int64
	InviteCode   string
	RealName     string
	IdentifyCard string
}
