package client

import (
	"github.com/ouqiang/timewheel"
	"time"
	"fmt"
	"log"
	"digicon/wallet_service/model"
	"digicon/wallet_service/utils"
	"github.com/tidwall/gjson"
	"math/big"
	"github.com/shopspring/decimal"
	"encoding/json"
	. "digicon/wallet_service/model"
	"unsafe"
	"strconv"
	"bytes"
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
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
		log.Println("init error",err)
	}
	if !exists {
		log.Println("token not exists btc ...")
	}
	p.Url = tokenModel.Node

	//初始化同步区块时间轮
	p.ethCheckTranNewTW = timewheel.New(1 * time.Second, 3600, func(data timewheel.TaskData) {
		fmt.Println("start eth.check.watch.new...")
		//处理交易验证
		p.checkTransactionDeal()
		//继续添加定时器
		p.ethCheckTranNewTW.AddTimer(ETH_UPDATE_INTERVAL_TW * time.Second, "eth_check_tibi", timewheel.TaskData{})
	})
	p.ethCheckTranNewTW.Start()
	//开始一个事件处理
	p.ethCheckTranNewTW.AddTimer(ETH_UPDATE_INTERVAL_TW * time.Second, "eth_check_tibi", timewheel.TaskData{})
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
	temp, _ := new(big.Int).SetString(status[2:], 16)
	statuss := decimal.NewFromBigInt(temp, 0).IntPart()
	if statuss == 1 {
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
			return
		}
		//新增数据
		p.WriteEthInRecord(data)
		//修改提币申请订单
		new(models.TokenInout).BteUpdateAppleDone(txhash)
		return
	}
	//暂未打包成功，重新放入队列，等待下次执行
	p.PushRedisList(txhash)
}

//保存数据到redis队列
func (p *EthTiBiWatch) PushRedisList(txhash string) {
	fmt.Println("收到一个eth交易检查请求：",txhash)
	redis := utils.Redis
	err := redis.RPush(ETH_CHECK_LIST_KEY,txhash).Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}

//从redis队列读取数据
func (p *EthTiBiWatch) GetDataFromRedis() (error,string) {
	redis := utils.Redis
	query := redis.LPop(ETH_CHECK_LIST_KEY)
	if query.Err() != nil {
		fmt.Println(query.Err().Error())
		return query.Err(),""
	}
	data := query.Val()
	return nil,data
}

//写入一条订单数据到表：token_chain_inout
func (p *EthTiBiWatch) WriteEthInRecord(data TranInfo) {
	tmp1,_ := new(big.Int).SetString(data.Value,10)
	value := decimal.NewFromBigInt(tmp1, int32(8)).IntPart()

	var inOutToken = new(models.TokenInout)

	var walletToken = new(models.WalletToken)
	err := walletToken.GetByAddress(data.From)
	if err != nil {
		log.Println("WriteBtcInRecord address not exists",err.Error())
		return
	}

	inOutToken.Id = 0
	inOutToken.Txhash = data.Hash
	inOutToken.From = data.From
	inOutToken.To = data.To
	//inOutToken.Value = value
	inOutToken.Amount = value
	inOutToken.Tokenid = 1  //提币
	inOutToken.TokenName = "ETH"
	inOutToken.Uid = walletToken.Uid
	inOutToken.Tokenid = walletToken.Tokenid
	affected, err := utils.Engine_wallet.InsertOne(inOutToken)
	if err != nil {
		log.Println("WriteEthInRecord error",err.Error())
	}
	fmt.Println(affected)
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
	Url string //节点链接
	Chainid     int
	BlockNumber int //当前处理到的高度

	WalletTokenModel *WalletToken     //钱包详情
	TxModel          *TokenChainInout //链上交易记录
	TokenInoutModel  *TokenInout      //平台交易记录
	TokenModel       *Tokens          //币种类
	ContextModel     *Context         //处理上下文
}

