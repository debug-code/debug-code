package routers

import (
	"debug-code/controllers"
	"debug-code/lib"
	"debug-code/tools/jwt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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

	//token check
	var checkToken = func(ctx *context.Context) {

		mySigningKey := []byte("hzwy23")
		token := ctx.Request.Header.Get("Authorization")
		//fmt.Println("Authorization", token)

		if token == "" {
			rmsg := lib.RetrunMessage("401.1", "token失效或错误", nil)
			ctx.Output.JSON(&rmsg, true, false)
			return
		}
		resToken := jwt.CheckToken(mySigningKey, token)
		if resToken.Valid {
			//token ok
			return
		} else {
			rmsg := lib.RetrunMessage("401.1", "token失效或错误", nil)
			ctx.Output.JSON(&rmsg, true, false)
			return
		}

	}

	//static
	beego.InsertFilter("/api/*", beego.BeforeRouter, checkToken)
	beego.InsertFilter("/static", beego.BeforeRouter, checkToken)

	//views
	beego.Router("/open/push", &controllers.PushController{})

	//Test
	beego.Router("/open/test", &controllers.TestController{})

	//open
	//login
	beego.Router("/open/login", &controllers.LoginController{})
	//articles
	beego.Router("/open/articles/?:id", &controllers.OpenArticlesController{})

	//ip search info
	beego.Router("/open/ip/?:ip", &controllers.IpToolController{})

	beego.Router("/open/socket", &controllers.MyWebSocketController{})

	//api
	//articles
	beego.Router("/api/articles/?:id", &controllers.ArticleController{})

	//articleTypes
	beego.Router("/api/articleTypes/?:id", &controllers.ArticleTypeController{})

	//token
	beego.Router("/api/token", &controllers.TokenController{})

}
