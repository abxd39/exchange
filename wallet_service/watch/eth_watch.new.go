package watch

import (
	"github.com/ouqiang/timewheel"
	"time"
	"fmt"
	"digicon/wallet_service/model"
	"digicon/wallet_service/utils"
	"github.com/tidwall/gjson"
	"math/big"
	"github.com/shopspring/decimal"
	"encoding/json"
	. "digicon/wallet_service/model"
	"strconv"
	"bytes"
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
	log "github.com/sirupsen/logrus"
)

//逻辑
//分两部分：充币和提币
//提币验证逻辑：
//	1、提供一个函数，供外部写入交易hash，存储到redis中(防止因程序挂掉，导致进程内存数据丢失)
//	2、开启时间轮，定时从redis中读取数据验证，验证结果需要修改提币申请表，同时记录区块信息
//充币逻辑：
//	1、按照之前的修改，会涉及到频繁访问数据库的问题，需要优化

type EthTiBiWatch struct {
	ethCheckTranNewTW *timewheel.TimeWheel  //时间轮，检查交易状态
	Url string //节点链接
}

const (
	ETH_CHECK_LIST_KEY = "eth_check_list_key"
	ETH_UPDATE_INTERVAL_TW = 10 //时间轮定时器间隔时间
)

//交易信息
type TranInfo struct {
	Hash string `json:"hash"`
	From string `json:"from"`
	To string `json:"to"`
	Value string `json:"value"`
	Input string `josn:"input"`
	Gas string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

func StartEthCheckNew() {
	p := new(EthTiBiWatch)
	p.Init()
}

//初始化
func (p *EthTiBiWatch) Init() {

	tokenModel := new(models.Tokens)

	exists, err := tokenModel.GetByName("ETH")
	if err != nil {
		log.Info("init error",err)
	}
	if !exists {
		log.Info("token not exists btc ...")
	}
	p.Url = tokenModel.Node

	//初始化同步区块时间轮
	p.ethCheckTranNewTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		log.Info("start eth.check.watch.new...")
		//处理交易验证
		p.checkTransactionDeal()
		//继续添加定时器
		p.ethCheckTranNewTW.AddTimer(ETH_UPDATE_INTERVAL_TW * time.Second, "eth_check_tibi", timewheel.TaskData{})
	})
	p.ethCheckTranNewTW.Start()
	//开始一个事件处理
	p.ethCheckTranNewTW.AddTimer(ETH_UPDATE_INTERVAL_TW * time.Second, "eth_check_tibi", timewheel.TaskData{})


	tokens := new(Tokens)
	boo,err := tokens.GetByid(18)
	log.Info("---------",boo,err)
	if boo != true || err != nil || tokens.Id <= 0 {
		log.Printf("get token by id error,tokenid:%d,error:%s",err)
	}

}

//处理交易验证
func (p *EthTiBiWatch) checkTransactionDeal() {
	err,txhash := p.GetDataFromRedis()
	if err != nil {
		return
	}
	//开始验证数据
	err,data := utils.RpcGetTransactionReceipt(p.Url,txhash)
	if err != nil {
		return
	}
	status := gjson.Get(data,"status").String()
	if status == "" {
		//暂未打包成功，重新放入队列，等待下次执行
		p.PushRedisList(txhash)
		return
	}

	//计算实际消耗的手续费
	gasUsed := gjson.Get(data,"gasUsed").String()
	gasUsedTemp, _ := new(big.Int).SetString(gasUsed[2:], 16)
	realFee := decimal.NewFromBigInt(gasUsedTemp, 0).IntPart()

	temp, _ := new(big.Int).SetString(status[2:], 16)
	statuss := decimal.NewFromBigInt(temp, 0).IntPart()
	log.Info("package status：",statuss)
	if statuss == 1 {
		//进行六次确认
		if p.TranSeixCheck(p.Url,data,txhash) == false {
			//未达到六次确认
			p.PushRedisList(txhash)
			return
		}
		//成功
		//查询交易信息
		err,tranInfo := utils.RpcGetTransaction(p.Url,txhash)
		if err != nil {
			return
		}
		//解析数据
		var data TranInfo
		err = json.Unmarshal([]byte(tranInfo),&data)
		if err != nil {
			log.Info("eth tibi unmatshal error",err)
			return
		}

		//修改提币申请订单
		_,err = new(models.TokenInout).BteUpdateAppleDone2(txhash,realFee)
		if err != nil {
			log.Error("修改提币申请单报错：",err)
			return
		}

		//判断是否代币转账
		if strings.Count(data.Input, "") < 138 || strings.Compare(data.Input[0:10], "0xa9059cbb") != 0 {
			//eth转账
			p.ETHDeal(data)
		} else {
			//erc20代币转账
			p.ERC20Deal(data)
		}
		//汇总手续费
		new(Common).GatherFee(txhash)
		return
	}
	//暂未打包成功，重新放入队列，等待下次执行
	p.PushRedisList(txhash)
}

