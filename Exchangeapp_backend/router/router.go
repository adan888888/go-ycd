package router

import (
	"exchangeapp/controllers"
	_ "exchangeapp/docs" //引用docs.go
	"exchangeapp/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//参数1：其实就是路由访问的地址，如果你监听的 /assets那你就访问这个地址，后面可以跟你具体的那个路径所存在的路径是多少。
	//如：http://127.0.0.1:3000/assets/img1/img2.jpeg,实际，img2文件夹是没有监听的他是动态生成的这个文件。所以可以直接动态访问。
	//参数2：其实就是你监听的是哪个文件夹的名字，以及那个文件夹所在的路径。
	r.Static("/assets", "./assets") //图片访问

	r.Use(cors.New(cors.Config{
		// 修改这里 - 允许前端应用访问
		AllowOrigins: []string{"*"}, // 允许所有源访问,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Device-Type",
			"X-Device-Id",
			"X-Lang",
			"X-Platform-Id",
			"X-App-Terminal-Id",
			"UserId",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //swagger

	// 第1组：不需要认证
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// 第2组：不需要认证
	api := r.Group("/api")
	{
		//汇率换算 部分
		api.GET("/exchangeRates", controllers.GetExchangeRates)
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticleByID)
		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
		api.GET("/banners", controllers.GetBanners)
		api.GET("/hotgames", controllers.GetHotgames)
		api.GET("/testmq/:msg", controllers.SendRabbitMsg) //http://localhost:3000/api/testmq/你好
	}
	// 第3组：需要认证
	ycd := r.Group("/api/ycd")
	ycd.Use(middlewares.AuthMiddleWare()) // 添加认证中间件
	{
		//ycd 部分
		ycd.POST("/createtable", controllers.CreateTables)
		ycd.GET("/table1", controllers.GetTable1)
		ycd.GET("/table2", controllers.GetTable2)
		ycd.PUT("/inserttable1", controllers.InsertTable1)
		ycd.PUT("/inserttable2", controllers.InsertTable2)
		ycd.DELETE("/deletelast", controllers.DeleteLast)
		ycd.POST("/restart", controllers.Restart)
		ycd.POST("/sortxiaoshu", controllers.SortXiaoShu)
		ycd.POST("/xiaoshu", controllers.Xiaoshu) //消数
		ycd.DELETE("/deleteall", controllers.DeleteAll)
		ycd.POST("/resetliushui", controllers.ResetLiushui)
		ycd.POST("/updateqiwangvalue", controllers.UpdateQiWangValue)
		ycd.POST("/updateodds", controllers.UpdateOdds)
		ycd.POST("/updatebenjin", controllers.UpdateBenjin)
		ycd.GET("/getusers", controllers.Getusers)
		ycd.GET("/loadmore", controllers.LoadMore) //加载更多历史数据 //http://localhost:3000/api/ycd/loadMore?last_value=836
		ycd.GET("/getStatisticalAreasData", controllers.GetStatisticalAreasData)
		ycd.GET("/linechartData", controllers.LinechartData)              //折线图数据
		ycd.POST("/cleanDataD", controllers.CleanDataD)                   //清除数据（消数列数据全部清除）
		ycd.GET("/randomBankerPlayer", controllers.GetRandomBankerPlayer) //随机庄闲接口
		// 管理后台统计与列表（需登录）
		ycd.GET("/today/users", controllers.GetTodayBettingUsers)
		ycd.GET("/today/amount", controllers.GetTodayBettingAmount)
		ycd.GET("/today/count", controllers.GetTodayBettingCount)
		ycd.GET("/stats", controllers.GetBettingStats)
		ycd.GET("/zhuangzhanbi", controllers.GetZhuangZhanBi)
		ycd.POST("/zhuangzhanbi", controllers.UpdateZhuangZhanBiPublic)
		ycd.POST("/updatezhuangzhanbi", controllers.UpdateZhuangZhanBi)
		ycd.GET("/table1/list", controllers.GetTable1List)
		ycd.PUT("/table1/config", controllers.UpdateTable1Config)
		ycd.GET("/betting-record/list", controllers.GetTable2List)
		ycd.GET("/betting-record/by-id", controllers.GetTable2ByID)
		ycd.PUT("/betting-record/config", controllers.UpdateTable2Config)
	}

	// 第4组：买币记录管理（需认证）
	buyRecords := r.Group("/api/buyRecords")
	buyRecords.Use(middlewares.AuthMiddleWare())
	{
		buyRecords.GET("", controllers.GetBuyRecords)           // 获取买币记录列表
		buyRecords.POST("", controllers.CreateBuyRecord)      // 录入买币记录
		buyRecords.DELETE("/:id", controllers.DeleteBuyRecord) // 删除买币记录
	}

	// 第5组：密码本管理（无需认证）
	passwordBook := r.Group("/api/password-book")
	{
		passwordBook.POST("", controllers.CreatePasswordItem)                    // 创建密码项
		passwordBook.GET("", controllers.GetPasswordItems)                       // 获取密码列表
		passwordBook.GET("/:id", controllers.GetPasswordItem)                    // 获取单个密码项
		passwordBook.PUT("/:id", controllers.UpdatePasswordItem)                 // 更新密码项
		passwordBook.DELETE("/:id", controllers.DeletePasswordItem)              // 删除密码项
		passwordBook.POST("/batch-delete", controllers.BatchDeletePasswordItems) // 批量删除密码项
	}

	// 第6组：数据库备份管理（无需认证）
	backupController := controllers.NewBackupController()
	backup := r.Group("/api/backup")
	{
		backup.POST("/manual", backupController.ManualBackup)     // 手动触发备份
		backup.GET("/list", backupController.GetBackupList)       // 获取备份文件列表
		backup.GET("/status", backupController.GetBackupStatus)   // 获取备份状态
		backup.DELETE("/clean", backupController.CleanOldBackups) // 清理旧备份
	}

	//cookie
	index := r.Group("/index")
	index.Use(middlewares.CheckUser)
	{
		index.GET("/test", func(context *gin.Context) {
			context.JSON(200, gin.H{"msg": "成功！"})
		})
	}

	return r

}
