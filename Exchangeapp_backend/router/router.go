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
	// 第3组：需要认证（jsq 计数器；/api/ycd 为旧路径兼容）
	mountJsqRoutes(r.Group("/api/jsq"))
	mountJsqRoutes(r.Group("/api/ycd"))

	// 第4组：用户管理（仅超级管理员 Admin，需认证）
	adminUsers := r.Group("/api/admin/users")
	adminUsers.Use(middlewares.AuthMiddleWare())
	{
		adminUsers.GET("", controllers.AdminListUsers)
		adminUsers.DELETE("/:uid", controllers.AdminDeleteUser)
		adminUsers.POST("/:uid/restore", controllers.AdminRestoreUser)
		adminUsers.PUT("/:uid/username", controllers.AdminUpdateUsername)
		adminUsers.PUT("/:uid/password", controllers.AdminUpdatePassword)
		adminUsers.PUT("/:uid/expires-at", controllers.AdminUpdateExpiresAt)
		adminUsers.PUT("/:uid/role", controllers.AdminUpdateUserRole)
	}

	// 第5组：买币记录（需登录且专业版及以上；数据始终按当前 JWT 登录 uid 隔离）
	buyRecords := r.Group("/api/buyRecords")
	buyRecords.Use(middlewares.AuthMiddleWare())
	buyRecords.Use(middlewares.ProOrAboveMiddleware())
	{
		buyRecords.GET("", controllers.GetBuyRecords)          // 获取买币记录列表
		buyRecords.POST("", controllers.CreateBuyRecord)       // 录入买币记录
		buyRecords.DELETE("/:id", controllers.DeleteBuyRecord) // 删除买币记录
	}

	// 第6组：百家乐开奖模拟（需登录）
	baccaratGroup := r.Group("/api/baccarat")
	baccaratGroup.Use(middlewares.AuthMiddleWare())
	{
		baccaratGroup.GET("/state", controllers.GetBaccaratState)
		baccaratGroup.POST("/shuffle", controllers.ShuffleBaccaratShoe)
		baccaratGroup.POST("/cut-card", controllers.CutBaccaratCard)
		baccaratGroup.POST("/deal", controllers.DealBaccaratHand)
		baccaratGroup.POST("/reset", controllers.ResetBaccaratSession)
		baccaratGroup.POST("/bulk-simulate", controllers.BulkSimulateBaccarat)
		baccaratGroup.POST("/bulk-collision", controllers.BulkCollisionSimulate)
		baccaratGroup.POST("/bulk-mean-reversion", controllers.BulkMeanReversionSimulate)
		baccaratGroup.POST("/bulk-cable", controllers.BulkCableSimulate)
	}

	// 第7组：密码本管理（需登录且专业版及以上）
	passwordBook := r.Group("/api/password-book")
	passwordBook.Use(middlewares.AuthMiddleWare())
	passwordBook.Use(middlewares.ProOrAboveMiddleware())
	{
		passwordBook.POST("", controllers.CreatePasswordItem)                    // 创建密码项
		passwordBook.GET("", controllers.GetPasswordItems)                       // 获取密码列表
		passwordBook.GET("/:id", controllers.GetPasswordItem)                    // 获取单个密码项
		passwordBook.PUT("/:id", controllers.UpdatePasswordItem)                 // 更新密码项
		passwordBook.DELETE("/:id", controllers.DeletePasswordItem)              // 删除密码项
		passwordBook.POST("/batch-delete", controllers.BatchDeletePasswordItems) // 批量删除密码项
	}

	// 第8组：数据库备份管理（无需认证）
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
			controllers.Success(context, "成功！", gin.H{})
		})
	}

	return r

}

// mountJsqRoutes 挂载计数器（jsq）相关路由，供 Flutter App 与 Vue 后台共用；需登录且订阅有效。
func mountJsqRoutes(g *gin.RouterGroup) {
	g.Use(middlewares.AuthMiddleWare())
	g.Use(middlewares.JsqSubscriptionMiddleware())

	// 表初始化与基础读写（table1 配置、table2 投注明细）
	g.POST("/createtable", controllers.CreateTables)
	g.GET("/table1", controllers.GetTable1)
	g.GET("/table2", controllers.GetTable2)
	g.PUT("/inserttable1", controllers.InsertTable1)
	g.PUT("/inserttable2", controllers.InsertTable2)
	g.PUT("/updaterestartstatsnapshot", controllers.UpdateLastRowRestartStatSnapshot)

	// 下注流程：删最后一笔、重启、消数、清空等
	g.DELETE("/deletelast", controllers.DeleteLast)
	g.POST("/restart", controllers.Restart)
	g.POST("/sortxiaoshu", controllers.SortXiaoShu)
	g.POST("/xiaoshu", controllers.Xiaoshu)
	g.DELETE("/deleteall", controllers.DeleteAll)
	g.POST("/resetliushui", controllers.ResetLiushui)

	// 参数配置：期望值、赔率、本金
	g.POST("/updateqiwangvalue", controllers.UpdateQiWangValue)
	g.POST("/updateodds", controllers.UpdateOdds)
	g.POST("/updatebenjin", controllers.UpdateBenjin)

	// App 端：用户列表、分页加载、统计区、折线图、消数值清理
	g.GET("/getusers", controllers.Getusers)
	g.GET("/loadmore", controllers.LoadMore)
	g.GET("/getStatisticalAreasData", controllers.GetStatisticalAreasData)
	g.GET("/linechartData", controllers.LinechartData)
	g.POST("/cleanDataD", controllers.CleanDataD)

	// 随机庄闲、今日投注概览、按日期区间统计
	g.GET("/randomBankerPlayer", controllers.GetRandomBankerPlayer)
	g.GET("/today/users", controllers.GetTodayBettingUsers)
	g.GET("/today/amount", controllers.GetTodayBettingAmount)
	g.GET("/today/count", controllers.GetTodayBettingCount)
	g.GET("/stats", controllers.GetBettingStats)

	// 庄占比查询与更新
	g.GET("/zhuangzhanbi", controllers.GetZhuangZhanBi)
	g.POST("/zhuangzhanbi", controllers.UpdateZhuangZhanBiPublic)
	g.POST("/updatezhuangzhanbi", controllers.UpdateZhuangZhanBi)

	// Vue 后台：table1 列表与配置
	g.GET("/table1/list", controllers.GetTable1List)
	g.PUT("/table1/config", controllers.UpdateTable1Config)

	// Vue 后台：投注记录列表、数据统计聚合、按 ID 检索、单条编辑
	g.GET("/betting-record/list", controllers.GetTable2List)
	g.GET("/betting-record/data-stats", controllers.GetDataStats)
	g.GET("/betting-record/by-id", controllers.GetTable2ByID)
	g.PUT("/betting-record/config", controllers.UpdateTable2Config)
}
