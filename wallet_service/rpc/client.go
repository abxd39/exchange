package rpc

import "digicon/wallet_service/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	UserSevice *client.UserRPCCli
	TokenSevice *client.TokenRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserSevice: client.NewUserRPCCli(),
		TokenSevice:client.NewTokenRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
