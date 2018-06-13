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

type ArticlesStruct struct {
	ID             int32
	Description    string //重要 、一般
	Title          string
	CreateDateTime string
}

type ArticlesDetailStruct struct {
	ID            int32
	Title         string
	Description   string
	Content       string
	Covers        []byte
	ContentImages []byte
	Type          int32
	TypeName      string
	Author        string
	Weight        int32
	Shares        int32
	Hits          int32
	Comments      int32
	DisplayMark   bool
	CreateTime    string
	UpdateTime    string
	AdminID       int32
	AdminNickname string
}

