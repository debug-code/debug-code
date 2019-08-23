package routers

import (
	"geekerblog/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {

	//跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	//static
	beego.SetStaticPath("/static", "static")

	//views
	beego.Router("/view/index", &controllers.IndexController{})

	//Test
	beego.Router("/open/test", &controllers.TestController{})

	//open
	//login
	beego.Router("/open/login", &controllers.LoginController{})
	//ip search info
	beego.Router("/open/ip/?:ip", &controllers.IpToolController{})

	beego.Router("/open/socket", &controllers.MyWebSocketController{})

	//api
	//articles
	beego.Router("/api/articles/?:id", &controllers.ArticleController{})

	//articleTypes
	beego.Router("/api/articleTypes/?:id", &controllers.ArticleTypeController{})

	//tools
	//sql to go struct
	beego.Router("/tools/stgs/?:id", &controllers.StgsController{})
}
