package utils

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Logger

func NewLogger() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)

	//同时写到多个输出
	w1 := os.Stdout                                                               //写到控制台
	w2, _ := os.OpenFile("./like.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //写到文件赋予读写追加
	Logger.SetOutput(io.MultiWriter(w1, w2))
}
