package client

import (
	"context"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
)

func (s *UserRPCCli) CallNoticeDesc(id int32) (rsp *proto.NoticeDetailResponse, err error) {
	rsp, err = s.conn.NoticeDetail(context.TODO(), &proto.NoticeDetailRequest{
		Id: id,
	})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (s *UserRPCCli) CallnoticeList() (rsp *proto.NoticeListResponse, err error) {
	rsp, err = s.conn.NoticeList(context.TODO(), &proto.NoticeListRequest{})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