func (p *EthTiBiWatch) formatGas(gas string,gasPrice string) (int64,int64,int64) {
	gas_num, err := strconv.ParseInt(gas, 0, 64)
	if err != nil {
		fmt.Println(err)
		gas_num = 0
	}
	gas_price, err := strconv.ParseInt(gasPrice, 0, 64)
	if err != nil {
		fmt.Println(err)
		gas_price = 0
	}

	t1 := decimal.NewFromFloat(float64(gas_num))
	t1_c := decimal.NewFromFloat(float64(gas_price))
	real_fee := t1.Mul(t1_c).IntPart()
	return gas_num,gas_price,real_fee
}

//保存数据到redis队列
func (p *EthTiBiWatch) PushRedisList(txhash string) {
	redis := utils.Redis
	log.Info("收到一个交易监控：",txhash)
	err := redis.RPush(ETH_CHECK_LIST_KEY,txhash).Err()
	if err != nil {
		log.Error("PushRedisList error:",txhash,err)
	}
}

//从redis队列读取数据
func (p *EthTiBiWatch) GetDataFromRedis() (error,string) {
	redis := utils.Redis
	query := redis.LPop(ETH_CHECK_LIST_KEY)
	if query.Err() != nil {
		log.Error(query.Err())
		return query.Err(),""
	}
	data := query.Val()
	return nil,data
}

//写入一条订单数据到表：token_chain_inout
func (p *EthTiBiWatch) WriteEthInRecord(data TranInfo) {
	//tmp1,_ := new(big.Int).SetString(data.Value,10)
	//value := decimal.NewFromBigInt(tmp1, int32(8)).IntPart()

	var inOutToken = new(models.TokenInout)

	var walletToken = new(models.WalletToken)
	err := walletToken.GetByAddress(data.From)
	if err != nil {
		log.Error("WriteEthInRecord WriteBtcInRecord address not exists",err)
		return
	}

	inOutToken.Id = 0
	inOutToken.Txhash = data.Hash
	inOutToken.From = data.From
	inOutToken.To = data.To
	inOutToken.Opt = 1  //提币
	//inOutToken.Value = value
	//inOutToken.Amount = value
	inOutToken.Tokenid = 3  //以太坊tokenid
	inOutToken.TokenName = "ETH"
	inOutToken.Uid = walletToken.Uid
	inOutToken.Tokenid = walletToken.Tokenid
	affected, err := utils.Engine_wallet.InsertOne(inOutToken)
	if err != nil {
		log.Error("WriteEthInRecord error",err)
	}
	log.Println(affected)
}

