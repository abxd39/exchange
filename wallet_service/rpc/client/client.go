package client


var InnerService *RPCClient

type RPCClient struct {
	//UserSevice *UserRPCCli
	TokenSevice *TokenRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		//UserSevice: NewUserRPCCli(),
		TokenSevice:NewTokenRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
