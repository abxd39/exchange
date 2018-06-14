package client

import (
	"context"
	cf "digicon/gateway/conf"
	. "digicon/gateway/log"
	proto "digicon/proto/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

type WalletRPCCli struct {
	conn proto.Gateway2WallerService
}

func NewWalletRPCCli() (w *WalletRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("wallet.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_wallet")
	greeter := proto.NewGateway2WallerService(service_name, service.Client())
	w = &WalletRPCCli{
		conn: greeter,
	}
	return

}

func (this *WalletRPCCli) Callhello(name string) (rsp *proto.HelloResponse2, err error) {
	rsp, err = this.conn.Hello(context.TODO(), &proto.HelloRequest2{Name: name})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallCreateWallet(userid int, tokenid int) (rsp *proto.CreateWalletResponse, err error) {
	rsp, err = this.conn.CreateWallet(context.TODO(), &proto.CreateWalletRequest{Userid: int32(userid), Tokenid: int32(tokenid)})
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
