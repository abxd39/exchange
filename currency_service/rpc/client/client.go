package client


var InnerService *RPCClient

type RPCClient struct {
	UserSevice *UserRPCCli
	//PriceRPCCli *PriceRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserSevice: NewUserRPCCli(),
		//PriceRPCCli: NewPriceRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
