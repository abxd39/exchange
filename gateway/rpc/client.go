package rpc

import "digicon/gateway/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	UserSevice      *client.UserRPCCli
	CurrencyService *client.CurrencyRPCCli
	TokenService    *client.TokenRPCCli
	WallService     *client.WalletRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserSevice:      client.NewUserRPCCli(),
		CurrencyService: client.NewCurrencyRPCCli(),
		TokenService:    client.NewTokenRPCCli(),
		WallService:     client.NewWalletRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
