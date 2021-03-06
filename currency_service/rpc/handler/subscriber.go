package handler


import (
	proto "digicon/proto/rpc"
	"digicon/currency_service/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Subscriber struct{}

func (sub *Subscriber) Process(ctx context.Context, data *proto.CnyPriceResponse) error {
	log.Println("Picked up a new message")
	log.Info("Picked up a new message")
	for _, v := range data.Data {
		model.CnyPriceMap[v.TokenId] = v
	}
	return nil
}

