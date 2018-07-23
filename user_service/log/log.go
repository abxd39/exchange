package log

import (
	"github.com/sirupsen/logrus"
	"os"

	cf "digicon/price_service/conf"
)

var Log *logrus.Logger
/*
func InitLogger() {

	path := cf.Cfg.MustValue("log", "log_path")
	filename := cf.Cfg.MustValue("log", "log_path")
	baseLogPath := fmt.Sprintf(path,filename)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic("config local file system logger error. ")
	}

	//log.SetFormatter(&log.TextFormatter{})
	level := cf.Cfg.MustValue("log", "log_path")
	switch level  {

	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	case "info":
		setNull()
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		setNull()
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		setNull()
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		setNull()
		logrus.SetLevel(logrus.InfoLevel)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	})
	logrus.AddHook(lfHook)
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err!= nil{
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logrus.SetOutput(writer)
}

*/
func InitLog() {
	Log = logrus.New()
	s := cf.Cfg.MustValue("log", "switch")
	if s == "1" {
		filename := cf.Cfg.MustValue("log", "log_path")

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			Log.Out = file
		} else {
			Log.Out = os.Stdout
		}
	} else {
		Log.Out = os.Stdout
	}
}
