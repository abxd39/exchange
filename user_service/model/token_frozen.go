package model

type TokenFrozen struct {
	Uid        uint64 `xorm:"uid"`
	Ukey       string `xorm:"ukey"`
	Num        int64  `xorm:"num"`
	TokenId    int    `xorm:"token_id"`
	Type       int    `xorm:"type"`
	CreateTime int64  `xorm:"create_time"`
	Opt        int    `xorm:"opt"`
}

func (*TokenFrozen) TableName() string {
	return "g_token.frozen"
}
