package watch

import (
	"fmt"
	"strings"
	."digicon/wallet_service/model"
	"errors"
	"digicon/common/convert"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	cf "digicon/wallet_service/conf"
	"github.com/tidwall/gjson"
	"strconv"
	log "github.com/sirupsen/logrus"
)

//通知
type Notice struct{}

func NewNotice() *Notice {
	return &Notice{}
}

//比特币充币提醒
func (p *Notice) BTCCBiNotice(data TranItem) {
	walletToken := new(WalletToken)
	err := walletToken.GetByAddress(data.Address)
	if err != nil || walletToken.Uid <= 0 {
		return
	}
	tokens := new(Tokens)
	boo,err := tokens.GetByid(walletToken.Tokenid)
	if boo != true || err != nil {
		return
	}
	amount :=  strconv.FormatFloat(data.Amount, 'E', -1, 64)
	p.CBiNotice(walletToken.Uid,amount,tokens.Mark)
}

//提币提醒
func (p *Notice) TiBiNoticeByTxHash(txhash string) (err error) {
	defer func() {
		if err != nil {
			log.Error("比特币提币通知失败：",txhash,err)
		}
	}()
	tokenInout := new(TokenInout)
	err = tokenInout.GetByHash(txhash)
	if err != nil {
		return
	}
	p.TiBiCompleteSendSms(tokenInout.Id)
	return
}

//eth和erc20代币充币提醒
func (p *Notice) ETHERC20CBiNotice(uid int,tokenId int,tokenName string,amount string) {
	p.CBiNotice(uid,amount,tokenName)
}

//usdt代币充币提醒
func (p *Notice) USDTCBiNotice(uid int,tokenId int,tokenName string,amount string) {
	p.CBiNotice(uid,amount,tokenName)
}

//充币到账短信通知
func (p *Notice) CBiNotice(uid int,amount string,tokenName string) (err error) {
	user := new(User)
	boo,err := user.GetUser(uint64(uid))
	if err != nil {
		return
	}
	if boo != true {
		err = errors.New("用户数据为空")
		return
	}

	gateway_ip := cf.Cfg.MustValue("hosts","gateway_ip","")
	if gateway_ip == "" {
		return
	}
	url := gateway_ip + "/user/send_notice"

	postData := make(map[string]interface{})
	postData["phone"] = user.Phone
	mark := tokenName
	num := amount
	postData["content"] = strings.Join([]string{"充币已经完成，币种：",mark,"，到账数量：",num},"")
	postData["auth"] = p.GetAuth()

	result,err := p.RpcPost(url,postData)
	if err != nil {
		return err
	}
	if res := gjson.Get(string(result),"code").Int();res != 0 {
		return errors.New(gjson.Get(string(result),"msg").String())
	}

	log.Info("TiBiCompleteSendSms complete")
	return

}

//提币完成短信通知
func (p *Notice) TiBiCompleteSendSms(apply_id int) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"apply_id":apply_id,
				"err":err,
			}).Error("TiBiCompleteSendSms error")
		}
		fmt.Println("结果：",err,apply_id)
	}()

	var boo bool

	tokenInout := new(TokenInout)
	err = tokenInout.GetByApplyId(apply_id)
	if err != nil {
		return
	}
	tokens := new(Tokens)
	boo,err = tokens.GetByid(tokenInout.Tokenid)
	if err != nil {
		return
	}
	if boo != true {
		err = errors.New("token not found!")
		return
	}

	user := new(User)
	boo,err = user.GetUser(uint64(tokenInout.Uid))
	if err != nil {
		return
	}
	if boo != true {
		err = errors.New("用户数据为空")
		return
	}

	gateway_ip := cf.Cfg.MustValue("hosts","gateway_ip","")
	if gateway_ip == "" {
		return
	}
	url := gateway_ip + "/user/send_notice"

	postData := make(map[string]interface{})
	postData["phone"] = user.Phone
	mark := tokens.Mark
	num := convert.Int64ToStringBy8Bit(tokenInout.Amount)
	postData["content"] = strings.Join([]string{"你申请的提币已经完成，币种：",mark,"，到账数量：",num},"")
	postData["auth"] = p.GetAuth()

	result,err := p.RpcPost(url,postData)
	if err != nil {
		return err
	}
	if res := gjson.Get(string(result),"code").Int();res != 0 {
		return errors.New(gjson.Get(string(result),"msg").String())
	}

	log.Info("TiBiCompleteSendSms complete")
	return
}

//获取auth
func (p *Notice) GetAuth() string {
	return "3b588ad9403a9f8356e1d8639153eb89"
}

func (p *Notice) RpcPost(url string, send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println("rpc post:", err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	//fmt.Println("resp:", resp)
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
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println(*str)
}
