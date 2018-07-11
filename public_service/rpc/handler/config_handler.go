package handler

import (
proto "digicon/proto/rpc"
"digicon/public_service/model"
"golang.org/x/net/context"
)

func (s *RPCServer) GetTokensList(ctx context.Context, req *proto.TokensRequest, rsp *proto.TokensResponse) error {
	u:=model.Tokens{}
	r:=u.GetTokens(req.Tokens)
	for _,v:=range r {

		rsp.Tokens=append(rsp.Tokens,&proto.TokensData{
			TokenId:int32(v.Id),
			Mark:v.Mark,
		})
	}
	return nil
}
