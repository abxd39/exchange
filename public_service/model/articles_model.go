package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/public_service/dao"
	//Dlog "digicon/public_service/log"
	Dlog "github.com/sirupsen/logrus"
	"log"
	"digicon/common/model"
	"fmt"
)

type Article struct {
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
	Astatus       int    `xorm:"default 1 comment('1 显示 0 不显示') TINYINT(1)"`
	CreateTime    string `xorm:"default '' comment('创建时间') VARCHAR(36)"`
	UpdateTime    string `xorm:"VARCHAR(36)"`
	AdminId       int    `xorm:"INT(4)"`
	AdminNickname string `xorm:"not null default '' comment('管理员名字') VARCHAR(50)"`
}

type Article_list struct {
	model.ModelList  `xorm:"-"`
	Id          int    `xorm:"not null pk autoincr comment('自增ID') INT(10)"`
	Description string `xorm:"default '' comment('描述') VARCHAR(1000)"`
	Title       string `xorm:"not null default '' comment('文章标题') VARCHAR(100)"`
	Covers        string `xorm:"default '' comment('封面图片') VARCHAR(1000)"`
	CreateTime  string `xorm:"default '' comment('创建时间') VARCHAR(36)"`
}
type ArticleType struct {
	Id       int    `xorm:"not null pk autoincr MEDIUMINT(6)"`
	TypeId   int    `xorm:"not null default 0 TINYINT(10)"`
	TypeName string `xorm:"not null default '' comment('类型名称 1关于我们，2媒体报道，3联系我们，4团队介绍，5数据资产介绍，6服务条款，7免责声明，8隐私保护9 业界新闻 10 公告 11 帮助手册 12 币种介绍') VARCHAR(100)"`
}

func (a *Article_list) TableName() string {
	return "article"
}

func (ArticleType) GetArticleTypeList() ([]ArticleType, error) {
	engine := dao.DB.GetMysqlConn()
	list := make([]ArticleType, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a Article_list) ArticleList(req *proto.ArticleListRequest, u *[]Article_list) (*model.ModelList, int32) {
	//err := s.mysql.im.Find(&u)
	//default_page_count := int(10)

	engine := dao.DB.GetMysqlConn()

	//没有指定 每页的行数
	count, err := engine.Where("type =? and astatus=1", req.ArticleType).Count(&Article_list{})
	if err != nil {
		Dlog.Errorln(err.Error())
		return nil, ERRCODE_UNKNOWN
	}
	offset,mList:= model.Paging(int(req.Page),int(req.PageNum),int(count))

	fmt.Println("total=",mList.Total, "type=", req.ArticleType, "page=", mList.PageSize, "起始行star_row=", mList.PageIndex, "page_num=", offset)
	err = engine.Desc("id").AllCols().Where("astatus=1 and type=?", req.ArticleType).Limit(mList.PageSize,offset).Find(u)
	if err != nil {
		log.Fatalf(err.Error())
	}


	return mList, ERRCODE_SUCCESS

}

func (Article) Article(Id int32, u *Article) int32 {
	engine := dao.DB.GetMysqlConn()
	ok, err := engine.Where("id=?", Id).Get(u)
	if err != nil {
		Dlog.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	if ok {
		return ERRCODE_SUCCESS

	}

	return ERRCODE_ARTICLE_NOT_EXIST

}
