package conf

import (
	"flag"
	"github.com/Unknwon/goconfig"
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
	flag.StringVar(&confPath, "conf", "websocket.ini", "config path")
}

func Init() {
	Cfg = NewConfig(confPath)
}
