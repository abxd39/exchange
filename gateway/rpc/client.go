package rpc

import "digicon/gateway/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	UserSevice      *client.UserRPCCli
	CurrencyService *client.CurrencyRPCCli
	TokenService    *client.TokenRPCCli
	WallService     *client.WalletRPCCli
	PublicService   *client.PublciRPCCli
	WalletSevice    *client.WalletRPCCli
	KineService    *client.KlineRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{

		UserSevice:      client.NewUserRPCCli(),
		CurrencyService: client.NewCurrencyRPCCli(),
		TokenService:    client.NewTokenRPCCli(),
		WallService:     client.NewWalletRPCCli(),
		PublicService:   client.NewPublciRPCCli(),
		WalletSevice:    client.NewWalletRPCCli(),
		KineService:     client.NewKlineRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
