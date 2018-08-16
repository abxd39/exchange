package xlog

import (
	"bufio"
	"digicon/common/hook"
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func InitLogger(path, name, level string) {
	baseLogPath := fmt.Sprintf("%s%s", path, name)
	writer, err := rotatelogs.New(
		baseLogPath+"_%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		panic("config local file system logger error. ")
	}


	switch level {

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

	g := &logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, g)
	logrus.AddHook(lfHook)
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logrus.SetOutput(writer)
}
