package handler

import (
	"context"
	proto "digicon/proto/rpc"

	"fmt"
)

func (s *WalletHandler) CreateUSDTWallet(ctx context.Context, req *proto.CreateWalletRequest, rsp *proto.CreateWalletResponse) error {
	fmt.Println("create usdt wallet ...")
	var err error
	fmt.Println(err)

	return nil
}
