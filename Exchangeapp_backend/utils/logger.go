package utils

import (
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
	logStore.SetLevel(logrus.DebugLevel)
	//同时写到多个输出
	w1 := os.Stdout                                                               //写到控制台
	w2, _ := os.OpenFile("./like.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //写到文件赋予读写追加
	logStore.SetOutput(io.MultiWriter(w1, w2))
	logStore.SetFormatter(&logrus.JSONFormatter{}) //以json的方式输出
	// 设置输出格式为带颜色的输出
	logStore.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	Logger = logStore.WithFields(logrus.Fields{ //提共自定义字段
		//"name": "我爱学习go", "app": "v2v2v2",
	})
	//提供hook函数
	//logStore.AddHook()//出现非常严重问题，直接邮箱或者微信报警，日志分割，塞入一些自定义的字段

	//context
	//logStore.WithContext()
}
