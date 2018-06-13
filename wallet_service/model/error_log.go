package models

import (
"time"
)

type ErrorLog struct {
	Id    int       `xorm:"INT(11)"`
	Param string    `xorm:"not null comment('参数') VARCHAR(255)"`
	Msg   string    `xorm:"not null comment('错误提示') VARCHAR(255)"`
	Ctime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
}

