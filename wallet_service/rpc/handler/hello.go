package handler

import (
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"log"
)

type Greeter struct{}

func (s *Greeter) Hello(ctx context.Context, req *proto.HelloRequest2, rsp *proto.HelloResponse2) error {
	log.Print("Received Say.Hello request")
	rsp.Greeting = "Hello " + req.Name
	return nil
}
