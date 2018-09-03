package handler

import (
	"github.com/apex/log"
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"digicon/user_service/model"
	. "digicon/proto/common"
)

//发送通知
func (s *RPCServer) SendNotice(ctx context.Context, req *proto.SendNoticeRequest, rsp *proto.CommonErrResponse) error {
	var ret int32
	var err error

	log.WithFields(log.Fields{
		"phone":   req.Phone,
		"email":   req.Email,
		"msg": req.Msg,
	}).Info("SendNotice")

	ret, err = model.SendNoticeSms(req.Phone, req.Msg)
	if err != nil {
		rsp.Err = ret
		rsp.Message = err.Error()
		return nil
	}

	//暂不处理邮件通知

	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}