//写一条数据到链记录表中
func (p *EthTiBiWatch) WriteERC20ChainTx(data TranInfo) {
	//交易是否已经收录
	exist, err := new(models.TokenChainInout).TxhashExist(data.Hash,0)
	if err != nil {
		return
	}
	if exist {
		return
	}
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Hash)
	if err != nil || tokenInout.Tokenid <= 0 {
		log.Error("WriteERC20ChainTx GetByHash error",data.Hash,err)
	}
	var opt int = 1  //提币

	//查询token数据
	tokens := new(Tokens)
	var boo bool
	boo,err = tokens.GetByid(tokenInout.Tokenid)
	if boo != true || err != nil || tokens.Id <= 0 {
		log.Error("WriteERC20ChainTx get token by id error",tokenInout.Tokenid,err)
	}

	//提币地址
	to := strings.Join([]string{"0x",data.Input[35:74]},"")

	//查询用户wallet_token
	//walletToken := new(models.WalletToken)
	//err = walletToken.GetByAddress(data.From)
	//if err != nil || walletToken.Uid <= 0 {
	//	log.Error("WriteERC20ChainTx GetByAddress error",err,walletToken.Uid,walletToken)
	//	return
	//}

	//数量

	amount := strings.TrimLeft(data.Input[74:],"0")
	temp, boo := new(big.Int).SetString(amount,16)
	if boo != true {
		log.Error("format error",amount)
	}
	value := decimal.NewFromBigInt(temp, int32(8 - tokens.Decimal)).String()

	txmodel := &models.TokenChainInout{
		Txhash:    data.Hash,
		From:      data.From,
		To:        to,
		Value:     value,
		Type:      opt,
		Tokenid:   tokenInout.Tokenid,
		TokenName: tokenInout.TokenName,
		Chainid:tokenInout.Chainid,
		Uid:tokenInout.Uid,
		Contract:data.To,
	}
	row, err := txmodel.InsertThis()
	if row <= 0 || err != nil {
		log.Error("WriteERC20ChainTx insert error",err)
	}
}

