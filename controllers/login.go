package controllers

import (
	"encoding/json"
	"fmt"
	"geekerblog/lib"
	"geekerblog/models"
	"geekerblog/tools/jwt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Post() {

	//get data
	body := this.Ctx.Input.RequestBody
	beego.Info("body", string(body))

	manager := models.Manager{}
	err := json.Unmarshal(body, &manager)
	if err != nil {
		fmt.Println(err)
	}

	//search manager from database
	o := orm.NewOrm()
	managerS := models.Manager{}
	err = o.QueryTable("manager").Filter("status", 1).
		Filter("account", manager.Account).One(&managerS)
	//fmt.Println(managerS)
	if err != nil || manager.Account != managerS.Account || managerS.Account == "" {
		beego.Error(err)
		lib.ReturnJson(this.Controller, "0", err)
		return
	}
	if manager.Passwd != managerS.Passwd {
		beego.Error(err)
		lib.ReturnJson(this.Controller, "0", err)
		return
	}

	//jwt
	mySigningKey := []byte("hzwy23")
	ss, err := jwt.GetNewToken(mySigningKey, strconv.Itoa(managerS.Id))
	if err != nil {
		lib.ReturnJson(this.Controller, "0", err)
		return
	}

	//return
	res := map[string]string{}
	res["token"] = ss
	lib.ReturnJson(this.Controller, "1", res)
	return
}
func (this *LoginController) Get() {

	tokens := this.Ctx.Input.Context.Request.Header.Get("Authorization")
	id, _ := jwt.GetTokenInfo(tokens)

	//jwt
	mySigningKey := []byte("hzwy23")
	ss, err := jwt.GetNewToken(mySigningKey, id)
	if err != nil {
		lib.ReturnJson(this.Controller, "0", err)
		return
	}

	//登录成功返回信息
	res := map[string]string{}
	res["token"] = ss
	lib.ReturnJson(this.Controller, "1", res)
	return

}
