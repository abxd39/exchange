package client

var InnerService *RPCClient

type RPCClient struct {
	TokenService *TokenRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		TokenService: NewTokenRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