//写一条数据到链记录表中
func (p *EthTiBiWatch) WriteETHChainTx(data TranInfo) {
	//交易是否已经收录
	exist, err := new(models.TokenChainInout).TxhashExist(data.Hash,0)

	if err != nil {
		log.Info("WriteETHChainTx error",exist, err)
		return
	}
	if exist {
		log.Info("WriteETHChainTx exists",exist, err)
		return
	}
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(data.Hash)
	if err != nil || tokenInout.Tokenid <= 0 {
		log.Error("WriteETHChainTx GetByHash error",data.Hash,err)
	}
	var opt int = 1  //提币

	//查询token数据
	tokens := new(Tokens)
	var boo bool
	boo,err = tokens.GetByid(tokenInout.Tokenid)
	if boo != true || err != nil || tokens.Id <= 0 {
		log.Error("get token by id error,tokenid:%d,error:%s",tokenInout.Tokenid,err)
	}

	//查询用户wallet_token
	//walletToken := new(models.WalletToken)
	//err = walletToken.GetByAddress(data.From)
	//if err != nil || walletToken.Uid <= 0 {
	//	log.Error("WriteETHChainTx GetByAddress error",err)
	//	return
	//}

	//格式化数量
	temp, boo := new(big.Int).SetString(data.Value[2:],16)
	if boo != true {
		log.Error("format data error:",data.Value)
	}
	value := decimal.NewFromBigInt(temp, int32(8 - tokens.Decimal)).String()

	txmodel := &models.TokenChainInout{
		Txhash:    data.Hash,
		From:      data.From,
		To:        data.To,
		Value:     value,
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
func (p *EthTiBiWatch) ETHDeal(data TranInfo) {
	contract := ""
	//写一条数据到链记录表中
	p.WriteETHChainTx(data)
	//确认消耗冻结数量
	new(Common).ETHConfirmSubFrozen(data.From,data.Hash,contract)
}

//erc代币处理
func (p *EthTiBiWatch) ERC20Deal(data TranInfo) {
	//data.Hash = strings.Join([]string{"0x",data.Input[35:74]},"")
	contract := data.To
	//写一条数据到链记录表中
	p.WriteERC20ChainTx(data)
	//确认消耗冻结数量
	new(Common).ETHConfirmSubFrozen(data.From,data.Hash,contract)
}

//交易六次确认检查
func (p *EthTiBiWatch) TranSeixCheck(url string,receipt string,txhash string) bool {
	blockNumber := gjson.Get(receipt,"blockNumber").String()
	if blockNumber == "" {
		return false
	}
	number, err := strconv.ParseInt(blockNumber, 0, 64)
	if err != nil {
		log.Error("blockNumber is null")
		return false
	}
	//查询当前区块数
	//当前最高块
	height, err := utils.GetBlockNumber(url)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Info("区块比较：",height,number)
	if height - number < 30 {
		return false
	}
	//如果需要，还可以判断number中有没有指定的交易，暂时不需要判断，通过以上就可以确认交易成功
	//查询交易
	data,err := utils.GetblockBynumber(url,int(number))
	if err != nil {
		return false
	}

	var block map[string]interface{}
	//fmt.Println(string(ret))
	json.Unmarshal(data, &block)
	txs := block["result"].(map[string]interface{})["transactions"].([]interface{})
	for i := 0; i < len(txs); i++ {
		tx := txs[i].(map[string]interface{})
		if tx["to"] == nil { //部署合约交易直接跳过
			continue
		}
		if tx["hash"].(string) == txhash {
			log.Info("交易存在:",txhash)
			return true
		}
	}
	return false
}


////////////////////////////////以下代币是充币验证逻辑////////////////////////////////
//逻辑说明：
//按照之前的修改，大概逻辑如下：
//1、初始化模型
//2、查询数据库中的区块高度
//3、从区块链中获取最新的区块高度
//4、对比区块高度，循环最近的交易
//5、按照根据交易数据，判断是以太币转账还是ERC20代币转账
//6、写入区块数据到表：token_chain_inout中，同时写入一条充币数据到表：token_inout中

type EthCBiWatch struct {
	ethCheckCBTranNewTW *timewheel.TimeWheel  //时间轮，检查交易状态
	ethUpdateWalletTokenNewTW *timewheel.TimeWheel  //更新wallet_token数据到redis中
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
	ETH_CBI_INTERVAL_TW = 10 //时间轮定时器间隔时间
	ETH_ADDRESS_INTERVAL_TW = 5 //时间轮定时器间隔时间
	ETH_CBI_ADDRESS_REDIS_KEY = "h_wallet_token"
)

//开始
func StartEthCBiWatch() {
	p := new(EthCBiWatch)
	p.Init()
}

//初始化
func (p *EthCBiWatch) Init() {
	//查询ETH节点
	var data = new(Tokens)
	bool, er := data.GetByName("ETH")
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
	var err error
	p.Chainid, err = p.getChainid()
	if err != nil {
		fmt.Println(err)
		return
	}

	p.BlockNumber, _ = p.ContextModel.MaxNumber(p.Url, p.Chainid)

	//初始化同步区块时间轮
	p.ethCheckCBTranNewTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		log.Info("----------------------------------start eth.cbi.watch.new...")
		//区块操作处理
		p.WorkerDone()
		//继续添加定时器
		p.ethCheckCBTranNewTW.AddTimer(ETH_CBI_INTERVAL_TW * time.Second, "eth_check_cbi", timewheel.TaskData{})
	})
	p.ethCheckCBTranNewTW.Start()
	//开始一个事件处理
	p.ethCheckCBTranNewTW.AddTimer(ETH_CBI_INTERVAL_TW * time.Second, "eth_check_cbi", timewheel.TaskData{})

	//读取wallet_token数据写到redis中
	p.WriteAllWalletTokenToRedis()

	//ethUpdateWalletTokenNewTW
	p.ethUpdateWalletTokenNewTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		log.Info("start eth.cbi.wallet_token.new...")
		//更新操作
		p.WriteIncrWalletTokenToRedis()
		//继续添加定时器
		p.ethUpdateWalletTokenNewTW.AddTimer(ETH_ADDRESS_INTERVAL_TW * time.Second, "eth_wallet_token_cbi", timewheel.TaskData{})
	})
	p.ethUpdateWalletTokenNewTW.Start()
	//开始一个事件处理
	p.ethUpdateWalletTokenNewTW.AddTimer(ETH_ADDRESS_INTERVAL_TW * time.Second, "eth_wallet_token_cbi", timewheel.TaskData{})

}

//处理区块
func (p *EthCBiWatch) WorkerDone() {
	//查询数据库中的区块数
	p.BlockNumber, _ = p.ContextModel.MaxNumber(p.Url, p.Chainid)
	//当前最高块
	temp, err := p.Get_blockNumber()
	hight := int(temp)
	if err != nil {
		fmt.Println(err)
		return
	}
	if p.BlockNumber <= 0 {
			p.BlockNumber = hight - 10
		}
	//p.BlockNumber = hight - 10

	log.Info("height：",p.BlockNumber,hight)

	if p.BlockNumber < hight-30 {
		for i := p.BlockNumber + 1; i <= hight-30; i++ {
			p.WorkerHander(i)
			//记录当前进度
			p.ContextModel.Save(p.Url, p.Chainid, i)

		}
	}
}

