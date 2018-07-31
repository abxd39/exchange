package utils


//

//import (
//	"github.com/sirupsen/logrus"
//	"os"
//)
//
//var log *logrus.Logger
//
//func init() {
//	log = logrus.New()
//	println("log 初始化")
//	filename := Cfg.MustValue("log", "log_path")
//	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
//	if err == nil {
//		log.Out = file
//	} else {
//		panic("Failed to log to file, using default stderr")
//	}
//}