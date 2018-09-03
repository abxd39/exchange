package watch

import (
	"github.com/tidwall/gjson"
	"strconv"
	"github.com/ouqiang/timewheel"
	"time"
	"encoding/json"
	"digicon/wallet_service/utils"
	. "digicon/wallet_service/model"
	log "github.com/sirupsen/logrus"
	"digicon/common/convert"
	"github.com/pkg/errors"
	"fmt"
	"math"
)

type USDTTiBiWatch struct {
	usdtCheckTranNewTW *timewheel.TimeWheel  //时间轮，检查交易状态
	Url string //节点链接
}

const (
	USDT_CHECK_LIST_KEY = "usdt_check_list_key"
	USDT_UPDATE_INTERVAL_TW = 10 //时间轮定时器间隔时间
	USDT_PROPERTYID = 1  //usdtid
	USDT_TIBI_SEND_BTC = 0.00000546  //usdt提币需要转出的比特币，这个加上手续费就是BTC总支出
)

//交易信息
type USDTTranInfo struct {
	Txid string `json:"txid"`
	Fee string `json:"fee"`
	Sendingaddress string `json:"sendingaddress"`
	Referenceaddress string `json:"referenceaddress"`
	Ismine bool `josn:"ismine"`
	Version int `json:"version"`
	Type_int int `json:"type_int"`
	Type string `json:"type"`
	Propertyid int `json:"propertyid"`
	Divisible bool `json:"divisible"`
	Amount string `json:"amount"`
	Valid bool `json:"valid"`
	Blockhash string `json:"blockhash"`
	Blocktime int64 `json:"blocktime"`
	Positioninblock int64 `json:"positioninblock"`
	Block int64 `json:"block"`
	Confirmations int64 `json:"confirmations"`
}



func StartUSDTTiBiCheckNew() {
	p := new(USDTTiBiWatch)
	p.Init()
}

//初始化
func (p *USDTTiBiWatch) Init() {

	tokenModel := new(Tokens)

	exists, err := tokenModel.GetByName("USDT")
	if err != nil {
		log.Info("init error",err)
	}
	if !exists {
		log.Info("token not exists usdt ...")
	}
	p.Url = tokenModel.Node

	//初始化同步区块时间轮
	p.usdtCheckTranNewTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		log.Info("start usdt.check.watch.new...")
		//处理交易验证
		p.checkTransactionDeal()
		//继续添加定时器
		p.usdtCheckTranNewTW.AddTimer(ETH_UPDATE_INTERVAL_TW * time.Second, "usdt_check_tibi", timewheel.TaskData{})
	})
	p.usdtCheckTranNewTW.Start()
	//开始一个事件处理
	p.usdtCheckTranNewTW.AddTimer(ETH_UPDATE_INTERVAL_TW * time.Second, "usdt_check_tibi", timewheel.TaskData{})


	tokens := new(Tokens)
	boo,err := tokens.GetByid(18)
	log.Info("---------",boo,err)
	if boo != true || err != nil || tokens.Id <= 0 {
		log.Printf("get token by id error,tokenid:%d,error:%s",err)
	}

}

