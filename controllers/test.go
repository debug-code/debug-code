package controllers

import (
	"geekerblog/lib"
	"geekerblog/models"
	"github.com/astaxie/beego"
	"github.com/debug-code/goDB/mysql"
)

type TestController struct {
	beego.Controller
}


func (this *TestController) Get(){

	db, err := mysql.Context()
	if err != nil {
		beego.Error("fac err:", err)
		lib.ReturnJson(this.Controller, "-1",nil)
		return
	}
	defer mysql.ReleaseC(db)

	manager := []models.Manager{}
	db.Limit(10).Find(&manager)
	lib.ReturnJson(this.Controller, "1",manager)
	return

}