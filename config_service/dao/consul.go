package dao

import (
	"github.com/hashicorp/consul/api"
	"digicon/config_service/conf"
	. "digicon/config_service/log"
)




type ConsulCli struct {
	ccon *api.Client
}


func NewConsulCli() *ConsulCli{
	addr := conf.Cfg.MustValue("consul", "addr")
	token := conf.Cfg.MustValue("consul", "token")

	client, err := api.NewClient(&api.Config{
		Token    :token,
		Address  :addr,
	})
	if err != nil {
		Log.Fatal("new consul client error!")
	}

	return &ConsulCli{
		ccon: client,
	}
}

func (s *Dao) GetConsulCli() *api.Client{
	return s.consul.ccon
}