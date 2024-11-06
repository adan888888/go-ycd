package utils

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

/*var Logger *logrus.Logger

func NewLogger() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetReportCaller(true)
	//同时写到多个输出
	w1 := os.Stdout                                                               //写到控制台
	w2, _ := os.OpenFile("./like.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //写到文件赋予读写追加
	Logger.SetOutput(io.MultiWriter(w1, w2))
	Logger.SetFormatter(&logrus.JSONFormatter{}) //以json的方式输出
}*/

// -----第2种 用法----
var Logger *logrus.Entry

func NewLogger() {
	logStore := logrus.New()
	logStore.SetLevel(logrus.DebugLevel) //输出的日志级别
	//同时写到多个输出
	w1 := os.Stdout                                                               //写到控制台
	w2, _ := os.OpenFile("./like.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //写到文件赋予读写追加
	logStore.SetOutput(io.MultiWriter(w1, w2))
	//logStore.SetFormatter(&logrus.JSONFormatter{}) //以json的方式输出
	// 设置输出格式为带颜色的输出
	logStore.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	//logrus.SetFormatter(&CustomFormatter{&logrus.TextFormatter{ForceColors: true}})
	Logger = logStore.WithFields(logrus.Fields{ //提共自定义字段
		//"name": "我爱学习go",
		//"app": "v2v2v2",
	})
	//提供hook函数
	//logStore.AddHook()//出现非常严重问题，直接邮箱或者微信报警，日志分割，塞入一些自定义的字段
	logStore.AddHook(&ColorizedHook{})
	//context
	//logStore.WithContext()
}

//type CustomFormatter struct {
//	*logrus.TextFormatter
//}

/*
	func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
		//例如，
		//Debug级别的日志会显示为青色，
		//Info级别的日志会显示为绿色，
		//Warn级别的日志会显示为黄色，
		//Error级别的日志会显示为红色。
		//其他级别的日志则使用默认颜色。
		//修改日志级别颜色
		switch entry.Level {
		case logrus.DebugLevel:
			return []byte(color.New(color.FgCyan).Sprint(f.TextFormatter.Format(entry))), nil
		case logrus.InfoLevel:
			return []byte(color.New(color.FgGreen).Sprint(f.TextFormatter.Format(entry))), nil
		case logrus.WarnLevel:
			return []byte(color.New(color.FgYellow).Sprint(f.TextFormatter.Format(entry))), nil
		case logrus.ErrorLevel:
			return []byte(color.New(color.FgHiWhite).Sprint(f.TextFormatter.Format(entry))), nil
		default:
			return f.TextFormatter.Format(entry)
		}
	}
*/
type ColorizedHook struct{}

func (hook ColorizedHook) Fire(entry *logrus.Entry) (err error) {
	var levelColor color.Attribute
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = color.FgBlue
	case logrus.WarnLevel:
		levelColor = color.FgYellow
	case logrus.ErrorLevel:
		levelColor = color.BgRed
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.FgMagenta
	default:
		levelColor = color.FgGreen
	}

	message := entry.Message
	if levelColor != 0 {
		message = color.New(levelColor).SprintFunc()(message)
	}

	entry.Message = message
	return nil
}

func (hook ColorizedHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
