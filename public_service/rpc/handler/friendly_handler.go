package handler

import (
	"context"
	proto "digicon/proto/rpc"
	"digicon/public_service/log"
	"digicon/public_service/model"
)

func (s *RPCServer) AddFriendlyLink(ctx context.Context, req *proto.AddFriendlyLinkRequest, rsp *proto.AddFriendlyLinkResponse) error {
	f := model.FriendlyLink{}
	err := f.Add(req, rsp)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) GetFriendlyLink(ctx context.Context, req *proto.FriendlyLinkRequest, rsp *proto.FriendlyLinkResponse) error {
	f := model.FriendlyLink{}
	err := f.GetFriendlyLinkList(req, rsp)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}
