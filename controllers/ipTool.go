package controllers

import (
	"geekerblog/lib"
	"github.com/astaxie/beego"
)

type IpToolController struct {
	beego.Controller
}

func (this *IpToolController) Get() {

	ip := this.Ctx.Input.Param(":ip")
	beego.Info("ip", ip)

	req := this.Ctx.Request.RemoteAddr
	beego.Info(req)
	beego.Info(this.Ctx.Request.Header.Get("x-Real-ip"))
	beego.Info(this.Ctx.Input.IP())

	lib.ReturnJson(this.Controller, "1", req)
	return
}
