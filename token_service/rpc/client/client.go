package client

var InnerService *RPCClient

type RPCClient struct {
	UserSevice   *UserRPCCli
	//PublicSevice *PublciRPCCli
	PriceService *PriceRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{
		UserSevice:   NewUserRPCCli(),
		//PublicSevice: NewPublciRPCCli(),
		PriceService: NewPriceRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()

	/*
		d := make([]model.QuenesConfig, 0)
		err := DB.GetMysqlConn().Find(&d)
		if err != nil {
			log.Fatalln(err.Error())
		}

		ids:=make([]int32,0)
		for _,v:=range d{
			ids=append(ids,int32(v.TokenId))
			ids=append(ids,int32(v.TokenTradeId))
		}

		rsp,err:=InnerService.PublicSevice.CallGetTokensList(ids)
		if err!=nil {
			log.Fatalln(err.Error())
		}

		t:=make(map[int]*proto.TokensData)
		for _,v:=range rsp.Tokens {
			t[int(v.TokenId)]=v
		}

		model.GetQueneMgr().Init(d,t)
	*/

}
