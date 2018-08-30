package gather

import(
	."digicon/wallet_service/model"
	"sync"
	"fmt"
	cf "digicon/wallet_service/conf"
	"digicon/wallet_service/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
	"github.com/shopspring/decimal"
	"encoding/hex"
	"errors"
)

//eth和erc20代币归总

const(
	NUM_LIMIT = 10  //每次读取十个
	ETH_NUM_LIMIT = 60000
)

type EthGather struct{
	Tokens map[int]Tokens    //map[id]Tokens  存储所有token和归总限制数量
	lock sync.Mutex         //tokens操作锁
	walletTokenTotal int64    //钱包token总数
	maxId int               //当前查询的最大id，使用这个作为查询limit，当limit==uidTotal的时候，需要重新初始化
}

func StartEthGather() {
	//初始化
	p := new(EthGather)
	p.init()
}

func (p *EthGather) init() {
	//初始化所有token
	p.initTokens()
	//初始化总数
	p.initTotal()
	//设置maxId
	p.setMaxId(0)
}

func (p *EthGather) initTokens() {
	p.Tokens = make(map[int]Tokens)
	tokens := new(Tokens)
	data,err := tokens.ListTokens()
	if err != nil {
		panic("查询币种错误："+err.Error())
	}
	for _,v := range data {
		p.Tokens[v.Id] = v
	}
}

func (p *EthGather) initTotal() {
	walletToken := new(WalletToken)
	num,err := walletToken.GetCount()
	if err != nil {
		fmt.Println("获取总数失败：",err)
	}
	p.walletTokenTotal = num
}

func (p *EthGather) setMaxId(id int) {
	p.maxId = id
}

func (p *EthGather) exec() {
	//查询数据
	walletToken := new(WalletToken)
	data,err := walletToken.GetWalletToken(p.maxId + 1,NUM_LIMIT)
	if err != nil {
		fmt.Println("查询钱包报错：",err)
		return
	}
	if len(data) == 0 {
		//已经不存在未查询的数据，重新从0开始
		//初始化总数
		p.initTotal()
		//设置maxId
		p.setMaxId(0)
		return
	}

	for _,v := range data {
		//检查是否官方账号
		if p.checkGFAddress(v.Address) == true {
			continue
		}
		//检查是否满足归总限制条件
		_,sure := p.checkGatherLimit(v)
		if sure == false {
			continue
		}
		//查询以太币数量
		boo,zhuanNum := p.checkEthNum(v)
		if boo == false {
			continue
		}
		if boo == true && zhuanNum > 0 {
			//需要转点以太币到这个账号
			p.sendETHToAddress(v,zhuanNum)
			continue
		}

		//执行归总操作
		p.sendToGFAddress(v)
	}
}


//检查是否是官方账号
func (p *EthGather) checkGFAddress(address string) bool {
	eth_address := cf.Cfg.MustValue("accounts","eth_address","")
	if address != eth_address {
		return false
	}
	return true
}

//检查是否满足归总条件
func (p *EthGather) checkGatherLimit(walletToken WalletToken) (err error,boo bool) {
	defer func() {
		if boo != true {
			log.WithFields(log.Fields{
				"err":err.Error(),
			}).Error("checkGatherLimit error")
		}
	}()
	boo,token := p.getTokenData(walletToken.Tokenid)
	if boo != true {
		log.Error("查询token数据失败：",walletToken.Tokenid)
		return nil,false
	}
	//查询数量
	num,err := utils.RpcGetValue(token.Node,walletToken.Address,walletToken.Contract,token.Decimal)
	if err != nil {
		log.Error("查询数量失败")
		return nil,false
	}

	int64_num,err := strconv.ParseInt(num,10,64)
	if int64_num > token.Gather_limit {
		return nil,true
	}
	return nil,false
}

//根据token_id获取基本数据
func (p *EthGather) getTokenData(token_id int) (err bool,tokens Tokens) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _,ok := p.Tokens[token_id];ok {
		err = true
		tokens = p.Tokens[token_id]
		return
	}
	return
}

//查询账号以太币数量
func (p *EthGather) checkEthNum(walletToken WalletToken) (boo bool,num float64) {

	tokenModel := new(Tokens)
	exists, err := tokenModel.GetByName("ETH")
	if err != nil {
		log.Info("tokenModel error",err)
		return false,0
	}
	if !exists {
		log.Info("token not exists eth ...")
		return false,0
	}

	boo,token := p.getTokenData(tokenModel.Id)
	if boo != true {
		log.Error("查询token数据失败：",walletToken.Tokenid)
		return false,0
	}
	//查询数量
	number,err := utils.RpcGetValue(token.Node,walletToken.Address,walletToken.Contract,token.Decimal)
	if err != nil {
		log.Error("查询数量失败")
		return false,0
	}
	int64_number,err := strconv.ParseInt(number,10,64)
	if err != nil {
		log.Error("checkEthNum error",err)
		return false,0
	}
	//查询gasPrice
	gasPrice,gasErr := utils.RpcGetGasPrice(tokenModel.Node)
	if gasErr != nil  {
		return false,0
	}
	if ETH_NUM_LIMIT * gasPrice > int64_number  {
		//需要从官方账号转以太币到该账号
		zhuanNum := float64(ETH_NUM_LIMIT * gasPrice - int64_number) * 1.1
		return true,zhuanNum
	}
	return true,0

}

//查询转账预估手续费
func (p *EthGather) getGas() {

}

