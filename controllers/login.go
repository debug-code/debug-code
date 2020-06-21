package controllers

import (
	"debug-code/lib"
	"debug-code/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Post() {

	beego.Info("/open/login Post")
	//get data
	body := this.Ctx.Input.RequestBody
	manager := models.Manager{}
	err := json.Unmarshal(body, &manager)
	if err != nil {
		fmt.Println(err)
	}

	//check data
	check, err := lib.CheckManager(manager)
	if err != nil {
		beego.Error(" lib.CheckManager(manager) error :", err)
		lib.ReturnJson(this.Controller, "-1", check)
		return
	}

	//get token
	token, err := lib.GetToken(check)
	if err != nil {
		beego.Error(" lib.GetToken() error :", err)
		lib.ReturnJson(this.Controller, "-1", check)
		return
	}

	//return
	res := map[string]string{}
	res["token"] = token
	lib.ReturnJson(this.Controller, "1", res)
	return
}
