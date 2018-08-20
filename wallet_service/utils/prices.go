package utils

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"digicon/common/errors"
)

//币种和汇率
//人民币价格获取
//market/cny_prices
type Price struct {
	CnyPrice string `json:"cny_price"`
	UsdPrice string `json:"usd_price"`
	TokenId int`json:"token_id"`
	CnyPriceInt int `json:"cny_price_int"`
}
type  temp struct {
	List []Price `json:"list"`
}
func GetTokenCnyPriceList(tid []int)([]Price,error)  {
	params := make(map[string]interface{})
	params["token_id"] = tid
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil,err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", "http://47.75.89.78:8069/market/cny_prices?",reader)
	if err != nil {
		return nil,err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return nil,err
	}
	body,err:=ioutil.ReadAll(result.Body)
	type base struct {
		Code int   `json:"code"`
		Data temp `json:"data"`
		Msg string     `json:"msg"`
	}
	b:=new(base)
	err= json.Unmarshal(body,b)
	if err!=nil{
		return nil,err
	}
	if b.Code!=0{
		return nil,errors.New(b.Msg)
	}
	return b.Data.List,nil
}

func GetCnyPrice(tokenId int) (error,int) {
	data,err := GetTokenCnyPriceList([]int{tokenId})
	if err != nil {
		return err,0
	}
	if len(data) == 0 {
		return errors.New("token price not find"),0
	}
	return nil,data[0].CnyPriceInt
}