const (
	ETH_CBI_INTERVAL_TW = 5 //时间轮定时器间隔时间
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
		fmt.Println("start eth.cbi.watch.new...")
		//区块操作处理
		p.WorkerDone()
		//继续添加定时器
		p.ethCheckCBTranNewTW.AddTimer(ETH_CBI_INTERVAL_TW * time.Second, "eth_check_cbi", timewheel.TaskData{})
	})
	p.ethCheckCBTranNewTW.Start()
	//开始一个事件处理
	p.ethCheckCBTranNewTW.AddTimer(ETH_CBI_INTERVAL_TW * time.Second, "eth_check_cbi", timewheel.TaskData{})

}

//处理区块
func (p *EthCBiWatch) WorkerDone() {
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

	if p.BlockNumber < hight-6 {
		for i := p.BlockNumber + 1; i <= hight-6; i++ {
			p.WorkerHander(i)
			//记录当前进度
			p.ContextModel.Save(p.Url, p.Chainid, i)

		}
	}
}

//具体处理区块
func (p *EthCBiWatch) WorkerHander(num int) error {
	ret, err := p.GetblockBynumber(num)
	if err != nil {
		return err
	}
	var block map[string]interface{}
	//fmt.Println(string(ret))
	json.Unmarshal(ret, &block)
	txs := block["result"].(map[string]interface{})["transactions"].([]interface{})
	//fmt.Println(txs)
	for i := 0; i < len(txs); i++ {
		tx := txs[i].(map[string]interface{})
		if tx["to"] == nil { //部署合约交易直接跳过
			continue
		}

		//检查eth转账
		ext, err := p.WalletTokenModel.AddrExist(tx["to"].(string), p.Chainid, "")

		if err != nil {
			fmt.Println(err)
			return err
		}
		if ext {
			fmt.Println("发现一个eth转账")
			//TODO:
			p.newOrder(p.WalletTokenModel.Uid, tx["from"].(string), tx["to"].(string), p.Chainid, "", tx["value"].(string), tx["hash"].(string))

			continue
		}

		input := tx["input"].(string)
		//不是token转账跳过
		if strings.Count(input, "") < 138 || strings.Compare(input[0:10], "0xa9059cbb") != 0 {
			//fmt.Println(input)
			continue
		}

		ext, err = p.WalletTokenModel.AddrExist(fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string))
		//fmt.Println(ext,err,this.WalletTokenModel)
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
		fmt.Println("发现一个token转账")

		ok, err := p.newOrder(p.WalletTokenModel.Uid, tx["from"].(string), fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string), fmt.Sprintf("0x%s", input[vstart:138]), tx["hash"].(string))
		fmt.Println(ok, err)
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
	fmt.Println(num, fmt.Sprintf("0x%s", strconv.FormatInt(int64(num), 16)))
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
	str := (*string)(unsafe.Pointer(&rsp))
	fmt.Println(*str)
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
	str := (*string)(unsafe.Pointer(&rsp))
	fmt.Println(*str)
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
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	url := p.Url
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//byte数组直接转成string，优化内存
	return respBytes, nil
}

//新增充值订单
func (p *EthCBiWatch) newOrder(uid int, from string, to string, chainid int, contract string, value string, txhash string) (bool,error) {
	//交易是否已经收录
	exist, err := p.TxModel.TxhashExist(txhash, p.Chainid)

	if err != nil {
		return false, err
	}
	if exist {
		return false, errors.New("tx already exists")
	}

	//
	tokenid, err := p.TokenModel.GetidByContract(contract, p.Chainid)
	deci, _ := p.TokenModel.GetDecimal(tokenid)
	fmt.Println(tokenid, err)
	if err != nil {
		return false, err
	}
	if tokenid == 0 {
		return false, errors.New("token not exist")
	}
	fmt.Println("this.TxModel.Insert")
	var opt int = 2  //充币
	p.TxModel.Insert(txhash, from, to, value, contract, chainid, uid, p.TokenModel.Id, p.TokenModel.Mark,opt)

	p.TokenInoutModel.Insert(txhash, from, to, value, contract, chainid, uid, p.TokenModel.Id, p.TokenModel.Mark, deci,opt)
	return true, nil
}