//转以太币到指定账号
func (p *EthGather) sendETHToAddress(walletToken WalletToken,num float64) error {
	boo,token := p.getTokenData(walletToken.Tokenid)
	if boo != true {
		log.Error("查询token数据失败：",walletToken.Tokenid)
		return nil
	}

	boo,from_uid := p.GetEthAccountUid()
	if boo != true {
		log.Error("查询官方以太坊地址错误")
		return nil
	}

	keystore := &WalletToken{Uid: int(from_uid), Tokenid: int(walletToken.Tokenid)}

	ok, err := utils.Engine_wallet.Where("uid=? and tokenid=?", from_uid, int(walletToken.Tokenid)).Get(keystore)
	if err != nil {
		log.Error("查询账户wallet_token出错：",err)
		return nil
	}
	if ok != true {
		log.Info("没有找到")
		return nil
	}

	//获取随机数
	nonce,nonce_err := utils.RpcGetNonce(token.Node,keystore.Address)

	if nonce_err != nil  {
		log.Error("获取随机数出错：",nonce_err)
		return nil
	}

	//查询gasPrice
	gasPrice,gasErr := utils.RpcGetGasPrice(token.Node)
	if gasErr != nil  {
		return nil
	}

	amount := decimal.NewFromFloat(num).Coefficient()

	signtxstr, err := keystore.Signtx(walletToken.Address, amount, gasPrice,nonce)
	if err != nil {
		log.Error("签名失败：",err)
		return nil
	}
	signtxStr := hex.EncodeToString(signtxstr)
	err,txhash := p.SendTx(token.Node,signtxStr)
	if err != nil {
		log.Error("发送交易报错：",err)
	}
	if txhash == "" {
		log.Error("交易hash为空")
		return nil
	}
	log.Info("转以太币成功：",txhash)
	return nil
}

//查询以太坊官方账户地址
func (p *EthGather) GetEthAccountUid() (bool,int64) {
	from_uid := cf.Cfg.MustInt64("accounts","eth_uid",0)
	if from_uid == 0 {
		return false,0
	}
	return true,from_uid
}

//发送交易
func (p *EthGather) SendTx(node string,signtx string) (error,string) {
	rets, err := utils.RpcSendRawTx(node, signtx)
	if err != nil {
		return err,""
	}
	txhash, ok := rets["result"]
	if ok != true {
		return errors.New("报错了"),""
	}
	return nil,txhash.(string)
}

//转账到官方地址
func (p *EthGather) sendToGFAddress(walletToken WalletToken) error {
	boo,token := p.getTokenData(walletToken.Tokenid)
	if boo != true {
		log.Error("查询token数据失败：",walletToken.Tokenid)
		return nil
	}

	keystore := &WalletToken{Uid: int(walletToken.Uid), Tokenid: int(walletToken.Tokenid)}
	ok, err := utils.Engine_wallet.Where("uid=? and tokenid=?", walletToken.Uid, int(walletToken.Tokenid)).Get(keystore)
	if err != nil {
		log.Error("查询账户wallet_token出错：",err)
		return nil
	}
	if ok != true {
		log.Info("没有找到")
		return nil
	}

	//获取随机数
	nonce,nonce_err := utils.RpcGetNonce(token.Node,keystore.Address)

	if nonce_err != nil  {
		log.Error("获取随机数出错：",nonce_err)
		return nil
	}

	//查询gasPrice
	gasPrice,gasErr := utils.RpcGetGasPrice(token.Node)
	if gasErr != nil  {
		return nil
	}

	//查询以太坊数量
	boo,ethNum := p.getEthNum(walletToken)
	if boo != true {
		log.Error("查询以太坊数量报错")
		return nil
	}

	num := ethNum

	if p.isEth(walletToken.Tokenid) {
		//以太坊转账，不能转太多
		num = ethNum - ETH_NUM_LIMIT * gasPrice
	}

	amount := decimal.New(num,0).Coefficient()

	//执行签名
	signtxstr, err := keystore.Signtx(walletToken.Address, amount, gasPrice,nonce)
	if err != nil {
		log.Error("签名失败：",err)
		return nil
	}
	signtxStr := hex.EncodeToString(signtxstr)
	err,txhash := p.SendTx(token.Node,signtxStr)
	if err != nil {
		log.Error("发送交易报错：",err)
	}
	if txhash == "" {
		log.Error("交易hash为空")
		return nil
	}
	log.Info("归总以太币成功：",txhash)
	return nil

}

//查询账号以太币数量
func (p *EthGather) getEthNum(walletToken WalletToken) (boo bool,num int64) {

	tokenModel := new(Tokens)
	exists, err := tokenModel.GetByName("ETH")
	if err != nil {
		log.Info("tokenModel error",err)
		return false,0
	}
	if !exists {
		log.Info("token not exists eth ...")
		return false,0
	}

	boo,token := p.getTokenData(tokenModel.Id)
	if boo != true {
		log.Error("查询token数据失败：",walletToken.Tokenid)
		return false,0
	}
	//查询数量
	number,err := utils.RpcGetValue(token.Node,walletToken.Address,walletToken.Contract,token.Decimal)
	if err != nil {
		log.Error("查询数量失败")
		return false,0
	}
	int64_number,err := strconv.ParseInt(number,10,64)
	if err != nil {
		log.Error("checkEthNum error",err)
		return false,0
	}
	return true,int64_number
}

//判断是否是以太坊
func (p *EthGather) isEth(token_id int) bool {
	if token_id == 3 {
		return true
	}
	return false
}
