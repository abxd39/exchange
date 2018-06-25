package conf

import (
	"flag"
	"github.com/Unknwon/goconfig"
)

var (
	confPath   string
	Cfg        *goconfig.ConfigFile
	SmsAccount string //短信平台账号
	SmsPwd     string //短信平台密码
	SmsWebUrl  string //短信平台网址
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
	flag.StringVar(&confPath, "conf", "user.ini", "config path")
}

func Init() {
	Cfg = NewConfig(confPath)
	SmsAccount = Cfg.MustValue("sms", "account", "N2562426")
	SmsPwd = Cfg.MustValue("sms", "pwd", "rSLFN2Io96772f")
	SmsWebUrl = Cfg.MustValue("sms", "sms_url", "http://smssh1.253.com/msg/send/json")

}
