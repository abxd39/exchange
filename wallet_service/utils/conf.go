package utils

import (
	"flag"
	"github.com/Unknwon/goconfig"
)

var Cfg *goconfig.ConfigFile

func NewConfig(path string) *goconfig.ConfigFile {
	ConfigFile, err := goconfig.LoadConfigFile(path)
	if err != nil {
		panic("load config err is " + err.Error())
		return nil
	}
	return ConfigFile
}

func init() {
	println("conf 初始化")
	var confPath string
	flag.StringVar(&confPath, "conf", "wallet.ini", "config path")
	Cfg = NewConfig(confPath)
}
