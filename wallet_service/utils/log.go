package utils

import (

	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	println("log 初始化")
	filename := Cfg.MustValue("log", "log_path")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = file
	} else {
		panic("Failed to log to file, using default stderr")
	}
}
