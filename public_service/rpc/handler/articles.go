package handler

import (
	proto "digicon/proto/rpc"
	. "digicon/public_service/dao"
	"digicon/user_service/model"
	"fmt"

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
	fmt.Println(result)
	//ntc := proto.ArticlesDetailResponse{}
	rsp.Id = result.ID
	rsp.Title = result.Title
	rsp.Description = result.Description
	rsp.Content = result.Content
	rsp.Covers = result.Covers
	rsp.ContentImages = result.ContentImages
	rsp.Type = result.Type
	rsp.TypeName = result.TypeName
	rsp.Author = result.Author
	rsp.Weight = result.Weight
	rsp.Shares = result.Shares
	rsp.Hits = result.Hits
	rsp.Comments = result.Comments
	rsp.DisplayMark = result.DisplayMark
	rsp.CreateTime = result.CreateTime
	rsp.UpdateTime = result.UpdateTime
	rsp.AdminId = result.AdminID
	rsp.AdminNickname = result.AdminNickname
	return nil
}
