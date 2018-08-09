package model

import (
	Err "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/public_service/dao"
	//"digicon/public_service/log"
	log "github.com/sirupsen/logrus"
)

type FriendlyLink struct {
	Id        int32  `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	Aorder    int32  `xorm:"not null comment('排序') INT(10)"`
	WebName   string `xorm:"not null default '' comment('网址名称') VARCHAR(100)"`
	LinkName  string `xorm:"not null default '' comment('网站链接') VARCHAR(100)"`
	LinkState int32  `xorm:"not null comment('1,上架2，下架') INT(10)"`
}

func (f *FriendlyLink) Add(req *proto.AddFriendlyLinkRequest, rsp *proto.AddFriendlyLinkResponse) error {
	engine := dao.DB.GetMysqlConn()
	flink := &FriendlyLink{
		Aorder:    req.Aorder,
		WebName:   req.WebName,
		LinkName:  req.LinkName,
		LinkState: req.LinkState,
	}
	result, err := engine.Insert(flink)
	if err != nil {
		log.Errorln(err.Error())
		rsp.Code = Err.ERRCODE_UNKNOWN
		return nil
	}
	if result == 0 {
		log.Errorf("friendly link insert fail!!")
		rsp.Code = Err.ERRCODE_UNKNOWN
		return nil
	}
	rsp.Code = Err.ERRCODE_SUCCESS
	return nil
}

func (f *FriendlyLink) GetFriendlyLinkList(req *proto.FriendlyLinkRequest, rsp *proto.FriendlyLinkResponse) error {
	//page !=0
	if 0 == req.Page {
		req.Page = 1
	}

	defa := req.Count
	if 0 == defa {
		req.Count = 100
	}
	engine := dao.DB.GetMysqlConn()
	u := &FriendlyLink{}
	total, err := engine.Count(u)
	if err != nil {
		log.Errorln("统计所有记录失败")
		rsp.Code = Err.ERRCODE_UNKNOWN
		return nil
	}

	rsp.Page = int32(total) / req.Count
	var limit int32
	if 1 == req.Page {
		limit = 1
	} else {
		limit = req.Page * req.Count
	}

	friendlist := make([]FriendlyLink, 0)
	err = engine.Desc("id").Limit(int(req.Count), int(limit)).Find(&friendlist)
	if err != nil {
		log.Errorln(err.Error())
		rsp.Code = Err.ERRCODE_UNKNOWN
		return nil
	}

	for _, frd := range friendlist {
		rsp.Friend = append(rsp.Friend, &proto.FriendlyLinkResponseFriendlylink{
			Id:        frd.Id,
			Aorder:    frd.Aorder,
			WebName:   frd.WebName,
			LinkName:  frd.LinkName,
			LinkState: frd.LinkState,
		})
	}
	return nil
}
