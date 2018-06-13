package model

type Articles struct {
	Id            int    `xorm:"not null pk autoincr comment('自增ID') INT(10)"`
	Title         string `xorm:"not null default '' comment('文章标题') VARCHAR(100)"`
	Description   string `xorm:"default '' comment('描述') VARCHAR(1000)"`
	Content       string `xorm:"not null comment('内容') TEXT"`
	Covers        string `xorm:"default '' comment('封面图片') VARCHAR(1000)"`
	ContentImages string `xorm:"comment('内容图片') TEXT"`
	Type          int    `xorm:"default 1 comment('类型 1 业界新闻 2 公告 3 帮助手册') TINYINT(4)"`
	TypeName      string `xorm:"default '' comment('类型名字') VARCHAR(50)"`
	Author        string `xorm:"default '' comment('作者名字') VARCHAR(150)"`
	Weight        int    `xorm:"not null default 0 comment('权重，排序字段') TINYINT(4)"`
	Shares        int    `xorm:"default 0 comment('分享数量') INT(11)"`
	Hits          int    `xorm:"default 0 comment('点击数量') INT(11)"`
	Comments      int    `xorm:"default 0 comment('评论数量') INT(11)"`
	DisplayMark   int    `xorm:"default 1 comment('1 显示 0 不显示') TINYINT(1)"`
	AdminId       int    `xorm:"INT(4)"`
	CreateTime    string `xorm:"default '' comment('创建时间') VARCHAR(36)"`
	UpdateTime    string `xorm:"VARCHAR(36)"`
	AdminNickname string `xorm:"not null default '' comment('管理员名字') VARCHAR(50)"`
}
