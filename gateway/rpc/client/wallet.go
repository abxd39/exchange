package client

import (
	"context"
	cf "digicon/gateway/conf"
	proto "digicon/proto/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	log "github.com/sirupsen/logrus"
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
		log.Errorln(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallCreateWallet(userid int, tokenid int) (rsp *proto.CreateWalletResponse, err error) {
	rsp, err = this.conn.CreateWallet(context.TODO(), &proto.CreateWalletRequest{Userid: int32(userid), Tokenid: int32(tokenid)})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallSigntx(userid int, tokenid int, to string, gasprice int64, mount string) (rsp *proto.SigntxResponse, err error) {
	rsp, err = this.conn.Signtx(context.TODO(), &proto.SigntxRequest{Userid: int32(userid), Tokenid: int32(tokenid), To: to, Gasprice: (gasprice), Mount: mount})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallSendRawTx(tokenid int32, signtx string, applyId int32) (rsp *proto.SendRawTxResponse, err error) {
	rsp, err = this.conn.SendRawTx(context.TODO(), &proto.SendRawTxRequest{TokenId: tokenid, Signtx: signtx, Applyid: applyId})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallTibi(uid int32, token_id int32, to string, gasprice int32, amount string) (rsp *proto.TibiResponse, err error) {
	rsp, err = this.conn.Tibi(context.TODO(), &proto.TibiRequest{Uid: (uid), Tokenid: (token_id), To: to, Gasprice: (gasprice), Amount: amount})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallTibiApply(uid int32, token_id int32, to string, gasprice string, amount string, real_amount string, sms_code string, email_code string, password string,phone string,email string) (rsp *proto.TibiApplyResponse, err error) {
	rsp, err = this.conn.TibiApply(context.TODO(), &proto.TibiApplyRequest{Uid: (uid), Tokenid: (token_id), To: to, Gasprice: (gasprice), Amount: amount, RealAmount: real_amount,Phone:phone,Email:email,SmsCode:sms_code,EmailCode:email_code,Password:password})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallCancelTiBi(uid int32, id int32) (rsp *proto.CancelTiBiResponse, err error) {
	rsp, err = this.conn.CancelTiBi(context.TODO(), &proto.CancelTiBiRequest{Uid: (uid), Id: (id)})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallGetValue(uid int32, token_id int32) (rsp *proto.GetValueResponse, err error) {
	rsp, err = this.conn.GetValue(context.TODO(), &proto.GetValueRequest{Uid: (uid), Tokenid: (token_id)})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallAddressSave(uid int32, tokenid int32, address string, mark string) (rsp *proto.AddressSaveResponse, err error) {
	rsp, err = this.conn.AddressSave(context.TODO(), &proto.AddressSaveRequest{Uid: (uid), Tokenid: (tokenid), Address: address, Mark: mark})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallAddressList(uid int32) (rsp *proto.AddressListResponse, err error) {
	rsp, err = this.conn.AddressList(context.TODO(), &proto.AddressListRequest{Uid: int32(uid)})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
func (this *WalletRPCCli) CallAddressDelete(uid int32, id int32) (rsp *proto.AddressDeleteResponse, err error) {
	rsp, err = this.conn.AddressDelete(context.TODO(), &proto.AddressDeleteRequest{Uid: (uid), Id: id})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

////////////////// btc  ///////////////////////////
func (this *WalletRPCCli) CallBtcSigntx(req *proto.BtcSigntxRequest) (rsp *proto.BtcSigntxResponse, err error) {
	rsp, err = this.conn.BtcSigntx(context.TODO(), req)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

//
func (this *WalletRPCCli) CallBtcTibi(req *proto.BtcTibiRequest) (rsp *proto.BtcResponse, err error) {
	rsp, err = this.conn.BtcTibi(context.TODO(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallInList(req *proto.InListRequest) (rsp *proto.InListResponse, err error) {
	rsp, err = this.conn.InList(context.TODO(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallOutList(req *proto.OutListRequest) (rsp *proto.OutListResponse, err error) {
	rsp, err = this.conn.OutList(context.TODO(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallGetAddress(userid int, tokenid int) (rsp *proto.GetAddressResponse, err error) {
	rsp, err = this.conn.GetAddress(context.TODO(), &proto.GetAddressRequest{Userid: int32(userid), Tokenid: int32(tokenid)})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallSyncBlockTx(block int32,key string) (rsp *proto.SyncEthBlockTxResponse, err error) {
	rsp, err = this.conn.SyncEthBlockTx(context.TODO(), &proto.SyncEthBlockTxRequest{Block:block,Key:key})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallGetOutTokenFee() (rsp *proto.GetOutTokenFeeResponse, err error) {
	rsp, err = this.conn.GetOutTokenFee(context.TODO(), &proto.GetOutTokenFeeRequest{})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (this *WalletRPCCli) CallCancelSubTokenWithFronze(uid int32,token_id int32,num int64,ukey string,key string) (rsp *proto.CancelSubTokenWithFronzeResponse, err error) {
	rsp, err = this.conn.CancelSubTokenWithFronze(context.TODO(), &proto.CancelSubTokenWithFronzeRequest{
		Uid:uid,
		Tokenid:token_id,
		Num:num,
		Ukey:ukey,
		Key:key,
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
