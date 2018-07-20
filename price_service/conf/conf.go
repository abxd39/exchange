package conf

import (
	"flag"
	"github.com/Unknwon/goconfig"
	"github.com/micro/go-plugins/registry/consul"
)

var (
	confPath string
	Cfg      *goconfig.ConfigFile
)

func NewConfig(path string) *goconfig.ConfigFile {
	ConfigFile, err := goconfig.LoadConfigFile(path)
	if err != nil {
		panic("load config err is " + err.Error())
		return nil
	}
	return ConfigFile
}

func init() {
	flag.StringVar(&confPath, "conf", "price.ini", "config path")
}

func Init() {
	Cfg = NewConfig(confPath)
}
