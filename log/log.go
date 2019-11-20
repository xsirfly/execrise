package log

import (
	"errors"
	"exercise/conf"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"time"
	"os"
)

func Init() error {
	if conf.GetConf() == nil {
		return errors.New("config has not init")
	}
	logConf := conf.GetConf().Log
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	level, err := logrus.ParseLevel(logConf.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	if logConf.Output == "stdout" {
		logrus.SetOutput(os.Stdout)
		return nil
	} else {
		writer, err := rotatelogs.New(
			logConf.Output+".%Y%m%d%H",
			rotatelogs.WithLinkName(logConf.Output),   // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			return err
		}

		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},&logrus.TextFormatter{DisableColors: true})
		logrus.AddHook(lfHook)
	}
	return nil
}
