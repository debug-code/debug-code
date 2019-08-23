package controllers

import "github.com/astaxie/beego"

type IndexController struct {
	beego.Controller
}

func (mc *IndexController) Get() {
	//mc.Render("/index.html")

	mc.TplName = "socket.html"
	//mc.Ctx.WriteString("hello world")

}