package log

import (
	"github.com/sirupsen/logrus"
	"os"
	cf "digicon/ws_service/conf"
)

var Log *logrus.Logger

func InitLog(){
	Log = logrus.New()
	filename := cf.Cfg.MustValue("log", "log_path")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = file
	} else {
		panic("Failed to log to file, using default stderr")
	}
}