//具体处理区块
func (p *EthCBiWatch) WorkerHander(num int) error {
	//log.Info("start WorkerHander",num)
	ret, err := p.GetblockBynumber(num)
	if err != nil {
		return err
	}
	var block map[string]interface{}
	//fmt.Println(string(ret))
	json.Unmarshal(ret, &block)
	txs := block["result"].(map[string]interface{})["transactions"].([]interface{})

	//fmt.Println(txs)
	//log.Info("start for")
	//log.Info("tx_data:",txs)
	for i := 0; i < len(txs); i++ {
		tx := txs[i].(map[string]interface{})
		if tx["to"] == nil { //部署合约交易直接跳过
			continue
		}
		//log.Info("交易数据：",txs[i])

		//检查eth转账
		ext := p.ExistsAddress(tx["to"].(string), p.Chainid, "")
		log.Info("是否存在：",ext,tx["to"].(string),",",p.Chainid)
		//fmt.Println("是否存在：",ext,tx["to"].(string) == "0x870f49783e9d8c9707a72b252a0e56d3b7628f31",p.Chainid)
		//ext, err := p.WalletTokenModel.AddrExist(tx["to"].(string), p.Chainid, "")

		//if err != nil {
		//	fmt.Println(err)
		//	return err
		//}
		if ext {
			log.Info("find_a_eth")
			//TODO:
			p.newOrder(p.WalletTokenModel.Uid, tx["from"].(string), tx["to"].(string), p.Chainid, "", tx["value"].(string), tx["hash"].(string),tx["gas"].(string),tx["gasPrice"].(string))

			continue
		}

		input := tx["input"].(string)
		//log.Info("input：",input)
		//不是token转账跳过
		if strings.Count(input, "") < 138 || strings.Compare(input[0:10], "0xa9059cbb") != 0 {
			//fmt.Println(input)
			continue
		}

		ext = p.ExistsAddress(fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string))
		//log.Info("是否存在-----：",fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string))
		//log.Info("existssss：",ext,fmt.Sprintf("0x%s", input[34:74]),p.Chainid,tx["to"].(string))
		//ext, err = p.WalletTokenModel.AddrExist(fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string))
		//fmt.Println(ext,err,this.WalletTokenModel)
		//fmt.Println("是否存在123：",ext,p.Chainid,fmt.Sprintf("0x%s", input[34:74]))
		if !ext {
			continue
		}
		var vstart int
		for i := 74; i < 138; i++ {
			if input[i:i+1] != "0" {
				vstart = i
				break
			}
		}
		if vstart == 0 {
			continue
		}
		log.Info("find_a_token")

		ok, err := p.newOrder(p.WalletTokenModel.Uid, tx["from"].(string), fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string), fmt.Sprintf("0x%s", input[vstart:138]), tx["hash"].(string),tx["gas"].(string),tx["gasPrice"].(string))
		log.Info("newOrder complete",ok, err)
		continue

	}
	return nil
}

func (p *EthCBiWatch) GetblockBynumber(num int) ([]byte, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_getBlockByNumber"
	strconv.FormatInt(int64(num), 16)
	//str:=fmt.Sprintf("0x%s",strconv.FormatInt(int64(num),16))
	//fmt.Println(num, fmt.Sprintf("0x%s", strconv.FormatInt(int64(num), 16)))
	send["params"] = []interface{}{fmt.Sprintf("0x%s", strconv.FormatInt(int64(num), 16)), true}
	send["id"] = p.Chainid
	rsp, err := p.post(send)
	//str := (*string)(unsafe.Pointer(&rsp))
	//fmt.Println(*str)

	return rsp, err
}

