package controllers

import (
	"debug-code/lib"
	"fmt"
	"github.com/astaxie/beego"
)

type PushController struct {
	beego.Controller
}

func (this *PushController) Post() {

	body := this.Ctx.Input.RequestBody
	fmt.Println(string(body))
	fmt.Println(this.Ctx.Request.RemoteAddr)
	fmt.Println(this.Ctx.Request.Header.Get("X-Real-ip"))
	lib.ReturnJson(this.Controller, "1", nil)
	return

}

func (this *PushController) Get() {

	body := this.Ctx.Input.RequestBody
	fmt.Println(this.Ctx.Request.Header.Get("X-Real-ip"))
	fmt.Println(this.Ctx.Request.RemoteAddr)
	fmt.Println(string(body))
	lib.ReturnJson(this.Controller, "1", nil)
	return

}
