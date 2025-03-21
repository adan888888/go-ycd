package main

import (
	"context"
	"exchangeapp/config"
	"exchangeapp/global"
	"exchangeapp/router"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title           Swagger Example API
// @version         1.0
// @contact.name   Like API
// @contact.email  support@swagger.io
// @host      192.168.100.133:3000
// @BasePath
func main() {
	config.InitConfig()
	r := router.SetupRouter()

	port := global.AppConfig.App.Port

	if port == "" {
		port = ":8080"
	}

	//服务器实例
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		//srv.ListenAndServe()启动服务的监听
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { //如果不是服务器的关闭错误
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)   //申明一个类型为 信号 的通道
	signal.Notify(quit, os.Interrupt) //假如发送一个 ctrl+c这样的一个信号 到通道中
	<-quit                            //还没能写入消息的时候，通道入一个阻塞,后面代码不会执行,一旦接收到信号 ，程序就会往下走
	log.Println("Shutdown 服务 ...")    //打印消息服务器正在关闭

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //创建一个带有5秒钟超时 可取消的上下文，5秒内的不管有没有处理好，都停止进程
	defer cancel()                                                          //取消上下文。这里延迟调用，确保函数在返回前取消上下文释放资源
	if err := srv.Shutdown(ctx); err != nil {                               //调用srv.Shutdown(ctx)关闭服务，优雅的退出程序。会一直等待还没有处理完的任务，但是受限于ctx上下文中设置的超时时间
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("服务优雅退出")

}
