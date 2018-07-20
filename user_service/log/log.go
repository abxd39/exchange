package log

import (
	cf "digicon/user_service/conf"
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLog() {
	Log = logrus.New()
	s := cf.Cfg.MustValue("log", "switch")
	if s == "1" {
		filename := cf.Cfg.MustValue("log", "log_path")

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			Log.Out = file
		} else {
			Log.Out = os.Stdout
		}
	} else {
		Log.Out = os.Stdout
	}
}