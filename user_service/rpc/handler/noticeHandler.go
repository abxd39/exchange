package handler

import (
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
)

func (s *RPCServer) NoticeList(ctx context.Context, req *proto.NoticeListRequest, rsp *proto.NoticeListResponse) error {
	/*
		result := make([]model.NoticeStruct, 0)
		ok := DB.NoticeList(&result)
		if ok {
			len := len(result)
			for _, value := range result {
				ntc := proto.NoticeListResponse_Notice{}
				ntc.Id = value.Id
			}

		}
	*/
	return nil
}

func (s *RPCServer) NoticeDetail(ctx context.Context, req *proto.NoticeDetailRequest, rsp *proto.NoticeDetailResponse) error {

	return nil
}
