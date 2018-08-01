package utils

import (
	"github.com/Unknwon/goconfig"
	"path/filepath"
	"sync"
	"fmt"
)

var (
	//火币网api配置
	HuoBiApi map[string]string
	Api map[string]string
	UserData map[string]map[string]string
)

var lock1 *sync.Mutex
var lock2 *sync.Mutex

func init() {
	HuoBiApi = make(map[string]string)
	HuoBiApi = GetHuoBiApiConfig()
	Api = make(map[string]string)
	Api = GetApiConfig()
	UserData = make(map[string]map[string]string)
	UserData["sell_user"] = make(map[string]string)
	UserData["sell_user"] = GetUserConfig("sell_user")
	UserData["buy_user"] = make(map[string]string)
	UserData["buy_user"] = GetUserConfig("buy_user")
}

//读取配置
func GetHuoBiApiConfig() map[string]string {
	defer PanicRecover()
	configfile,_ := filepath.Abs("./src/exc_order/utils/config.ini")

	cfg, err := goconfig.LoadConfigFile(configfile)
	if err != nil {
		panic("读取配置文件失败[config.ini]")
		return map[string]string{}
	}
	section,err := cfg.GetSection("huobi")
	if err != nil {
		panic("读取urls失败")
	}
	return section
}

//读取配置
func GetApiConfig() map[string]string {
	defer PanicRecover()
	configfile,_ := filepath.Abs("./src/exc_order/utils/config.ini")

	cfg, err := goconfig.LoadConfigFile(configfile)
	if err != nil {
		panic("读取配置文件失败[config.ini]")
		return map[string]string{}
	}
	section,err := cfg.GetSection("api")
	if err != nil {
		panic("读取urls失败")
	}
	return section
}

//读取配置
func GetUserConfig(section_name string) map[string]string {
	defer PanicRecover()
	configfile,_ := filepath.Abs("./src/exc_order/utils/config.ini")

	cfg, err := goconfig.LoadConfigFile(configfile)
	if err != nil {
		panic("读取配置文件失败[config.ini]")
		return map[string]string{}
	}
	section,err := cfg.GetSection(section_name)
	if err != nil {
		panic("读取urls失败")
	}
	return section
}

//读取配置
func GetHuoBiTradeApiUrl(key,symbol string) string {
	defer PanicRecover()
	//lock1.Lock()
	//defer lock1.Unlock()
	if _,ok := HuoBiApi[key];!ok {
		return ""
	}
	return fmt.Sprintf(HuoBiApi[key],symbol)
}

//读取配置
func GetApiUrl(key string) string {
	defer PanicRecover()
	//lock1.Lock()
	//defer lock1.Unlock()
	if _,ok := Api[key];!ok {
		return ""
	}
	return Api[key]
}

//读取配置
func GetUser(section_name string) map[string]string {
	defer PanicRecover()
	//lock1.Lock()
	//defer lock1.Unlock()
	if _,ok := UserData[section_name];!ok {
		return map[string]string{}
	}
	return UserData[section_name]
}

func GetGoConfigP() (bool,*goconfig.ConfigFile) {
	defer PanicRecover()
	configfile,_ := filepath.Abs("./src/exc_order/utils/config.ini")

	cfg, err := goconfig.LoadConfigFile(configfile)
	if err != nil {
		return false,cfg
	}
	return true,cfg
}