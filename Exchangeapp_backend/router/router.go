package router

import (
	"exchangeapp/controllers"
	_ "exchangeapp/docs" //引用docs.go
	"exchangeapp/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//参数1：其实就是路由访问的地址，如果你监听的 /assets那你就访问这个地址，后面可以跟你具体的那个路径所存在的路径是多少。
	//如：http://127.0.0.1:3000/assets/img1/img2.jpeg,实际，img2文件夹是没有监听的他是动态生成的这个文件。所以可以直接动态访问。
	//参数2：其实就是你监听的是哪个文件夹的名字，以及那个文件夹所在的路径。
	r.Static("/assets", "./assets") //图片访问

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //swagger

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)

		auth.POST("/register", controllers.Register)
	}

	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchangeRates)
	//api.Use(middlewares.AuthMiddleWare())
	{
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
	//cookie
	index := r.Group("/index")
	index.Use(middlewares.CheckUser)
	{
		index.GET("/test", func(context *gin.Context) {
			context.JSON(200, gin.H{"msg": "成功！"})
		})
	}

	//统计数据
	{
		api.POST("/ycd/createtable", controllers.CreateTables)
		api.GET("/ycd/table1", controllers.GetTable1)
		api.GET("/ycd/table2", controllers.GetTable2)
		api.PUT("/ycd/inserttable1", controllers.InsertTable1)
		api.PUT("/ycd/inserttable2", controllers.InsertTable2)
		api.DELETE("/ycd/deletelast", controllers.DeleteLast)
		api.POST("/ycd/restart", controllers.Restart)
		api.POST("/ycd/sortxiaoshu", controllers.SortXiaoShu)
		api.POST("/ycd/xiaoshu", controllers.Xiaoshu)
		api.DELETE("/ycd/deleteall", controllers.DeleteAll)
		api.POST("/ycd/resetliushui", controllers.ResetLiushui)
		api.POST("/ycd/updateqiwangvalue", controllers.Updateqiwangvalue)
		api.POST("/ycd/updateodds", controllers.UpdateOdds)
		api.POST("/ycd/updatebenjin", controllers.UpdateBenjin)
		api.GET("/ycd/getusers", controllers.Getusers)
	}
	return r

}
