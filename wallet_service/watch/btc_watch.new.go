package watch

import(
	."digicon/wallet_service/model"
	"fmt"
	"github.com/ouqiang/timewheel"
	"time"
	"digicon/wallet_service/utils"
	"encoding/json"
	"strconv"
	log "github.com/sirupsen/logrus"
	"github.com/shopspring/decimal"
	"digicon/common/convert"
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
func (p *BtcWatch) TranDeal(data TranItem) bool {
	//判断地址是否存在
	walletToken := new(WalletToken)
	boo,err := walletToken.CheckExists2(data.Address)
	if err != nil || boo != true {
		return false
	}
	//判断交易是否存在
	exists, err := p.tkChainInOutModel.TxIDExist(data.Txid)
	if exists == true || err != nil {
		log.Error("交易id存在：",data.Txid)
		return false
	}
	if data.Category != "send" && data.Category != "receive" {
		return false
	}
	//判断确认次数
	if data.Confirmations < 3 {
		return false
	}
	//写入交易记录到链记录表
	p.WriteChainTx(data)

	log.Info("比特币交易处理：",data.Txid)

	if data.Category == "send" {  //提币
		//更新完成状态
		_,err := new(TokenInout).BteUpdateAppleDone2(data.Txid,p.GetFee(data.Fee))
		if err != nil {
			log.Error("更新比特币申请状态失败：",data.Txid)
			return false
		}
		log.Info("btc send update complete：",data.Txid)
		//确认消耗冻结
		new(Common).BTCConfirmSubFrozen(data)
		//汇总手续费
		new(Common).GatherFee(data.Txid)
	}
	if data.Category == "receive" {  //充币
		log.Info("比特币充币：",data)
		//更新用户账户数量
		new(Common).AddBTCTokenNum(data)
		//添加一条充币记录到表：token_inout
		p.WriteBtcInRecord(data)
	}

	return true
}

//计算比特币提币手续费
func (p *BtcWatch) GetFee(fee float64) int64 {
	aa := decimal.NewFromFloat(fee)
	bb := decimal.NewFromFloat(float64(100000000))
	return aa.Mul(bb).Abs().IntPart()
}

//写入充币记录
func (p *BtcWatch) WriteBtcInRecord(data TranItem) {
	log.Info("写入比特币充币记录：",data)

	//交易是否已经收录
	exist, errr := new(TokenInout).TxhashExist(data.Txid, 0)

	if errr != nil {
		return
	}
	if exist {
		return
	}

	//tmp1,_ := new(big.Int).SetString(data.Amount,10)
	//amount := decimal.NewFromBigInt(tmp1, int32(8)).IntPart()
	//amount := int64(data.Amount * 100000000)

	t1 := decimal.NewFromFloat(data.Amount)
	t1_c := decimal.NewFromFloat(100000000)
	amount := t1.Mul(t1_c).IntPart()


	var inOutToken = new(TokenInout)

	var walletToken = new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil {
		log.Error("WriteBtcInRecord address not exists",err.Error())
		return
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
	}
	log.Info("交易已添加",affected)
}

//写入链交易记录
func (p *BtcWatch) WriteChainTx(data TranItem) {
	log.Info("写入链交易记录：",data)
	//交易是否已经收录
	exist, err := new(TokenChainInout).TxhashExist(data.Txid,0)

	if err != nil {
		return
	}
	if exist {
		return
	}
	//写入数据
	//新增数据
	var opt int = 1
	if data.Category == "send" {
		opt = 2 //提币
	} else if data.Category == "receive" {
		opt = 1 //充币
	}

	amount := convert.Float64ToInt64By8Bit(data.Amount)
	amount1 := strconv.FormatInt(amount,10)

	txmodel := &TokenChainInout{
		Txhash:    data.Txid,
		Address:   data.Address,
		Value:     amount1,
		Type:      opt,
		Tokenid:   2,
		TokenName: "BTC",  //这里仅仅是为了记录txhash用，所以TokenName并无实际用途
	}
	row, err := txmodel.InsertThis()
	if row <= 0 || err != nil {
		fmt.Println(err.Error())
	}
}
