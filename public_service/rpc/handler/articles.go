package handler

import (
	proto "digicon/proto/rpc"
	. "digicon/public_service/dao"
	"digicon/user_service/model"

	"golang.org/x/net/context"
)

func (s *RPCServer) ArticlesList(ctx context.Context, req *proto.ArticlesListRequest, rsp *proto.ArticlesListResponse) error {
	result := make([]model.ArticlesStruct, 0)
	rsp.Err = DB.ArticlesList(req.ArticlesType, req.StartRow, req.EndRow, &result)
	//len := len(result)
	ntc := proto.ArticlesListResponse_Articles{}
	for _, value := range result {
		ntc.Id = value.ID
		ntc.Title = value.Title
		ntc.Description = value.Description
		ntc.CreateDateTime = value.CreateDateTime
		rsp.Articles = append(rsp.Articles, &ntc)
	}

	return nil
}

func (s *RPCServer) ArticlesDetail(ctx context.Context, req *proto.ArticlesDetailRequest, rsp *proto.ArticlesDetailResponse) error {
	result := model.ArticlesDetailStruct{}
	rsp.Err = DB.ArticlesDescription(req.Id, &result)
	ntc := proto.ArticlesDetailResponse{}
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
	return nil
}