//获取区块高度
func (p *EthCBiWatch) Get_blockNumber() (int64, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_blockNumber"
	send["params"] = []string{}
	send["id"] = p.Chainid
	rsp, err := p.post(send)
	//str := (*string)(unsafe.Pointer(&rsp))
	//fmt.Println(*str)
	if err != nil {
		return 0, err
	}
	//
	data := make(map[string]interface{})
	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return 0, err
	}
	result, ok := data["result"]
	if !ok {
		return 0, nil
	}
	var balance string
	balance, ok = result.(string)
	number, err := strconv.ParseInt(balance, 0, 64)
	//fmt.Println(data["result"],err)

	return number, nil
}

func (p *EthCBiWatch) getChainid() (int, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "net_version"
	send["params"] = []string{}
	send["id"] = 67
	rsp, err := p.post(send)
	//str := (*string)(unsafe.Pointer(&rsp))
	//fmt.Println(*str)
	if err != nil {
		return 0, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return 0, err
	}
	result, ok := data["result"]
	if !ok {
		return 0, nil
	}
	var balance string
	balance, ok = result.(string)
	number, err := strconv.ParseInt(balance, 0, 64)
	//fmt.Println(data["result"],err)

	return int(number), nil

}

//post RPC数据
func (p *EthCBiWatch) post(send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	url := p.Url
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//byte数组直接转成string，优化内存
	return respBytes, nil
}

//新增充值订单
func (p *EthCBiWatch) newOrder(uid int, from string, to string, chainid int, contract string, value string, txhash string,gas string,gasPrice string) (bool,error) {
	log.Info("newOrder_params:",uid,",",from,",",to,",",chainid,",",contract,",",value,",",txhash)
	//交易是否已经收录
	exist, err := p.TxModel.TxhashExist(txhash, p.Chainid)
	fmt.Println("交易是否存在----------------------：",exist,err,txhash,p.Chainid)
	if err != nil {
		fmt.Println("txhash not exists",txhash,p.Chainid,exist,err)
		return false, err
	}
	if exist {
		fmt.Println("tx already exists",exist)
		return false, errors.New("tx already exists")
	}

	//
	//根据to查询数据

	//查询uid
	walletToken := new(models.WalletToken)
	//err = walletToken.GetByAddress(to)
	var bo bool
	bo,err = walletToken.GetByAddressContract(to,contract)
	log.Info("walletToken_data",walletToken.Uid,",",walletToken.Tokenid,",",to,",",contract)
	if err != nil || bo != true || walletToken.Uid <= 0 {
		log.Error("get uid by address error",err,to,contract)
		return false,err
	}

	tokensModel := new(Tokens)

	boo, err := tokensModel.GetByid(walletToken.Tokenid)
	if boo != true || err != nil {
		log.Error("walletToken.Tokenid not find token",walletToken.Tokenid,boo,err)
		return false, err
	}

	deci, err := tokensModel.GetDecimal(walletToken.Tokenid)
	if err != nil {
		log.Error("GetDecimal error",err)
		return false,err
	}

	//把总量转成本系统使用的十进制保留八位小数的整数部分
	temp, _ := new(big.Int).SetString(value[2:],16)
	value = decimal.NewFromBigInt(temp, int32(8 - tokensModel.Decimal)).String()

	log.Info("find tokenid：",walletToken.Uid,contract,tokensModel.Id, tokensModel.Mark,tokensModel,",",value)

	var opt int = 1  //充币
	_,err = p.TxModel.Insert(txhash, from, to, value, contract, chainid, walletToken.Uid, tokensModel.Id, tokensModel.Mark,opt)
	if err != nil {
		log.Info("insert into tx order error:",err)
		log.WithFields(log.Fields{
			"txhash":txhash,
			"from":from,
			"to":to,
			"value":value,
			"contract":contract,
			"chainid":chainid,
			"uid":walletToken.Uid,
			"id":tokensModel.Id,
			"mark":tokensModel.Mark,
			"opt":opt,
		}).Info("insert into tx order error:",err)
		return false,err
	}

	_,err = p.TokenInoutModel.Insert(txhash, from, to, value, contract, chainid, walletToken.Uid, tokensModel.Id, tokensModel.Mark, deci,opt)
	if err != nil {
		log.Info("insert into inout order error:",err)
		log.WithFields(log.Fields{
			"txhash":txhash,
			"from":from,
			"to":to,
			"value":value,
			"contract":contract,
			"chainid":chainid,
			"uid":walletToken.Uid,
			"id":tokensModel.Id,
			"mark":tokensModel.Mark,
			"deci":deci,
			"opt":opt,
		}).Info("insert into inout order error:",err)
		return false,nil
	}

	//添加用户token
	//intValue := decimal.NewFromBigInt(temp, int32(8 - p.TokenModel.Decimal)).IntPart()
	fmt.Println("添加用户token：00000000000000000000:",to,",",walletToken.Tokenid,",",value,",",txhash)
	boo,errr := new(Common).AddETHTokenNum(walletToken.Uid,to,walletToken.Tokenid,value,txhash)
	if boo != true {
		log.Error("AddETHTokenNum err:",errr)
	}

	log.WithFields(log.Fields{
		"uid":uid,
		"real_uid":walletToken.Uid,
		"from":from,
		"to":to,
		"chainid":chainid,
		"contract":contract,
		"value":value,
		"txhash":txhash,
	}).Info("add chain complete")

	return true, nil
}

