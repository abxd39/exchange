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
	flag.StringVar(&confPath, "conf", "price.ini", "config path")
}

func Init() {
	Cfg = NewConfig(confPath)
/*
	viper.SetConfigName("price") // name of config file (without extension)
	//viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	//viper.AddConfigPath("./price_service")   // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	g:=viper.Get("etcd.addr")
	fmt.Println(g)
	*/
}
