package model

import (
	proto "digicon/proto/rpc"
	"digicon/public_service/dao"
	"fmt"
)

type Banner struct {
	Id          int    `xorm:"not null pk autoincr INT(11)"`
	Order       int    `xorm:"not null default 1 comment('排序') TINYINT(4)"`
	PictureName string `xorm:"not null default '' comment('图片名称') VARCHAR(255)"`
	UploadTime  string `xorm:"not null comment('上传时间') DATETIME"`
	LinkPath    string `xorm:"not null default '' comment('链接地址') VARCHAR(255)"`
	PicturePath string `xorm:"not null default '' comment('图片路径') VARCHAR(255)"`
	Status      int    `xorm:"not null default 1 comment('上架状态 1 上架 2下架') TINYINT(4)"`
}

func (b *Banner) GetBannerList(req *proto.BannerRequest, rsp *proto.BannerResponse) error {
	fmt.Println("xxx")
	engine := dao.DB.GetMysqlConn()
	ban := make([]Banner, 0)
	fmt.Println("000000000000000")
	err := engine.Desc("order").Where("status=1").Find(&ban)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("11111111111111111")
	for _, v := range ban {
		banner := &proto.BannerResponse_List{
			Order:       int32(v.Order),
			PictureName: v.PictureName,
			TimeStart:   v.UploadTime,
			TimeEnd:     v.UploadTime,
			LinkPath:    v.LinkPath,
			PicturePath: v.PicturePath,
		}
		rsp.List = append(rsp.List, banner)

	}
	fmt.Printf("22222222222222222222222%#v\n", rsp.List)
	return nil
}
