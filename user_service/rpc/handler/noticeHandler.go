package handler

import (
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"digicon/user_service/model"

	"golang.org/x/net/context"
)

func (s *Greeter) NotiveList(ctx context.Context, req *proto.NoticeListRequest, rsp *proto.NoticeListResponse) error {
	result := make([]model.NoticeStruct, 0)
	ok := DB.NoticeList(&result)
	if ok {
		len := len(result)
		for _, value := range result {
			ntc := proto.NoticeListResponse_Notice{}
			ntc.Id = value.Id
		}

	}
	return nil
}

func (s *Greeter) NotiveDetail(ctx context.Context, req *proto.NoticeDetailRequest, rsp *proto.NoticeDetailResponse) error {

	return nil
}