//处理交易验证
func (p *USDTTiBiWatch) checkTransactionDeal() (err error) {

	var txhash string

	err,txhash = p.GetDataFromRedis()
	log.Info("---------------------从redis中读取数据：",err,txhash)
	if err != nil {
		log.Info("USDT查询redis txhash失败：",err)
		return
	}

	//查询交易数据
	err,tranInfo := utils.UsdtOmniGettransaction(p.Url,txhash)
	if err != nil {
		log.Error("USDT从链上查询交易数据失败：",err,p.Url,txhash)
		p.PushRedisList(txhash)
		return
	}
	propertyid := gjson.Get(tranInfo,"propertyid").Int()
	if propertyid != USDT_PROPERTYID {
		log.Info("USDT不是USDT代币",propertyid)
		p.PushRedisList(txhash)
		//不是usdt代币交易
		return
	}
	confirmations := gjson.Get(tranInfo,"confirmations").Int()
	if confirmations < 6 {
		log.Error("USDT未达到六次确认：",confirmations)
		//未达到六次确认
		err = errors.New("confirmations < 6")
		p.PushRedisList(txhash)
		return
	}

	log.Info("00000000000000000000000000000USDT开始记录",tranInfo)

	//交易成功，处理数据
	//手续费
	fee := gjson.Get(tranInfo,"fee").Float()
	fee = math.Abs(fee) + USDT_TIBI_SEND_BTC  //总支出btc
	realFee := convert.Float64ToInt64By8Bit(fee)
	//修改提币申请订单
	_,err = new(TokenInout).BteUpdateAppleDone2(txhash,realFee)
	if err != nil {
		log.Error("修改提币申请单报错：",err)
		return
	}
	//解析数据
	var data USDTTranInfo
	err = json.Unmarshal([]byte(tranInfo),&data)
	if err != nil {
		log.Info("eth tibi unmatshal error",err)
		return
	}
	//写记录处理
	p.USDTDeal(data)

	//汇总手续费
	new(Common).GatherFee(txhash)

	err = nil
	return
}

//保存数据到redis队列
func (p *USDTTiBiWatch) PushRedisList(txhash string) {
	redis := utils.Redis
	log.Info("收到一个USDT交易监控：",txhash)
	err := redis.RPush(USDT_CHECK_LIST_KEY,txhash).Err()
	if err != nil {
		log.Error("usdt PushRedisList error:",txhash,err)
	}
}

//从redis队列读取数据
func (p *USDTTiBiWatch) GetDataFromRedis() (error,string) {
	redis := utils.Redis
	query := redis.LPop(USDT_CHECK_LIST_KEY)
	if query.Err() != nil {
		log.Error(query.Err())
		return query.Err(),""
	}
	data := query.Val()
	return nil,data
}

//写一条数据到链记录表中
func (p *USDTTiBiWatch) WriteUSDTChainTx(data USDTTranInfo) {
	//交易是否已经收录
	exist, err := new(TokenChainInout).TxhashExist(data.Txid,0)

	if err != nil {
		log.Info("WriteUSDTChainTx error",exist, err)
		return
	}
	if exist {
		log.Info("WriteUSDTChainTx exists",exist, err)
		return
	}
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Txid)
	if err != nil {
		log.Error("WriteUSDTChainTx GetByHash error",data.Txid,err)
	}

	if tokenInout.Tokenid <= 0 {
		log.Error("WriteUSDTChainTx GetByHash Tokenid error",data.Txid,err)
	}

	var opt int = 1  //提币

	//查询token数据
	value,err := convert.StringToInt64By8Bit(data.Amount)
	if err != nil {
		log.Error("StringToInt64By8Bit error",err)
		return
	}

	txmodel := &TokenChainInout{
		Txhash:    data.Txid,
		From:      data.Sendingaddress,
		To:        data.Referenceaddress,
		Value:     strconv.FormatInt(value,10),
		Type:      opt,
		Tokenid:   tokenInout.Tokenid,
		TokenName: tokenInout.TokenName,
		Chainid:tokenInout.Chainid,
		Uid:tokenInout.Uid,
	}
	row, err := txmodel.InsertThis()
	if row <= 0 || err != nil {
		log.Error("WriteETHChainTx insert error",err)
	}
}

//eth提币处理
func (p *USDTTiBiWatch) USDTDeal(data USDTTranInfo) {
	//写一条数据到链记录表中
	p.WriteUSDTChainTx(data)
	//确认消耗冻结数量
	new(Common).USDTConfirmSubFrozen(data)
}

////////////////////////////////////////////////////////USDT充币监控/////////////////////////////////////////////////////////////////

