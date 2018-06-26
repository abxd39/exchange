package conf

import "github.com/Unknwon/goconfig"
import "flag"

var (
	Inifpath string
	Inif     *goconfig.ConfigFile
)

func init() {
	flag.StringVar(&Inifpath, "conf", "backstage.ini", "config path")
	if l := len(Inifpath); l == 0 {
		panic("数据库配置文件路径获取失败")
	}
	NewConfi(Inifpath)
}

func NewConfi(path string) *goconfig.ConfigFile {
	ini, err := goconfig.LoadConfigFile(path)
	if err != nil {
		panic(err.Error())
	}
	return ini
}
