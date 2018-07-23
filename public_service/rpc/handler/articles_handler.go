package handler

import (
	proto "digicon/proto/rpc"
	"digicon/public_service/model"
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/context"
)

func (s *RPCServer) ArticleList(ctx context.Context, req *proto.ArticleListRequest, rsp *proto.ArticleListResponse) error {
	var total int
	result := make([]model.Article_list, 0)

	total, rsp.Code = new(model.Article_list).ArticleList(req, &result)
	rsp.TotalPage = int32(total)
	for _, value := range result {
		rsp.Article = append(rsp.Article, &proto.ArticleListResponse_Article{
		Id : int32(value.Id),
		Title :value.Title,
		Description:value.Description,
		CreateDateTime:value.CreateTime,
		})

	}
	//fmt.Println("ArticleList 列表为", ntc)
	return nil
}

func (s *RPCServer) Article(ctx context.Context, req *proto.ArticleRequest, rsp *proto.ArticleResponse) error {
	result := &model.Article{}
	//result := new(model.Article)
	rsp.Code = new(model.Article).Article(req.Id, result)
	fmt.Println(result)
	js, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("struct model.Article Marshal Fatalf!!")
		return err
	}
	//json.Unmarshal
	rsp.Data = string(js)
	//fmt.Println(rsp.Data)
	return nil
}

func (s *RPCServer) ArticleTypeList(ctx context.Context, req *proto.ArticleTypeRequest, rsp *proto.ArticleTypeListResponse) error {
	list, err := new(model.ArticleType).GetArticleTypeList()
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	for _, v := range list {
		at := &proto.ArticleTypeListResponse_ArticleType{
			Id:   int32(v.TypeId),
			Name: v.TypeName,
		}
		rsp.Type = append(rsp.Type, at)
	}
	return nil
}