//USDT充币监控
type USDTCBiWatch struct {
	usdtCheckCBTranNewTW *timewheel.TimeWheel  //时间轮，检查交易状态
	usdtUpdateWalletTokenNewTW *timewheel.TimeWheel  //更新wallet_token数据到redis中
	Url string //节点链接
	Chainid     int
	BlockNumber int //当前处理到的高度

	WalletTokenModel *WalletToken     //钱包详情
	TxModel          *TokenChainInout //链上交易记录
	TokenInoutModel  *TokenInout      //平台交易记录
	TokenModel       *Tokens          //币种类
	ContextModel     *Context         //处理上下文
	GetWalletTokenLastTime time.Time     //获取wallet_token最后时间，用于增量更新
}

const (
	USDT_CBI_INTERVAL_TW = 10 //时间轮定时器间隔时间
	USDT_ADDRESS_INTERVAL_TW = 5 //时间轮定时器间隔时间
	USDT_CBI_ADDRESS_REDIS_KEY = "h_wallet_token"   //和以太坊共用一个
	USDT_NAME = "USDT"
)

//开始
func StartUSDTCBiWatch() {
	p := new(USDTCBiWatch)
	p.Init()
}

//初始化
func (p *USDTCBiWatch) Init() {
	//查询ETH节点
	var data = new(Tokens)
	bool, er := data.GetByName(USDT_NAME)
	if bool != true || er != nil {
		fmt.Println("start fail")
		return
	}

	p.Url = data.Node

	log.Info("init data：",data.Node,data)

	//model初始化
	//this.WalletToken = new(Blocks)
	p.WalletTokenModel = new(WalletToken)
	p.TxModel = new(TokenChainInout)
	p.TokenInoutModel = new(TokenInout)
	p.TokenModel = new(Tokens)
	p.ContextModel = new(Context)
	//获取chainid
	p.Chainid = 0

	p.BlockNumber, _ = p.ContextModel.MaxNumber(p.Url, p.Chainid)

	//初始化同步区块时间轮
	p.usdtCheckCBTranNewTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		log.Info("----------------------------------start usdt.cbi.watch.new...")
		//区块操作处理
		p.WorkerDone()
		//继续添加定时器
		p.usdtCheckCBTranNewTW.AddTimer(USDT_CBI_INTERVAL_TW * time.Second, "usdt_check_cbi", timewheel.TaskData{})
	})
	p.usdtCheckCBTranNewTW.Start()
	//开始一个事件处理
	p.usdtCheckCBTranNewTW.AddTimer(USDT_CBI_INTERVAL_TW * time.Second, "usdt_check_cbi", timewheel.TaskData{})

	//读取wallet_token数据写到redis中[和以太坊的共用，这里不需要一步更新]

}

//处理区块
func (p *USDTCBiWatch) WorkerDone() {
	//查询数据库中的区块数
	p.BlockNumber, _ = p.ContextModel.GetMaxNumber(USDT_NAME)
	//当前最高块
	err,height := utils.USDTGetBlockCount(p.Url)
	if err != nil {
		return
	}
	if p.BlockNumber <= 0 {
		p.BlockNumber = height - 10
	}
	//p.BlockNumber = hight - 10

	log.Info("usdt height：",p.BlockNumber,height)

	if p.BlockNumber <= height-6 {
		for i := p.BlockNumber + 1; i <= height-2; i++ {
			log.Info("USDT循环次数：",i)
			//p.WorkerHander(i)
			p.parseBlock(i)
			//记录当前进度
			p.ContextModel.SaveMaxNumber(p.Url, p.Chainid, USDT_NAME,i)

		}
	}
}

//解析区块
func (p *USDTCBiWatch) parseBlock(num int) error {
	log.Info("parse USDT TX:",num)
	err,blockTx := utils.USDTOmniListblocktransactions(p.Url,num)
	log.Info("解析usdt区块：",err,blockTx)
	if err != nil {
		log.Error("usdt error",err)
		return err
	}
	if len(blockTx) == 0 {
		return errors.New("not found usdt tx")
	}

	for _,v := range blockTx {
		//根据交易hash，查询交易详情
		p.parseTx(v)
	}
	return nil
}

