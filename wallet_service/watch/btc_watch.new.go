package watch

import(
	."digicon/wallet_service/model"
	"fmt"
	"github.com/ouqiang/timewheel"
	"time"
	"digicon/wallet_service/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/shopspring/decimal"
	"github.com/pkg/errors"
	"strconv"
)

//新版监听区块变化，更新数据
//包括转入和转出
//-------------------逻辑-----------------------
//使用时间轮定时器驱动
//拉取区块交易数据，判断交易id是否存在，不存在则进行如下操作
//1、判断是否提币或充币，如果满足其一，则进行如下操作
//2、记录数据到token_chain_inout，用于下次判断
//3、如果是提币，需要更新提币申请表
//4、如果是充币，需要增加用户账户余额
type BtcWatch struct {
	item []BtcWatchItem
	tranData []TranItem
	Url string
	updateBlockTW *timewheel.TimeWheel
	getAddressTW *timewheel.TimeWheel
	tkChainInOutModel *TokenChainInout
}

type TranItem struct {
	Account string `json:"account"`
	Address string `json:"address"`
	Category string `json:"category"`
	Amount float64 `json:"amount"`  //单位是1bitcoin，小数点之后保留八位小数
	Vout int `json:"vout"`
	Fee float64 `json:"fee"`
	Confirmations int64 `json:"confirmations"`
	Blockhash string `json:"blockhash"`
	Blockindex int `json:"blockindex"`
	Blocktime int `json:"blocktime"`
	Txid string `json:"txid"`
	Walletconflicts []string `json:"walletconflicts"`
	Time int64 `json:"time"`
	Timereceived int64 `json:"timereceived"`
	Abandoned bool `json:"abandoned"`
}

type BtcWatchItem struct {
	Uid int
	Address string
}

const (
	BTC_INTERVAL_TW = 10 //时间轮定时器间隔时间
)

//时间轮
var btcNewTW *timewheel.TimeWheel

func NewBtcWatch() *BtcWatch {
	return new(BtcWatch)
}

func StartBtcWatch() {
	btcWatchP := new(BtcWatch)
	//初始化
	btcWatchP.Init()
	log.Println("btc watch start ...")
}

//初始化
func (p *BtcWatch) Init() {
	tokenModel := new(Tokens)

	exists, err := tokenModel.GetByName("BTC")
	if err != nil {
		log.Error("init error",err)
	}
	if !exists {
		log.Error("token not exists btc ...")
	}
	p.Url = tokenModel.Node

	//初始化同步区块时间轮
	p.updateBlockTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		fmt.Println("start btc.watch.new...")
		//区块操作处理
		p.BlockUpdateDeal()
		//继续添加定时器
		p.updateBlockTW.AddTimer(BTC_INTERVAL_TW * time.Second, "btc_check_tran", timewheel.TaskData{})
	})
	p.updateBlockTW.Start()
	//开始一个事件处理
	p.updateBlockTW.AddTimer(BTC_INTERVAL_TW * time.Second, "btc_check_tran", timewheel.TaskData{})

	//初始化模型
	p.tkChainInOutModel = new(TokenChainInout)
}

//拉取数据
func (p *BtcWatch) GetTranData() {
	err,jsonData := utils.BtcListtransactions(p.Url)
	if err != nil {
		log.Error("GetTranData error",err.Error())
		return
	}
	//解析数据
	err = json.Unmarshal([]byte(jsonData),&p.tranData)
	if err != nil {
		log.Error("json unmarchal",err.Error())
		return
	}
	return
}

//区块操作处理
func (p *BtcWatch) BlockUpdateDeal() {
	//拉取交易
	p.GetTranData()
	data := p.tranData
	for _,v := range data {
		p.TranDeal(v)
	}
	//重置
	p.tranData = []TranItem{}
}

