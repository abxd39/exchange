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

type NoticeStruct struct {
	ID             int32
	Description    string //重要 、一般
	Title          string
	CreateDateTime string
}

type NoticeDetailStruct struct {
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
