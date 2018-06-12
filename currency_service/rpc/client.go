package rpc

import "digicon/currency_service/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	UserSevice *client.UserRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserSevice: client.NewUserRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
