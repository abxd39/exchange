package conf

import (
	"github.com/Unknwon/goconfig"
	"github.com/golang/glog"
)

var Cfg *goconfig.ConfigFile

func NewConfig() *goconfig.ConfigFile {
	ConfigFile, err := goconfig.LoadConfigFile("wallet.ini")
	if err != nil {
		glog.Fatalf("load config err is %s", err.Error())
		return nil
	}
	return ConfigFile
}

func InitConf() {
	Cfg = NewConfig()
}
