package client

//import "digicon/price_service/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	TokenSevice *TokenRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		TokenSevice: NewTokenRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()

}