//解析交易
func (p *USDTCBiWatch) parseTx(txhash string) {
	log.Info("开始解析交易数据，parseTx：",txhash)
	//查询交易数据
	err,tranInfo := utils.UsdtOmniGettransaction(p.Url,txhash)
	if err != nil {
		return
	}
	propertyid := gjson.Get(tranInfo,"propertyid").Int()
	log.Info("USDT代币交易key:",propertyid)
	if propertyid != USDT_PROPERTYID {
		//不是usdt代币交易
		return
	}
	confirmations := gjson.Get(tranInfo,"confirmations").Int()
	log.Info("USDT代币交易确认次数:",confirmations)
	if confirmations < 6 {
		//未达到六次确认
		err = errors.New("confirmations < 6")
		return
	}
	//交易成功，处理数据
	//解析数据
	var data USDTTranInfo
	err = json.Unmarshal([]byte(tranInfo),&data)
	if err != nil {
		log.Info("usdt tibi unmatshal error",err)
		return
	}
	//新增订单并处理添加token
	boo,err := p.newOrder(data)
	fmt.Println("usdt add order result",boo,err)
}

//新增订单记录
func (p *USDTCBiWatch) newOrder(data USDTTranInfo) (bool,error) {
	//交易是否已经收录
	exist, err := p.TxModel.TxhashExist(data.Txid,0)
	if err != nil {
		fmt.Println("txhash not exists",data.Txid,p.Chainid,exist,err)
		return false, err
	}
	if exist {
		fmt.Println("tx already exists",exist)
		return false, errors.New("tx already exists")
	}

	walletToken := new(WalletToken)
	err = walletToken.GetByAddress(data.Referenceaddress)  //根据接受者地址，获取用户数据
	if err != nil {
		return false,nil
	}
	if walletToken.Uid <= 0 {
		return false,nil
	}

	//查询tokenid
	tokens := new(Tokens)
	boo,err := tokens.GetByName(USDT_NAME)
	if err != nil {
		return false,err
	}
	if boo != true {
		return false,nil
	}

	var opt int = 1  //充币
	amount,err := convert.StringToInt64By8Bit(data.Amount)
	if err != nil {
		log.Error("StringToInt64By8Bit error",err)
	}
	value := strconv.FormatInt(amount,10)
	_,err = p.TxModel.Insert(data.Txid, data.Sendingaddress, data.Referenceaddress, value, "", 0, walletToken.Uid, tokens.Id, tokens.Mark,opt)
	if err != nil {
		log.WithFields(log.Fields{
			"txhash":data.Txid,
			"from":data.Sendingaddress,
			"to":data.Referenceaddress,
			"value":value,
			"contract":"",
			"chainid":0,
			"uid":walletToken.Uid,
			"id":tokens.Id,
			"mark":tokens.Mark,
			"opt":opt,
		}).Info("insert usdt into tx order error:",err)
	}
	_,err = p.TokenInoutModel.Insert(data.Txid, data.Sendingaddress, data.Referenceaddress, value, "", 0, walletToken.Uid, tokens.Id, tokens.Mark, tokens.Decimal,opt)
	if err != nil {
		log.WithFields(log.Fields{
			"txhash":data.Txid,
			"from":data.Sendingaddress,
			"to":data.Referenceaddress,
			"value":value,
			"contract":"",
			"chainid":0,
			"uid":walletToken.Uid,
			"id":tokens.Id,
			"mark":tokens.Mark,
			"deci":tokens.Decimal,
			"opt":opt,
		}).Info("insert usdt into inout order error:",err)
	}

	//给用户添加token
	boo,errr := new(Common).AddETHTokenNum(walletToken.Uid,data.Referenceaddress,walletToken.Tokenid,value,data.Txid)
	if boo != true {
		log.Error("AddUSDTTokenNum err:",errr)
	}
	if errr != nil {
		log.Error("add usdt error:",errr)
	}
	log.WithFields(log.Fields{
		"uid":walletToken.Uid,
		"real_uid":walletToken.Uid,
		"from":data.Sendingaddress,
		"to":data.Referenceaddress,
		"chainid":0,
		"contract":"",
		"value":value,
		"txhash":data.Txid,
	}).Info("add usdt chain complete")
	return true,nil
}
