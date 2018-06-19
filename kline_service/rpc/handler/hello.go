package handler

import (
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"log"
)

type RPCServer struct{}

func (s *RPCServer) Hline(ctx context.Context, req *proto.KineRequest, rsp *proto.KlineResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Data = "Hello " + req.Cmd
	return nil
}
