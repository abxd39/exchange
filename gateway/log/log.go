package log

import (
	cf "digicon/gateway/conf"
	"github.com/liudng/godump"
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLog() {
	Log = logrus.New()

	filename := cf.Cfg.MustValue("log", "log_path")
	godump.Dump(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Out = os.Stdout
		//panic("Failed to log to file, using default stderr")
	}
}