//格式化gas和gasPrice
func (p *EthCBiWatch) formatGas(gas string,gasPrice string) (int64,int64,int64) {
	gas_num, err := strconv.ParseInt(gas, 0, 64)
	if err != nil {
		fmt.Println(err)
		gas_num = 0
	}
	gas_price, err := strconv.ParseInt(gasPrice, 0, 64)
	if err != nil {
		fmt.Println(err)
		gas_price = 0
	}

	t1 := decimal.NewFromFloat(float64(gas_num))
	t1_c := decimal.NewFromFloat(float64(gas_price))
	real_fee := t1.Mul(t1_c).IntPart()
	return gas_num,gas_price,real_fee
}

//写入所有wallet_token到redis中
func (p *EthCBiWatch) WriteAllWalletTokenToRedis() {
	walletToken := new(models.WalletToken)
	err,data := walletToken.GetAllAddress()
	if err != nil {
		log.Error("GetAllAddress error",err)
		return
	}
	for _,v := range data {
		key := ETH_CBI_ADDRESS_REDIS_KEY
		field := "%s:%s"  //chainid:address:contract
		field = fmt.Sprintf(field,v.Address,v.Contract)
		field = strings.ToLower(field)
		err = utils.Redis.HSet(key,field,strings.ToLower(v.Address)).Err()
		if err != nil {
			log.Error("redis error ",err)
		}
		//修改时间
		p.GetWalletTokenLastTime = v.CreatedTime
		//测试查询数据
		query := utils.Redis.HGet(key,field)
		log.Info("redis_result:",key,field,query.Err(),query.Val())
	}
}

//增量更新写入wallet_token到redis中
func (p *EthCBiWatch) WriteIncrWalletTokenToRedis() {
	createdTime := p.GetWalletTokenLastTime.Format("2006-01-02 15:04:05")
	walletToken := new(models.WalletToken)
	err,data := walletToken.GetAddressByTime(createdTime)
	if err != nil {
		log.Error("WriteIncrWalletTokenToRedis error",err)
		return
	}
	for _,v := range data {
		key := ETH_CBI_ADDRESS_REDIS_KEY
		field := "%s:%s"  //chainid:address:contract
		field = fmt.Sprintf(field,v.Address,v.Contract)
		field = strings.ToLower(field)
		err := utils.Redis.HSet(key,field,strings.ToLower(v.Address)).Err()
		if err != nil {
			log.Error("WriteIncrWalletTokenToRedis hset error:",err)
		}
		//修改时间
		p.GetWalletTokenLastTime = v.CreatedTime
	}
}

//判断地址是否存在
func (p *EthCBiWatch) ExistsAddress(address string,chainid int,contract string) bool {
	key := ETH_CBI_ADDRESS_REDIS_KEY
	field := "%s:%s"  //chainid:address:contract
	field = fmt.Sprintf(field,address,contract)
	field = strings.ToLower(field)
	return utils.Redis.HExists(key,field).Val()
}

