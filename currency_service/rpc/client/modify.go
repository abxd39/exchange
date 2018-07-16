package client

import (
	"context"
	proto "digicon/proto/rpc"
)

/*
	func: 验证码验证rpc
*/
func (s *UserRPCCli) CallAuthVerify(req *proto.AuthVerifyRequest) (rsp *proto.AuthVerifyResponse, err error) {
	return s.userconn.AuthVerify(context.TODO(), req)
}
