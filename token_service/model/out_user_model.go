package model

type OutUser struct {
	Uid    int64 `xorm:"uid pk"`
	Status int   `xorm:"status"`
}

func (s *OutUser) TableName() string {
	return "g_common.user"
}
