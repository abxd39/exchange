package models

import (
	"digicon/wallet_service/utils"
)

type UidAndInviteID struct {
	Uid      uint64
	InviteId uint64
}

func (*User) TableName() string {
	return "user"
}


func (s *User) GetUser(uid uint64) (ok bool, err error) {
	ok, err = utils.Engine_common.Where("uid=?", uid).Get(s)
	return
}