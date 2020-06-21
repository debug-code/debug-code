package controllers

import (
	"debug-code/lib"
	"github.com/astaxie/beego"
)

type TokenController struct {
	beego.Controller
}

func (this *TokenController) Get() {

	token, err := lib.GetTokenByContorller(this.Controller)
	if err != nil {
		beego.Error("Get, GetTokenByContorller err:", err)
		lib.ReturnJson(this.Controller, "-1", "token 获取失败")
		return
	}

	//成功返回token
	res := map[string]string{}
	res["token"] = token
	lib.ReturnJson(this.Controller, "1", res)
	return

}
