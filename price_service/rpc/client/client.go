package client

//import "digicon/price_service/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	UserSevice *UserRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserSevice: NewUserRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
