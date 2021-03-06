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

type PriceRPCCli struct {
	conn proto.PriceRPCService
}

func (s *PriceRPCCli) CallCurrentPrice(p *proto.CurrentPriceRequest) (rsp *proto.CurrentPriceResponse, err error) {
	rsp, err = s.conn.CurrentPrice(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *PriceRPCCli) CallSymbols(p *proto.NullRequest) (rsp *proto.SymbolsResponse, err error) {
	rsp, err = s.conn.Symbols(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *PriceRPCCli) CallQuotation(p *proto.QuotationRequest) (rsp *proto.QuotationResponse, err error) {
	rsp, err = s.conn.Quotation(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}


func (s *PriceRPCCli) CallSymbolsTitle() (rsp *proto.SymbolTitleResponse, err error) {
	rsp, err = s.conn.SymbolTitle(context.TODO(), &proto.NullRequest{})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *PriceRPCCli) CallSymbolsById(p *proto.SymbolsByIdRequest) (rsp *proto.SymbolsByIdResponse, err error) {
	rsp, err = s.conn.SymbolsById(context.TODO(),p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (s *PriceRPCCli) CallCnyPrices(p *proto.CnyPriceRequest) (rsp *proto.CnyPriceResponse, err error) {
	rsp, err = s.conn.GetCnyPrices(context.TODO(),p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func NewPriceRPCCli() (u *PriceRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("token.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_price")
	greeter := proto.NewPriceRPCService(service_name, service.Client())
	u = &PriceRPCCli{
		conn: greeter,
	}
	return
}

func (s *PriceRPCCli) CallVolume(p *proto.VolumeRequest) (rsp *proto.VolumeResponse, err error) {
	rsp, err = s.conn.Volume(context.TODO(), p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
