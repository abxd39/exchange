package watch

import (
	"encoding/json"
	"strings"
	"fmt"
	log "github.com/sirupsen/logrus"
	. "digicon/wallet_service/model"
	"digicon/common/errors"
	"strconv"
)

//eth手动入账
type EthOperate struct {
	Url string
	Chainid int
	WalletTokenModel *WalletToken     //钱包详情
	cbiP             *EthCBiWatch
}

func (p *EthOperate) init() (bool,error) {
	//初始化
	p.WalletTokenModel = new(WalletToken)
	//查询ETH节点
	var data = new(Tokens)
	bool, er := data.GetByName("ETH")
	if bool != true || er != nil {
		return false,errors.New("查询token失败")
	}

	p.Url = data.Node
	p.Chainid = data.Chainid

	p.cbiP = new(EthCBiWatch)
	p.cbiP.Url = data.Node
	p.cbiP.Chainid = data.Chainid

	p.cbiP.WalletTokenModel = new(WalletToken)
	p.cbiP.TxModel = new(TokenChainInout)
	p.cbiP.TokenInoutModel = new(TokenInout)
	p.cbiP.TokenModel = new(Tokens)
	p.cbiP.ContextModel = new(Context)


	return true,nil
}

//具体处理区块
func (p *EthOperate) WorkerHander(num int) (error,string) {
	//初始化
	boo,err := p.init()
	fmt.Println("node:",p.Url,p.Chainid)
	if boo != true || err != nil {
		return err,"初始化失败"
	}
	//log.Info("start WorkerHander",num)
	ret, err := p.cbiP.GetblockBynumber(num)
	if err != nil {
		return err,""
	}
	var block map[string]interface{}
	json.Unmarshal(ret, &block)
	txs := block["result"].(map[string]interface{})["transactions"].([]interface{})


	syncEthNum := 0
	syncTokenNum := 0

	for i := 0; i < len(txs); i++ {
		tx := txs[i].(map[string]interface{})
		if tx["to"] == nil { //部署合约交易直接跳过
			continue
		}

		//检查eth转账
		ext := p.cbiP.ExistsAddress(tx["to"].(string), p.Chainid, "")

		if ext {
			log.Info("find_a_eth")
			syncEthNum++
			//TODO:
			p.cbiP.newOrder(p.WalletTokenModel.Uid, tx["from"].(string), tx["to"].(string), p.Chainid, "", tx["value"].(string), tx["hash"].(string),tx["gas"].(string),tx["gasPrice"].(string))

			continue
		}

		input := tx["input"].(string)
		//不是token转账跳过
		if strings.Count(input, "") < 138 || strings.Compare(input[0:10], "0xa9059cbb") != 0 {
			continue
		}

		ext = p.cbiP.ExistsAddress(fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string))
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
		fmt.Println("find_a_token")

		p.cbiP.newOrder(p.WalletTokenModel.Uid, tx["from"].(string), fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string), fmt.Sprintf("0x%s", input[vstart:138]), tx["hash"].(string),tx["gas"].(string),tx["gasPrice"].(string))
		syncTokenNum++
		continue

	}
	return nil,strings.Join([]string{"同步完成：",strconv.Itoa(syncEthNum),"条eth交易,",strconv.Itoa(syncTokenNum),"条token交易"},"")
}
