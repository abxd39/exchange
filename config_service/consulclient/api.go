package consulclient

import (
	"digicon/config_service/dao"
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Client struct {
	c *api.Client
}

func NewClient() *Client {
	c := dao.DB.GetConsulCli()
	return &Client{c: c}
}

/*
	设置单个值
*/
func (this *Client) Put(k string, value []byte) (err error) {
	kv := this.c.KV()
	p := &api.KVPair{Key: k, Value: value}
	wm, err := kv.Put(p, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(wm)
	return
}

/*
	result 结果根据对应的struct数据结构来解析
*/
func (this *Client) List(key string) (result api.KVPairs, pair *api.QueryMeta, err error) {
	kv := this.c.KV()
	result, pair, err = kv.List(key, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

/*
	删除单个键值
*/
func (this *Client) Delete(key string) (wmeta *api.WriteMeta, err error) {
	kv := this.c.KV()
	wmeta, err = kv.Delete(key, nil)
	return
}
