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

	EmailAppKey    string
	EmailSecretKey string

	GtPrivateKey  string
	GtCaptchaID  string

	SmsTitle string
	MailTitle string
	Subject	string
	Alias string
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
	SmsAccount = Cfg.MustValue("sms", "account", "I1757342")
	SmsPwd = Cfg.MustValue("sms", "pwd", "i1PYZXVaWt2de6")
	//SmsWebUrl = Cfg.MustValue("sms", "sms_url", "http://smssh1.253.com/msg/send/json")
	SmsWebUrl = Cfg.MustValue("sms", "sms_url", "https://intapi.253.com")

	EmailAppKey = Cfg.MustValue("email", "app_key", "LTAIcJgRedhxruPq")
	EmailSecretKey = Cfg.MustValue("email", "secret_key", "d7p6tWRfy0B2QaRXk7q4mb5seLROtb")

	GtPrivateKey = Cfg.MustValue("gree", "key", "668d6d27cb1186d138eb9b225436e4b9")
	GtCaptchaID = Cfg.MustValue("gree", "id", "73909f4a67161216debdcb3de16ef6c5")

	SmsTitle = Cfg.MustValue("sms", "title", "【UNT】")
	MailTitle = Cfg.MustValue("mail", "title", "您好，您正在注册UNT账号。【UNT】安全验证:")
	Subject = Cfg.MustValue("mail", "subject", "欢迎注册UNT")
	Alias = Cfg.MustValue("mail", "alias", "shendun")
}
