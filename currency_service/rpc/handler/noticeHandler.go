package handler

import (
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"
	"digicon/user_service/model"

	"golang.org/x/net/context"
)

type Greeter struct{}

func (s *Greeter) NotiveList(ctx context.Context, req *proto.NoticeListRequest, rsp *proto.NoticeListResponse) error {
	result := make([]model.NoticeStruct, 0)
	ok := DB.NoticeList(&result)
	if ok {
		for _, value := range result {
			ntc := proto.NoticeListResponse_Notice{}
			ntc.Id = value.Id
			ntc.Title = value.Title
			ntc.Description = value.Description
			ntc.CreateDateTime = value.CreateDateTime
			rsp.Notice = append(rsp.Notice, &ntc)
		}
	}
	return nil
}

func (s *Greeter) NotiveDetail(ctx context.Context, req *proto.NoticeDetailRequest, rsp *proto.NoticeDetailResponse) error {

	return nil
}
