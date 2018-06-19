package handler

import (
	proto "digicon/proto/rpc"
	"digicon/public_service/model"
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/context"
)

func (s *RPCServer) ArticlesList(ctx context.Context, req *proto.ArticlesListRequest, rsp *proto.ArticlesListResponse) error {
	var total int
	result := make([]model.Articles_list, 0)

	total, rsp.Err = new(model.Articles_list).ArticlesList(req, &result)
	rsp.TotalPage = int32(total)
	for _, value := range result {
		ntc := proto.ArticlesListResponse_Articles{}
		ntc.Id = int32(value.Id)
		ntc.Title = value.Title
		ntc.Description = value.Description
		ntc.CreateDateTime = value.CreateTime
		rsp.Articles = append(rsp.Articles, &ntc)

	}
	//fmt.Println("ArticlesList 列表为", ntc)
	return nil
}

func (s *RPCServer) ArticlesDetail(ctx context.Context, req *proto.ArticlesDetailRequest, rsp *proto.ArticlesDetailResponse) error {
	result := &model.ArticlesCopy1{}
	rsp.Err = new(model.ArticlesCopy1).ArticlesDescription(req.Id, result)
	fmt.Println(result)
	js, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("struct model.Articlescopy Marshal Fatalf!!")
	}
	//json.Unmarshal
	rsp.Data = string(js)
	//fmt.Println(rsp.Data)
	return nil
}
