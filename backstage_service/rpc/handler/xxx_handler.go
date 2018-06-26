package handler

import (
	"context"
	def "digicon/proto/common"
	proto "digicon/proto/rpc"
)

type RPCServer struct{}

func (s *RPCServer) Test(ctx context.Context, req *proto.TestRequest, rsp *proto.TestResponse) error {
	rsp.Code = def.ERRCODE_SUCCESS
	rsp.Msg = "Hello BACKSTAGE SERVICE ^v^ !!"
	return nil
}
