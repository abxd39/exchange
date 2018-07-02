package model

import (
	proto "digicon/proto/rpc"
	"digicon/public_service/dao"
	"fmt"
)

type Banner struct {
	Id          int    `xorm:"not null pk INT(11)"`
	Order       int    `xorm:"not null default 1 comment('排序') TINYINT(4)"`
	PictureName string `xorm:"not null default '' comment('图片名称') VARCHAR(255)"`
	TimeStart   string `xorm:"not null comment('展示开始日期') DATETIME"`
	TimeEnd     string `xorm:"not null comment('展示结束日期') DATETIME"`
	LinkPath    string `xorm:"not null default '' comment('链接地址') VARCHAR(255)"`
	PicturePath string `xorm:"not null default '' comment('图片路径') VARCHAR(255)"`
	State       int    `xorm:"not null default 1 comment('上架状态 1 上架 0下架') TINYINT(4)"`
}

func (b *Banner) GetBannerList(req *proto.BankPayRequest, rsp *proto.BannerResponse) error {
	fmt.Println("xxx")
	engine := dao.DB.GetMysqlConn()
	ban := make([]*Banner, 0)
	err := engine.Find(ban)
	if err != nil {
		return err
	}
	for _, v := range ban {
		banner := proto.BannerResponse_List{
			Order:       int32(v.Order),
			PictureName: v.PictureName,
			TimeStart:   v.TimeStart,
			TimeEnd:     v.TimeEnd,
			LinkPath:    v.LinkPath,
			PicturePath: v.PicturePath,
		}
		rsp.List = append(rsp.List, &banner)

	}
	return nil
}
