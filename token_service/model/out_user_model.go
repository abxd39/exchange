package model

type OutUser struct {
	Uid int64 `xorm:"uid pk"`
}

func (s *OutUser) TableName() string {
	return "g_common.user"
}
