package exchange

/*
import (
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	log "github.com/sirupsen/logrus"
	"digicon/token_service/model"
	"digicon/token_service/rpc"
)

var Symbols map[int]*proto.TokensData

func LoadSymbol(token_id int32) (err error) {

	//t := new(model.QuenesConfig).GetQuenes(ty)

	d := make([]model.QuenesConfig, 0)
	err = DB.GetMysqlConn().Where("switch=", 1).Find(&d)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	ids := make([]int32, 0)
	for _, v := range d {
		ids = append(ids, int32(v.TokenId))
		ids = append(ids, int32(v.TokenTradeId))
	}

	h, err := rpc.InnerService.PublicSevice.CallGetTokensList(ids)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	//Symbols = make(map[int]*proto.TokensData)
	for _, v := range h.Tokens {
		Symbols[int(v.TokenId)] = v
	}
	return nil
}

func init() {
	Symbols = make(map[int]*proto.TokensData)
}
*/

/*
var InitChan chan model.ConfigQuenes

func InitExchange() {
	fmt.Println("begin to load all config")
	InitChan = make(chan model.ConfigQuenes, 1000)

	d := new(model.ConfigQuenes).GetAllQuenes()
	for _, v := range d {
		InitChan <- v
	}

	for v := range InitChan {
		rsp, err := client.InnerService.PriceService.CallLastPrice(v.Name)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(1 * time.Second)
			InitChan <- v
			continue
		}
		var e *model.EntrustQuene
		cny := model.GetTokenCnyPrice(v.TokenId)
		if !rsp.Ok {
			e = model.NewEntrustQueue(v.TokenId, v.TokenTradeId,100000000, v.Name, cny, 0, 0, 0)
		}else {
			data := rsp.Data
			e = model.NewEntrustQueue(v.TokenId, v.TokenTradeId, data.Price, v.Name, cny, data.Amount, data.Vol, data.Count)
		}

		model.GetQueneMgr().AddQuene(e)

		if len(InitChan) == 0 {
			fmt.Println("finish load all config")
			return
		}
	}
}
*/
