package log

import (
	cf "digicon/token_service/conf"
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLog() {
	Log = logrus.New()

	filename := cf.Cfg.MustValue("log", "log_path")

	_, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = os.Stdout
		//Log.Out = file
	} else {
		panic("Failed to log to file, using default stderr")
	}
}
