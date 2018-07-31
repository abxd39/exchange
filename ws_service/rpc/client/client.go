package client


var InnerService *RPCClient

type RPCClient struct {
	UserService *UserRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserService: NewUserRPCCli(),
	}
	return c
}


func InitInnerService(){
	InnerService = NewRPCClient()
}

