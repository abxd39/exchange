package handler

import (
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
)

func (s *RPCServer) NoticeList(ctx context.Context, req *proto.NoticeListRequest, rsp *proto.NoticeListResponse) error {
<<<<<<< HEAD
	result := make([]model.NoticeStruct, 0)
	rsp.Err = DB.NoticeList(req.NoticeType, req.StartRow, req.EndRow, &result)
	//len := len(result)
	ntc := proto.NoticeListResponse_Notice{}
	for _, value := range result {
		ntc.Id = value.ID
		ntc.Title = value.Title
		ntc.Description = value.Description
		ntc.CreateDateTime = value.CreateDateTime
		rsp.Notice = append(rsp.Notice, &ntc)
	}

=======
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
>>>>>>> 22bc5f02800c57be8a373bbe6d6be3b14bc01bed
	return nil
}

func (s *RPCServer) NoticeDetail(ctx context.Context, req *proto.NoticeDetailRequest, rsp *proto.NoticeDetailResponse) error {
<<<<<<< HEAD
	result := model.NoticeDetailStruct{}
	rsp.Err = DB.NoticeDescription(req.Id, &result)
	ntc := proto.NoticeDetailResponse{}
	ntc.Id = result.ID
	ntc.Title = result.Title
	ntc.Description = result.Description
	ntc.Content = result.Content
	ntc.Covers = result.Covers
	ntc.ContentImages = result.ContentImages
	ntc.Type = result.Type
	ntc.TypeName = result.TypeName
	ntc.Author = result.Author
	ntc.Weight = result.Weight
	ntc.Shares = result.Shares
	ntc.Hits = result.Hits
	ntc.Comments = result.Comments
	ntc.DisplayMark = result.DisplayMark
	ntc.CreateTime = result.CreateTime
	ntc.UpdateTime = result.UpdateTime
	ntc.AdminId = result.AdminID
	ntc.AdminNickname = result.AdminNickname
=======

>>>>>>> 22bc5f02800c57be8a373bbe6d6be3b14bc01bed
	return nil
}