//交易处理
func (p *BtcWatch) TranDeal(data TranItem) (boo bool,err error) {
	defer func() {
		if err != nil {
			log.Error("比特币提币处理失败：",boo,err)
		}
	}()
	//判断地址是否存在
	walletToken := new(WalletToken)
	boo,err = walletToken.CheckExists2(data.Address)
	if err != nil {
		return false,err
	}
	if boo != true {
		return false,nil
	}

	if data.Category != "send" && data.Category != "receive" {
		return false,nil
	}
	//判断确认次数
	if data.Confirmations < 2 {
		return false,nil
	}

	log.Info("比特币交易处理：",data.Txid)
	log.Info("比特币交易分类：",data.Category,data)

	if data.Category == "send" {  //提币
		//判断交易是否存在
		exists, err := p.tkChainInOutModel.TxIDExist2(2,data.Txid)
		if err != nil {
			log.Error("send 交易id存在：",2,"---",data.Txid)
			return false,err
		}
		if exists == true {
			log.Error("send 交易id存在：",2,"---",data.Txid)
			return false,nil
		}
		//判断地址是否存在
		boo,err := new(TibiAddress).TiBiAddressExists(data.Address)
		if err != nil || boo != true {
			log.Error("send 比特币提币地址不存在或出错：",boo,err)
			return false,nil
		}
		//写入交易记录到链记录表
		boo,err = p.WriteChainTx(2,data)
		if boo != true || err != nil {
			log.Error("比特币写入链记录失败：",boo,err)
			return false,nil
		}
		//更新完成状态
		_,err = new(TokenInout).BteUpdateAppleDone2(data.Txid,p.GetFee(data.Fee))
		if err != nil {
			log.Error("更新比特币申请状态失败：",data.Txid)
			return false,err
		}
		log.Info("btc send update complete：",data.Txid)
		//确认消耗冻结
		err = new(Common).BTCConfirmSubFrozen(data)
		if err != nil {
			return false,err
		}
		//汇总手续费
		err = new(Common).GatherFee(data.Txid)
		if err != nil {
			return false,err
		}
		//提币成功提醒
		NewNotice().TiBiNoticeByTxHash(data.Txid)
	} else if data.Category == "receive" {  //充币
		//判断充币交易是否存在
		exists, err := p.tkChainInOutModel.TxIDExist2(1,data.Txid)
		if err != nil {
			log.Error("receive 交易id存在：",1,"------",data.Txid)
			return false,err
		}
		if exists == true {
			log.Error("receive 交易id存在：",1,"------",data.Txid)
			return false,nil
		}
		log.Info("receive 比特币充币：",1,"------",data)
		//判断地址是否是比特币地址
		boo,err := new(WalletToken).AddressExist(2,data.Address)
		if boo != true || err != nil {
			log.Error("不是比特币地址：",boo,err,2,data.Address)
			return false,nil
		}
		//写入交易记录到链记录表
		boo,err = p.WriteChainTx(1,data)
		if boo != true || err != nil {
			log.Error("比特币写入链记录失败：",boo,err)
			return false,nil
		}
		//更新用户账户数量
		err = new(Common).AddBTCTokenNum(data)
		if err != nil {
			log.Error("receive AddBTCTokenNum error:",err)
			return false,err
		}
		//添加一条充币记录到表：token_inout
		err = p.WriteBtcInRecord(1,data)
		if err != nil {
			log.Error("receive WriteBtcInRecord error:",err)
			return false,err
		}
		//添加充币提醒
		NewNotice().BTCCBiNotice(data)
	}

	return true,nil
}

//计算比特币提币手续费
func (p *BtcWatch) GetFee(fee float64) int64 {
	aa := decimal.NewFromFloat(fee)
	bb := decimal.NewFromFloat(float64(100000000))
	return aa.Mul(bb).Abs().IntPart()
}

//写入充币记录
func (p *BtcWatch) WriteBtcInRecord(optype int,data TranItem) error {
	log.Info("写入比特币充币记录：",data)

	//交易是否已经收录
	exist, errr := new(TokenInout).TxhashExist2(optype,data.Txid, 0)

	if errr != nil {
		return errr
	}
	if exist {
		return errors.New("已经存在")
	}

	//tmp1,_ := new(big.Int).SetString(data.Amount,10)
	//amount := decimal.NewFromBigInt(tmp1, int32(8)).IntPart()
	//amount := int64(data.Amount * 100000000)


	fs := strconv.FormatFloat(data.Amount, 'E', -1, 64)

	t1,_ := decimal.NewFromString(fs)
	t1_c := decimal.NewFromFloat(100000000)
	amount := t1.Mul(t1_c).IntPart()


	var inOutToken = new(TokenInout)

	var walletToken = new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil {
		log.Error("WriteBtcInRecord address not exists",err.Error())
		return err
	}

	//根据tokenid获取MARK
	tokens := new(Tokens)
	boo,err := tokens.GetByid(walletToken.Tokenid)
	tokenName := "BTC"
	if boo == true && err == nil {
		tokenName = tokens.Mark
	}

	inOutToken.Id = 0
	inOutToken.Txhash = data.Txid
	inOutToken.From = ""
	inOutToken.To = data.Address
	//inOutToken.Value = data.Amount
	inOutToken.Amount = amount
	inOutToken.Tokenid = 2
	inOutToken.TokenName = tokenName
	inOutToken.Uid = walletToken.Uid
	inOutToken.Tokenid = walletToken.Tokenid
	inOutToken.Opt = 1 ////充币
	affected, err := utils.Engine_wallet.InsertOne(inOutToken)
	if err != nil {
		log.Error("WriteBtcInRecord error",err.Error())
		return err
	}
	log.Info("交易已添加",affected)
	return nil
}

//写入链交易记录
func (p *BtcWatch) WriteChainTx(optype int,data TranItem) (bool,error) {
	log.Info("写入链交易记录：",optype,"-------",data)
	//交易是否已经收录
	exist, err := new(TokenChainInout).TxIDExist2(optype,data.Txid)

	if err != nil {
		return false,err
	}
	if exist {
		return false,nil
	}


	//amount := convert.Float64ToInt64By8Bit(data.Amount)
	//amount1 := strconv.FormatInt(amount,10)

	fs := strconv.FormatFloat(data.Amount, 'E', -1, 64)
	aa,_ := decimal.NewFromString(fs)
	bb := decimal.NewFromFloat(float64(100000000))
	amount1 := aa.Mul(bb).Abs().String()


	txmodel := &TokenChainInout{
		Txhash:    data.Txid,
		Address:   data.Address,
		Value:     amount1,
		Type:      optype,
		Tokenid:   2,
		TokenName: "BTC",  //这里仅仅是为了记录txhash用，所以TokenName并无实际用途
	}
	row, err := txmodel.InsertThis()
	if row <= 0 {
		return false,nil
	}
	if err != nil {
		return false,err
	}
	return true,nil
}
