package model

type MatchTrade struct {
}

func (s *MatchTrade) process() {
	for {
		GetQueneMgr().CallBackFunc(func(quene *EntrustQuene) {
			buyer, err := quene.GetFirstEntrust(0)
			if err != nil {
				return
			}

			seller, err := quene.GetFirstEntrust(ORDER_OPT_BUY)
			if err != nil {
				return
			}

			if buyer.OnPrice == seller.OnPrice {

			}
		})
	}

}
