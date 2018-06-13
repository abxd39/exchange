package handler

import (
	proto "digicon/proto/rpc"
	. "digicon/user_service/dao"

	"golang.org/x/net/context"
)

func (s *RPCServer) NotiveList(ctx context.Context, req *proto.NoticeListRequest, rsp *proto.NoticeListResponse) error {
	//result := make([]model.NoticeStruct, 0)
	result, ok := DB.NoticeList()
	if ok {
		len := len(result)
		for _, value := range result {
			ntc := proto.NoticeListResponse_Notice{}
			ntc.Id = value.ID
		}

	}
	return nil
}

func (s *RPCServer) NotiveDetail(ctx context.Context, req *proto.NoticeDetailRequest, rsp *proto.NoticeDetailResponse) error {

	return nil
}
